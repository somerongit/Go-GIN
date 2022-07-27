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

	/*
	 Accessing URL Params and Query Params
	 Return empty string if query not exsits
	*/
	routerGroup.GET("/users/:id", func(ctx *gin.Context) {
		// id := ctx.Param("id")
		// Or
		id := ctx.Params.ByName("id")
		info := ""

		if ctx.Query("isAddress") == "true" {
			info = info + ",\n\tPincode: 101010,\n\tCity: Kolkata"
		}
		if ctx.Query("isAdmin") == "true" {
			info = info + ",\n\tIsAdmin: true"
		}

		ctx.String(http.StatusOK, "{\n\tId: 0"+id+",\n\tName: S.Bakuli"+info+"\n}")
	})

	// Accessing Post data and header info
	routerGroup.POST("/users", func(ctx *gin.Context) {
		if ctx.GetHeader("X-Secure-Host") == "3000" {
			name := ctx.PostForm("name")
			pincode := ctx.PostForm("pincode")
			city := ctx.PostForm("city")
			isAdmin := ctx.DefaultPostForm("isAdmin", "false")

			ctx.String(http.StatusAccepted, "{\n\tId: 01,\n\tName: "+name+",\n\tPincode: "+pincode+",\n\tCity: "+city+",\n\tIsAdmin: "+isAdmin+"\n}")
		} else {
			ctx.String(http.StatusForbidden, "You are not authorized...")
		}
	})

	log.Fatal(router.Run(port))
}
