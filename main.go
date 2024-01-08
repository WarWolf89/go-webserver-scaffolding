package main

import (
	// ever since slog made it into the core pkg there's no need for external loggers
	"log/slog"
	"os"
	// gin is the most used and standard web framework in go

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"csaba.almasi.per/webserver/src/pkg/fruitservice/api"
	"csaba.almasi.per/webserver/src/pkg/fruitservice/fruitstore"
)

func main() {

	gengine := gin.Default()
	validate := validator.New(validator.WithRequiredStructEnabled())
	fsvc := fruitstore.ProvideSVC()

	// need to type assert since fsvc is an interface
	// close redis connection since the client creation is called in main
	defer func() {
		if err := fsvc.(*fruitstore.RedisStore).Client.Conn().Close(); err != nil {
			slog.Error("Failed to close Redis connection:", "error", err)
		}
	}()

	api := api.ProvideApi(gengine, fsvc, validate)

	api.RegisterAPIEndpoints()
	if err := api.Gengine.Run(":8080"); err != nil {
		slog.Error("Failed setting up webserver:", err)
		os.Exit(1)
	}
}
