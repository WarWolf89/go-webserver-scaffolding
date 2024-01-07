package api

import (
	models "csaba.almasi.per/webserver/src/pkg/models"
	svc "csaba.almasi.per/webserver/src/pkg/svc"

	// ever since slog made it into the core pkg there's no need for external loggers

	"encoding/json"
	"fmt"
	"net/http"

	// validator
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	// gin is the most used and standard web framework in go
	"github.com/gin-gonic/gin"
	// went with go-redis since it's better documented than redisgo
)

// move to api go
type Api struct {
	Gengine  *gin.Engine
	rscv     *svc.RedisService
	validate *validator.Validate
}

func ProvideApi(gengine *gin.Engine, rscv *svc.RedisService, validate *validator.Validate) *Api {
	api := &Api{
		Gengine:  gengine,
		rscv:     rscv,
		validate: validate,
	}
	return api
}

func (api *Api) RegisterAPIEndpoints() {
	ge := api.Gengine
	group := ge.Group("api/v1")
	group.GET("/fruits", api.GetFruits)
	group.GET("/fruits/:id", api.GetFruitByID)
	group.POST("/fruits", api.AddFruit)
}

func (api *Api) GetFruits(c *gin.Context) {
	fruits, err := api.rscv.GetFruits(c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Error when fetching all fruits.")
		return
	}
	c.IndentedJSON(http.StatusOK, fruits)
}

func (api *Api) AddFruit(c *gin.Context) {
	var fruit models.Fruit

	// Call BindJSON to bind the received JSON to
	// a new fruit.

	if err := c.ShouldBind(&fruit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// validate post request
	if err := api.validate.Struct(fruit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Request body has incorrect data: ": err.Error(),
		})
		return
	}

	fruit.ID = uuid.NewString()

	json, err := json.Marshal(fruit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Set("payload", json)
	c.Set("id", fruit.ID)

	id, err := api.rscv.AddFruit(c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Failed creating fruit.")
		return
	}

	// TODO add actual redirect URL to header instead of placeholder
	c.Header("Location", fmt.Sprintf("location/path/%s", id))
	c.IndentedJSON(http.StatusCreated, id)
}

func (api *Api) GetFruitByID(c *gin.Context) {

	id, exists := c.Params.Get("id")
	c.Set("id", id)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"no ID provided as a pathparam": id,
		})
		return
	}

	f, err := api.rscv.GetFruitByID(c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, fmt.Sprintf("Error getting fruit with id: %s", id))
		return
	}

	c.IndentedJSON(http.StatusOK, &f)

}
