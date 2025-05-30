package main

import (
	"gin-freemarket/controllers"
	"gin-freemarket/infra"
	"gin-freemarket/middlewares"

	"gin-freemarket/repositories"
	"gin-freemarket/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine{

  itemRepository := repositories.NewItemRepository(db)
  itemService:=services.NewItemService(itemRepository)
  itemController:=controllers.NewItemController(itemService)

  authRepository := repositories.NewAuthRepository(db)
  authService := services.NewAuthService(authRepository)
  authController := controllers.NewAuthController(authService)

  r := gin.Default() //ginのdefaultルーターを初期化し、routerに格納
  r.Use(cors.Default())
  itemRouter := r.Group("/items")
  itemRouterWithAuth := r.Group("/items",middlewares.AuthMiddleware(authService))
  authRouter := r.Group("/auth")


  itemRouter.GET("/", itemController.FindAll) 
  itemRouterWithAuth.GET("/:id", itemController.FindById) 
  itemRouterWithAuth.POST("",itemController.Create)
  itemRouterWithAuth.PUT("/:id",itemController.Update)
  itemRouterWithAuth.DELETE("/:id",itemController.DeleteById)

  authRouter.POST("/signup",authController.Signup)
  authRouter.POST("/login",authController.Login)

  return r
}


func main() {
  infra.Initialize()
  db :=infra.SetupDB()
  r:= setupRouter(db)
  r.Run("localhost:8080") 
}