package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	port := ":3000"

	routerGroup := router.Group("/api/v1")

	routerGroup.Handle(http.MethodGet, "/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	log.Fatal(router.Run(port))
}
