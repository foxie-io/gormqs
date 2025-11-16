package main

import (
	"context"
	"encoding/json"
	"example/userorder/models"
	"example/userorder/queries"
	"fmt"
	"log"
	"sync"

	"github.com/foxie-io/gormqs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func mustNotErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	db *gorm.DB
	mu sync.Once

	user_qs      *queries.UserQueries
	item_qs      *queries.ItemQueries
	order_qs     *queries.OrderQueries
	orderItem_qs *queries.OrderItemQueries
)

func getDB() *gorm.DB {
	mu.Do(func() {
		_db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
		db = _db.Debug()

		user_qs = queries.NewUserQueries(db)
		item_qs = queries.NewItemQueries(db)
		order_qs = queries.NewOrderQueries(db)
		orderItem_qs = queries.NewOrderItemQueries(db)
	})
	return db
}

func createUser(ctx context.Context, number uint) (*models.User, error) {
	user := &models.User{
		Username: fmt.Sprintf("user-%d", number),
		Balance:  1000 * float64(number),
	}
	if err := user_qs.CreateOne(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func createItems(ctx context.Context) ([]*models.Item, error) {
	items := []*models.Item{
		{Product: "item-1", Quantity: 1, Price: 10.50},
		{Product: "item-2", Quantity: 1, Price: 10.50},
		{Product: "item-3", Quantity: 1, Price: 10.50},
	}

	if err := item_qs.CreateMany(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func createUserOrder(ctx context.Context, userId uint, orderItems []*models.OrderItem) (*models.Order, error) {
	order := &models.Order{
		UserID: userId,
	}

	for _, orderItem := range orderItems {
		order.PayAmount += orderItem.PayAmount()
		order.DiscountAmount += orderItem.DiscountAmount()
	}

	// tx1
	user, err := user_qs.BlockBalance(ctx, userId, order.PayAmount)
	if err != nil {
		return nil, err
	}

	// tx2
	err = getDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// wrap tx into ctx, so queries can use it as tx instance
		ctx := gormqs.WrapContext(tx)

		// create an order
		if err := order_qs.CreateOne(ctx, order); err != nil {
			return err
		}

		// create orderItems
		for _, orderItem := range orderItems {
			orderItem.OrderID = order.ID
		}

		if err := orderItem_qs.CreateMany(ctx, &orderItems); err != nil {
			return err
		}

		// commit blocked balance = success
		// will use tx2 instance because of tx2Ctx
		_, err = user_qs.CommitBlockedBalance(ctx, user.ID, order.PayAmount)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		// tx 3
		user, err = user_qs.UnblockBalance(ctx, userId, order.PayAmount)
		if err != nil {
			return nil, err
		}
	}

	return order, nil
}

func main() {
	ctx := context.Background()
	db := getDB()
	err := db.AutoMigrate(&models.User{}, &models.Item{}, &models.Order{}, &models.OrderItem{})
	mustNotErr(err)

	user1, err := createUser(ctx, 1)
	mustNotErr(err)

	items, err := createItems(ctx)
	mustNotErr(err)

	orderItems := make([]*models.OrderItem, len(items))
	for i, item := range items {
		orderItems[i] = &models.OrderItem{
			ItemID:   item.ID,
			Quantity: 1 + uint(i),
			Price:    item.Price,
			Discount: 0.1 * float64(i),
		}
	}

	order, err := createUserOrder(ctx, user1.ID, orderItems)
	mustNotErr(err)

	// custom query
	orderWithDetails, err := order_qs.GetOneWithDetails(ctx, order.ID)
	mustNotErr(err)

	log.Println("order with details:")
	printJson(orderWithDetails)
}

func printJson(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	mustNotErr(err)
	fmt.Println(string(b))
}
