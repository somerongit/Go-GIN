package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type User struct {
	Id   string `json:"id"`   // `uri:"id"`  or `form:"id"` or `header:"X-Request-Id"`
	Name string `json:"name"` // `uri:"name"`  or `form:"name"` or `header:"X-Request-Name"`
}

type Customer struct {
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required,password"`
	Role          string `json:"role" binding:"required,oneof=Basic Admin"`
	StreetAddress string `json:"street_address"`
	StreetNumber  int    `json:"street_number" binding:"required_with=StreetAddress"`
}

func verifyPassword(fl validator.FieldLevel) bool {
	regx := regexp.MustCompile("\\w{8,}")
	password := fl.Field().String()
	return regx.MatchString(password)
}

func main() {

	router := gin.Default()
	port := ":3000"
	routerGroup := router.Group("/api/v1")
	staticRouterGroup := router.Group("/storage/v1")

	router.LoadHTMLGlob("./templates/*")

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/home")
	})

	router.GET("/home", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title":       "Home",
			"description": "This is home page.",
		})
	})

	router.GET("/about", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "about.html", gin.H{
			"title":       "About",
			"description": "This is about page.",
		})
	})

	// Static Storage Route
	staticRouterGroup.StaticFS("/file", http.Dir("./static"))

	if valid, ok := binding.Validator.Engine().(*validator.Validate); ok {
		valid.RegisterValidation("password", verifyPassword)
	}

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

	// Accessing Post form data and header info
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

	// Accessing Json body data with struct
	routerGroup.POST("/user", func(ctx *gin.Context) {
		var usr User
		/* 	Other options for Data binding
				ctx.ShouldBindUri(&usr)
		 		ctx.ShouldBindForm(&usr)
		 		ctx.ShouldBindHeader(&usr)
		 		ctx.ShouldBindQuery(&usr)
		*/
		if e := ctx.ShouldBindJSON(&usr); e != nil {
			ctx.String(http.StatusBadRequest, e.Error())
			return
		}
		fmt.Print(usr)
		ctx.String(http.StatusOK, "User is added...")
	})

	// Accessing Json body data with validation
	routerGroup.POST("/customers", func(ctx *gin.Context) {
		var usr Customer
		if e := ctx.ShouldBindJSON(&usr); e != nil {
			ctx.String(http.StatusBadRequest, e.Error())
			return
		}
		fmt.Print(usr)
		ctx.String(http.StatusOK, "Customer is added...")
	})
	log.Fatal(router.Run(port))
}
