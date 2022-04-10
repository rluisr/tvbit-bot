package external

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rluisr/tvbit-bot/pkg/adapter/controllers"
)

var Router *gin.Engine

func init() {
	logger := &Logger{}

	Router = gin.New()
	Router.ForwardedByClientIP = true

	tvController := controllers.NewTVController(logger)

	Router.POST("/tv", func(c *gin.Context) { tvController.Handle(c) })

	var addr string
	if os.Getenv("SERVER_ENV") == "local" {
		Router.Use(gin.Logger())
		addr = ":3001"
	} else {
		addr = ":8082"
	}
	if os.Getenv("NOMAD_PORT_api") != "" {
		addr = ":" + os.Getenv("NOMAD_PORT_api")
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	log.Println("Close connections ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	defer cancel()
	log.Println("Server exiting")
}
