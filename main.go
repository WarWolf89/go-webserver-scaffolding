package main

import (
	// ever since slog made it into the core pkg there's no need for external loggers
	"fmt"
	"log/slog"
	"os"

	// gin is the most used and standard web framework in go

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"csaba.almasi.per/webserver/src/pkg/fruitservice/api"
	"csaba.almasi.per/webserver/src/pkg/fruitservice/fruitstore"
	"csaba.almasi.per/webserver/src/pkg/util"
)

func main() {

	config, err := util.LoadConfig("prod_config")
	// exit immediately on non-recoverable(e.g. path issue) config load errors
	if err != nil {
		slog.Warn("Unrecoverable error when loading config file", "error", err)
		os.Exit(1)
	}

	gengine := gin.Default()
	validate := validator.New(validator.WithRequiredStructEnabled())
	fsvc := fruitstore.ProvideRedisStore(config)

	// need to type assert since fsvc is an interface
	// close redis connection since the client creation is called in main
	defer func() {
		if err := fsvc.Client.Conn().Close(); err != nil {
			slog.Error("Failed to close Redis connection:", "error", err)
		}
	}()

	api := api.ProvideApi(gengine, fsvc, validate)

	// panic so that redis connection can close on defer
	api.RegisterAPIEndpoints()
	if err := api.Gengine.Run(fmt.Sprintf(":%d", config.ServerPort)); err != nil {
		slog.Error("Failed setting up webserver:", err)
		panic(err)
	}
}
