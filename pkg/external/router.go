/*
tvbit-bot
Copyright (C) 2022  rluisr(Takuya Hasegawa)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

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

var (
	Router  *gin.Engine
	version string
)

func init() {
	logger := &Logger{}

	Router = gin.Default()
	//Router.ForwardedByClientIP = true

	tvController := controllers.NewTVController(logger)

	Router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"version": version,
			"repo":    "https://github.com/rluisr/tvbit-bot",
			"owner":   "rluisr / @rarirureluis",
		})
	})

	Router.POST("/tv", func(c *gin.Context) { tvController.Handle(c) })

	var addr string
	if os.Getenv("SERVER_ENV") == "local" {
		Router.Use(gin.Logger())
		addr = ":3001"
	} else {
		addr = ":8082"
	}
	if os.Getenv("PORT") != "" {
		addr = ":" + os.Getenv("PORT")
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
