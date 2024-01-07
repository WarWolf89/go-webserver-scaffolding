package main

import (
	"csaba.almasi.per/webserver/src/pkg/api"
	fruitservice "csaba.almasi.per/webserver/src/pkg/svc"

	// ever since slog made it into the core pkg there's no need for external loggers

	"log/slog"
	"os"

	// validator
	"github.com/go-playground/validator/v10"
	// gin is the most used and standard web framework in go
	"github.com/gin-gonic/gin"
)

func main() {

	gengine := gin.Default()
	validate := validator.New(validator.WithRequiredStructEnabled())
	rscv := fruitservice.ProvideSVC()

	api := api.ProvideApi(gengine, rscv, validate)
	api.RegisterAPIEndpoints()

	if err := api.Gengine.Run(":8080"); err != nil {
		slog.Error("Failed setting up webserver:", err)
		os.Exit(1)
	}
}
