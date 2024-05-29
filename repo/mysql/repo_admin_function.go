package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MinhSang97/order_app/model/admin_model"
	"github.com/MinhSang97/order_app/redis"
	"github.com/MinhSang97/order_app/repo"
	"gorm.io/gorm"
	"log"
)

type adminFunctionRepository struct {
	db *gorm.DB
}

var RedisClient = redis.ConnectRedis()

func (s adminFunctionRepository) GetAll(ctx context.Context) ([]admin_model.AdminFunctionModel, error) {
	var users []admin_model.AdminFunctionModel

	//RedisClient := redis.ConnectRedis()
	//
	// Đọc danh sách users từ Redis (nếu có)
	cachedUsersJSON, err := RedisClient.Get(ctx, "users").Result()
	if err == nil {
		var cachedUserss []admin_model.AdminFunctionModel
		err := json.Unmarshal([]byte(cachedUsersJSON), &cachedUserss)
		if err != nil {
			log.Println("Failed to unmarshal users from Redis:", err)
			return users, fmt.Errorf("Failed to unmarshal users from Redis: %w", err)
		}
		fmt.Println("Users fetched from Redis")
		return cachedUserss, nil
	}

	// Nếu không tìm thấy trong Redis, đọc từ cơ sở dữ liệu MySQL
	if err := s.db.Table("users").Scan(&users).Error; err != nil {
		return users, fmt.Errorf("get all users error: %w", err)
	}

	// Cache danh sách sinh viên vào Redis
	jsonUsers, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Failed to marshal students:", err)
		return users, fmt.Errorf("Failed to marshal users: %w", err)
	}

	err = redis.RedisClient.Set(ctx, "users", jsonUsers, 0).Err()
	if err != nil {
		fmt.Println("Failed to cache users in Redis:", err)
	}

	fmt.Println("Students query from MySQL")

	return users, nil
}

var instance adminFunctionRepository

func NewAdminFunctionUseCase(db *gorm.DB) repo.AdminFunctionRepo {
	if instance.db == nil {
		instance.db = db
	}
	return instance
}
