package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type fruit struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

var fruits = []fruit{
	{ID: "1", Name: "apple", Color: "red"},
	{ID: "2", Name: "orange", Color: "orange"},
	{ID: "3", Name: "pineapple", Color: "yellow"},
}

// getAlbums responds with the list of all albums as JSON.
func getFruits(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, fruits)
}

func addFruit(c *gin.Context) {
	var fruit fruit

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&fruit); err != nil {
		return
	}

	// Add the new album to the slice.
	fruits = append(fruits, fruit)
	c.IndentedJSON(http.StatusCreated, fruit)
}

func getFruitByID(c *gin.Context) {
	id, _ := c.Params.Get("id")

	// Loop over the list of fruits, looking for
	// a fruit whose ID value matches the parameter.
	for _, f := range fruits {
		if f.ID == id {
			c.IndentedJSON(http.StatusOK, f)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "fruit not found"})
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/fruits", getFruits)
	r.GET("/fruits/:id", getFruitByID)
	r.POST("/fruits", addFruit)
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
