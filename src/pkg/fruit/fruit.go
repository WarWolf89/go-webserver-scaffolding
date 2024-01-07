package fruit

import (
	"context"

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

type Service interface {
	GetFruits(ctx context.Context) ([]*Fruit, error)
	AddFruit(ctx context.Context, fruit *Fruit) (string, error)
	GetFruitByID(ctx context.Context, id string) (*Fruit, error)
}
