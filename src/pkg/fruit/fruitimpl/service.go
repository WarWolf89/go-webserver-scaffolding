package fruitimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"csaba.almasi.per/webserver/src/pkg/fruit"
)

type RedisService struct {
	Client *redis.Client
}

// type RedisStore struct {
// 	Client *redis.Client
// }

// type SQLStore struct {
// 	client *sql.DB
// }

// type FruitStore interface {
// 	GetFruit(context.Context) (*fruit.Fruit, error)
// }

// interface compliance check https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ fruit.Service = (*RedisService)(nil)

func ProvideSVC() fruit.Service {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rsvc := &RedisService{
		Client: client,
	}
	return rsvc
}

func (rsvc *RedisService) GetFruits(ctx context.Context) ([]*fruit.Fruit, error) {
	fruits := []*fruit.Fruit{}
	val, err := rsvc.Client.HGetAll(ctx, "basket").Result()
	if err != nil {
		slog.Error("Failed fetching item keys", err)
		return nil, err
	}

	for _, item := range val {
		fruit := &fruit.Fruit{}
		err := json.Unmarshal([]byte(item), fruit)
		if err != nil {
			slog.Error("Failed marshaling items", err)
			return nil, err
		}
		fruits = append(fruits, fruit)
	}
	return fruits, nil
}

func (rsvc *RedisService) AddFruit(ctx context.Context, fruit *fruit.Fruit) (string, error) {
	fruit.ID = uuid.NewString()

	json, err := json.Marshal(fruit)
	if err != nil {
		slog.Error("Error marshalling fruit", err)
		return "", err
	}

	// Basket is the grouping hash key we use for all the fruits
	if err := rsvc.Client.HSet(ctx, "basket", fmt.Sprintf("fruit:%s", fruit.ID), json).Err(); err != nil {
		slog.Error("Error writing to Redis:", err)
		return "", err
	}

	return fruit.ID, nil
}

func (rsvc *RedisService) GetFruitByID(ctx context.Context, id string) (*fruit.Fruit, error) {

	fruit := &fruit.Fruit{}

	val, err := rsvc.Client.HGet(ctx, "basket", fmt.Sprintf("fruit:%s", id)).Result()
	if err != nil {
		slog.Error("Error from Redis lookup for fruit:", id, "\n with error:", err)
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &fruit)
	if err != nil {
		slog.Error("Failed marshaling items:", val, "with error:", err)
		return nil, err
	}
	return fruit, nil
}
