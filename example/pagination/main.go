package main

import (
	"context"
	"encoding/json"
	"example/pagination/dto"
	"example/pagination/models"
	"example/pagination/queries"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/foxie-io/gormqs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
	mu sync.Once

	user_qs *queries.UserQueries
)

func getDB() *gorm.DB {
	mu.Do(func() {
		_db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

		if err != nil {
			log.Fatal(err)
		}
		db = _db.Debug()
		db.AutoMigrate(&models.User{})
		user_qs = queries.NewUserQueries(db)

	})
	return db
}

func prepare100Users() {
	users := []*models.User{}
	for i := 0; i < 100; i++ {
		users = append(users, &models.User{
			Username: "user_" + fmt.Sprintf("%v", i),
			Balance:  100,
		})
	}

	user_qs.CreateMany(context.Background(), &users)
}

func responseJson(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	getDB()

	prepare100Users()

	mux := http.NewServeMux()

	// list and count : `curl 'localhost:8080/users/page'`
	// count only : `curl 'localhost:8080/users/page?select=count'`
	// list only : `curl 'localhost:8080/users/page?select=list'`

	mux.HandleFunc("GET /users/page", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		pageParam := dto.NewPageParam[[]dto.BaseUser](r)

		// pageParam implement gormqs.ManyWithCountResulter
		// for pagination, sometimes we don't need to count, becuase count cost alot of performance
		if err := user_qs.GetManyWithCount(ctx, pageParam, pageParam.DBOption()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result := pageParam.Result()
		if result.Data != nil {
			hasNext := pageParam.HasNext(len(*result.Data))
			result.HasNext = &hasNext
		}

		responseJson(w, result)
	})

	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var countUsers int64

		// SELECT * FROM `users` LIMIT 2
		users, err := user_qs.GetMany(ctx, gormqs.Count(&countUsers), gormqs.LimitAndOffset(2, 0))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var baseUsers []dto.BaseUser
		var countBaseUser int64

		// dynamic select and maping
		// SELECT `users`.`id`,`users`.`username` FROM `users` LIMIT 1
		if err := user_qs.GetManyTo(ctx, &baseUsers, gormqs.Count(&countBaseUser), gormqs.LimitAndOffset(1, 0)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseJson(w, map[string]any{
			"users":          users,
			"totalUsers":     countUsers,
			"baseUsers":      baseUsers,
			"totalBaseUsers": countBaseUser,
		})
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
