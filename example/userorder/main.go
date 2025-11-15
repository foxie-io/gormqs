package main

import (
	"context"
	"encoding/json"
	"example/userorder/models"
	"example/userorder/queries"
	qsopt "example/userorder/queries/options"
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
		order_qs = queries.NewOrderQuerier(db)
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

func userOrderItemTransaction(userId uint, orderItems []*models.OrderItem, returnOrder *models.Order) func(*gorm.DB) error {
	return func(tx *gorm.DB) error {
		// context will carry transaction instance
		ctx := gormqs.ContextWithValue(tx.Statement.Context, tx)
		tx.Statement.Context = ctx

		order := &models.Order{
			UserID: userId,
		}

		// get and lock user
		user, err := user_qs.GetOne(ctx, qsopt.LockForUpdate(), qsopt.WhereID(userId))
		if err != nil {
			return err
		}

		// perepare order
		order.UserID = user.ID
		for _, orderItem := range orderItems {
			order.PayAmount += orderItem.PayAmount()
			order.DiscountAmount += orderItem.DiscountAmount()
		}

		// create an order
		if err := order_qs.CreateOne(ctx, order); err != nil {
			return err
		}

		// perepare order item
		for _, orderItem := range orderItems {
			orderItem.OrderID = order.ID
		}

		// create order items
		if err := orderItem_qs.CreateMany(ctx, &orderItems); err != nil {
			return err
		}

		// deduct user balance
		user.Balance -= order.PayAmount
		if _, err := user_qs.Updates(ctx, user, qsopt.UserSelect(qsopt.UserBalance)); err != nil {
			return err
		}

		*returnOrder = *order
		return nil
	}
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

	var order models.Order
	err = db.Transaction(userOrderItemTransaction(user1.ID, orderItems, &order))
	mustNotErr(err)

	// custom query
	orderWithItem, err := order_qs.GetOneWithOrderItems(ctx, order.ID)
	mustNotErr(err)

	log.Println("orderWithItem:")
	printJson(orderWithItem)

}

func printJson(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	mustNotErr(err)
	fmt.Println(string(b))
}
