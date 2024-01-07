package models

import (
	"github.com/gin-gonic/gin"
)

type Fruit struct {
	ID    string `json:"id"`
	Name  string `json:"name" validate:"alpha"`
	Color string `json:"color" validate:"alpha"`
}

type FruitHandler interface {
	GetFruits(*gin.Context)
}
