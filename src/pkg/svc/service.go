package fruitservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"

	models "csaba.almasi.per/webserver/src/pkg/models"
)

type RedisService struct {
	Client *redis.Client
}

func ProvideSVC() *RedisService {
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

func (rsvc *RedisService) GetFruits(ctx context.Context) ([]*models.Fruit, error) {
	fruits := []*models.Fruit{}
	val, err := rsvc.Client.HGetAll(ctx, "basket").Result()
	if err != nil {
		slog.Error("Failed fetching item keys", err)
		return nil, err
	}

	for _, item := range val {
		fruit := &models.Fruit{}
		err := json.Unmarshal([]byte(item), fruit)
		if err != nil {
			slog.Error("Failed marshaling items", err)
			return nil, err
		}
		fruits = append(fruits, fruit)
	}
	return fruits, nil
}

func (rsvc *RedisService) AddFruit(ctx context.Context) (string, error) {

	// Get the payload from the passed down gin context
	json := ctx.Value("payload").([]byte)
	fuid := (ctx.Value("id")).(string)

	// Basket is the grouping hash key we use for all the fruits
	if err := rsvc.Client.HSet(ctx, "basket", fmt.Sprintf("fruit:%s", fuid), json).Err(); err != nil {
		slog.Error("Error writing to Redis:", err)
		return "", err
	}

	return fuid, nil
}

func (rsvc *RedisService) GetFruitByID(ctx context.Context) (*models.Fruit, error) {
	id := ctx.Value("id").(string)
	fruit := &models.Fruit{}

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
