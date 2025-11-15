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

	userQueries  queries.UserQueries
	itemQueries  queries.ItemQueries
	orderQueries queries.OrderQueries
)

func getDB() *gorm.DB {
	mu.Do(func() {
		_db, err := gorm.Open(sqlite.Open("userorder.db"), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
		db = _db.Debug()

		userQueries = queries.NewUserQueries(db)
		itemQueries = queries.NewItemQueries(db)
		orderQueries = queries.NewOrderQueries(db)
	})
	return db
}

func createUser(ctx context.Context, number uint) (*models.User, error) {
	user := &models.User{
		Username: fmt.Sprintf("user-%d", number),
		Balance:  1000 * float64(number),
	}
	if err := userQueries.CreateOne(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func userOrderItems(ctx context.Context, userId uint) (*models.Order, error) {
	var (
		order = &models.Order{
			Discount: 0.5,
			UserID:   userId,
		}

		items = []*models.Item{
			{OrderID: order.ID, Product: "item-1", Quantity: 1, Price: 10.50},
			{OrderID: order.ID, Product: "item-2", Quantity: 1, Price: 10.50},
			{OrderID: order.ID, Product: "item-3", Quantity: 1, Price: 10.50},
		}
	)

	db.Transaction(gormqs.Tx(func(tx *gorm.DB) error {
		ctx := tx.Statement.Context

		// get and lock user
		user, err := userQueries.GetOne(ctx, qsopt.LockForUpdate(), qsopt.WhereID(userId))
		if err != nil {
			return err
		}

		for _, item := range items {
			item.OrderID = order.ID
			order.Amount += item.Price * float64(item.Quantity)
		}

		// create items
		if err := itemQueries.CreateMany(ctx, &items); err != nil {
			return err
		}

		// create an order
		order.UserID = user.ID
		if err := orderQueries.CreateOne(ctx, order); err != nil {
			return err
		}

		// deduct user balance
		user.Balance -= (order.Amount * (1 - order.Discount))
		if _, err := userQueries.Updates(ctx, user, qsopt.UserSelect(qsopt.UserBalance)); err != nil {
			return err
		}

		return nil
	}))

	return order, err
}

func main() {
	ctx := context.Background()
	db := getDB()
	err := db.AutoMigrate(&models.User{}, &models.Item{}, &models.Order{})
	mustNotErr(err)

	user1, err := createUser(ctx, 1)
	mustNotErr(err)

	order, err := userOrderItems(ctx, user1.ID)
	mustNotErr(err)

	// custom query
	orderWithItem, err := orderQueries.Querier().GetOneWithItems(ctx, order.ID)
	mustNotErr(err)

	log.Println("orderWithItem:")
	printJson(orderWithItem)
}

func printJson(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	mustNotErr(err)
	fmt.Println(string(b))
}
