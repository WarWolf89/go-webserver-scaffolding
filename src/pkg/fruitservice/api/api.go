package api

import (
	// ever since slog made it into the core pkg there's no need for external loggers
	"fmt"
	"log/slog"
	"net/http"

	// gin is the most used and standard web framework in go

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"csaba.almasi.per/webserver/src/pkg/fruitservice"
)

type Api struct {
	Gengine  *gin.Engine
	Fsvc     fruitservice.FruitService
	validate *validator.Validate
}

func ProvideApi(gengine *gin.Engine, rsvc fruitservice.FruitService, validate *validator.Validate) *Api {
	api := &Api{
		Gengine:  gengine,
		Fsvc:     rsvc,
		validate: validate,
	}
	return api
}

func (api *Api) RegisterAPIEndpoints() {
	v1 := api.Gengine.Group("api/v1")
	v1.GET("/fruits", api.GetFruits)
	v1.GET("/fruits:id", api.GetFruitByID)
	v1.POST("/fruits", api.AddFruit)
}

func (api *Api) GetFruits(c *gin.Context) {
	fruits, err := api.Fsvc.GetFruits(c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Error when fetching all fruits.")
		return
	}
	c.IndentedJSON(http.StatusOK, fruits)
}

func (api *Api) AddFruit(c *gin.Context) {
	var fruit fruitservice.Fruit

	// Call BindJSON to bind the received JSON to a new fruit.
	if err := c.ShouldBind(&fruit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate post request
	if err := api.validate.Struct(fruit); err != nil {
		slog.Warn("Validation failed for", "fruit", fruit, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Request body has incorrect data: ": err.Error(),
		})
		return
	}

	id, err := api.Fsvc.AddFruit(c, &fruit)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Failed creating fruit.")
		return
	}

	// Full url is cleaner, but this should be good enough
	c.Header("Location", fmt.Sprintf("v1/api/fruits:%s", id))
	c.IndentedJSON(http.StatusCreated, id)
}

func (api *Api) GetFruitByID(c *gin.Context) {

	id, exists := c.Params.Get("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"no ID provided as a pathparam": id,
		})
		return
	}

	f, err := api.Fsvc.GetFruitByID(c, id)
	if err != nil {
		slog.Error("Failed getting fruit from Redis", "error", err)
		c.IndentedJSON(http.StatusInternalServerError, fmt.Sprintf("Error getting fruit with id: %s", id))
		return
	}

	c.IndentedJSON(http.StatusOK, &f)
}
