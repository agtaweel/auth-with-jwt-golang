package routes

import (
	controller "github.com/agtaweel/golang-jwt-project/controllers"
	"github.com/agtaweel/golang-jwt-project/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRouters(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authnticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/user/:user_id", controller.GetUser())
	incomingRoutes.POST("/add-product", controller.AddProduct())
	incomingRoutes.GET("/products", controller.IndexProducts())
}
