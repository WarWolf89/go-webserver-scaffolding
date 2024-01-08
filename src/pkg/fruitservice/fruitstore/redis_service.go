package fruitstore

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"csaba.almasi.per/webserver/src/pkg/fruitservice"
)

// Basket is the  primary grouping hash key we use for all the fruits in Redis
const primary_key = "basket"

type RedisStore struct {
	Client *redis.Client
}

// interface compliance check https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ fruitservice.FruitService = (*RedisStore)(nil)

func ProvideSVC() fruitservice.FruitService {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rsvc := &RedisStore{
		Client: client,
	}
	return rsvc
}

func (rsvc *RedisStore) GetFruits(ctx context.Context) ([]*fruitservice.Fruit, error) {
	fruits := []*fruitservice.Fruit{}
	val, err := rsvc.Client.HGetAll(ctx, primary_key).Result()
	if err != nil {
		slog.Error("Failed fetching item keys", err)
		return nil, err
	}

	for _, item := range val {
		fruit := &fruitservice.Fruit{}
		err := json.Unmarshal([]byte(item), fruit)
		if err != nil {
			slog.Error("Failed marshaling items", err)
			return nil, err
		}
		fruits = append(fruits, fruit)
	}
	return fruits, nil
}

func (rsvc *RedisStore) AddFruit(ctx context.Context, fruit *fruitservice.Fruit) (string, error) {
	fruit.ID = uuid.NewString()

	json, err := json.Marshal(fruit)
	if err != nil {
		slog.Error("Error marshalling fruit", err)
		return "", err
	}

	if err := rsvc.Client.HSet(ctx, primary_key, fmt.Sprintf("fruit:%s", fruit.ID), json).Err(); err != nil {
		slog.Error("Error writing to Redis:", err)
		return "", err
	}

	return fruit.ID, nil
}

func (rsvc *RedisStore) GetFruitByID(ctx context.Context, id string) (*fruitservice.Fruit, error) {

	fruit := &fruitservice.Fruit{}
	// gin gets patparam with ':' hence omitting here, the key here is fruit:<id>
	val, err := rsvc.Client.HGet(ctx, primary_key, fmt.Sprintf("fruit%s", id)).Result()
	if err != nil {
		slog.Error("Error from Redis lookup for fruit:", id, "\n with error:", err)
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &fruit)
	if err != nil {
		slog.Error("Failed marshaling items:", "value", val, "error:", err)
		return nil, err
	}
	return fruit, nil
}
