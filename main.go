package main

import (

	// ever since slog made it into the core pkg there's no need for external loggers

	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	// validator
	"github.com/go-playground/validator/v10"
	// gin is the most used and standard web framework in go
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	// went with go-redis since it's better documented than redisgo
	"github.com/redis/go-redis/v9"
	// redis json helper lib
)

type Fruit struct {
	ID    string `json:"id"`
	Name  string `json:"name" validate:"alpha"`
	Color string `json:"color" validate:"iscolor"`
}

type Handler struct {
	redis_client *redis.Client
	validate     *validator.Validate
}

func (h *Handler) GetFruits(c *gin.Context) {

	fruits := []*Fruit{}
	val, err := h.redis_client.HGetAll(c, "basket").Result()
	if err != nil {
		slog.Error("Failed fetching item keys", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	for _, item := range val {
		fruit := &Fruit{}
		err := json.Unmarshal([]byte(item), fruit)
		if err != nil {
			slog.Error("Failed marshaling items", err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal server error")
		}
		fruits = append(fruits, fruit)
	}
	c.IndentedJSON(http.StatusOK, fruits)
}

func (h *Handler) AddFruit(c *gin.Context) {

	var fruit Fruit

	// Call BindJSON to bind the received JSON to
	// a new fruit.

	if err := c.ShouldBind(&fruit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// validate post request
	if err := h.validate.Struct(fruit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Request body has incorrect data: ": err.Error(),
		})
		return
	}

	fruit.ID = uuid.New().String()

	json, err := json.Marshal(fruit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Basket is the grouping hash key we use for all the fruits
	if err := h.redis_client.HSet(c, "basket", fmt.Sprintf("fruit:%s", fruit.ID), json).Err(); err != nil {
		slog.Error("Error writing to Redis:", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
	}

	c.IndentedJSON(http.StatusCreated, fruit)
}

func (h *Handler) GetFruitByID(c *gin.Context) {
	var fruit Fruit
	id, exists := c.Params.Get("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"no ID provided as a pathparam": id,
		})
		return
	}
	val, err := h.redis_client.HGet(c, "basket", fmt.Sprintf("fruit:%s", id)).Result()
	if err != nil {
		slog.Error("Error from Redis lookup", err)
		c.IndentedJSON(http.StatusInternalServerError, "Error when looking up fruit by id")
		return
	}

	err = json.Unmarshal([]byte(val), &fruit)

	if err != nil {
		slog.Error("Failed marshaling items", err)
		c.IndentedJSON(http.StatusInternalServerError, "Internal server error")
	}

	c.IndentedJSON(http.StatusOK, fruit)

}

func setupRedis() *redis.Client {
	// Set up the Redis Client
	redis_client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return redis_client
}

func setupClients() (*gin.Engine, *redis.Client) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	redis_client := setupRedis()
	router := gin.Default()
	group := router.Group("api/v1")
	group.GET("/fruits", (&Handler{redis_client: redis_client}).GetFruits)
	group.GET("/fruits/:id", (&Handler{redis_client: redis_client}).GetFruitByID)
	group.POST("/fruits", (&Handler{validate: validate, redis_client: redis_client}).AddFruit)
	return router, redis_client
}

func main() {
	// use a single instance of Validate, it caches struct info

	gin_router, redis_client := setupClients()
	if err := gin_router.Run(":8080"); err != nil {
		slog.Error("Failed setting up webserver:", err)
		os.Exit(1)
	}

	defer func() {
		if err := redis_client.Close(); err != nil {
			slog.Error("goredis - failed to communicate to redis-server: ", err)
		}
	}()

}
