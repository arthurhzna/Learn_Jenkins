package main

import (
	"Learn_Jenkins/config"
	"Learn_Jenkins/controllers"
	"Learn_Jenkins/domain/model"
	"Learn_Jenkins/middlewares"
	"Learn_Jenkins/repositories"
	"Learn_Jenkins/routes"
	"Learn_Jenkins/services"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	db, err := config.InitDatabase()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database")

	db.AutoMigrate(&model.User{})
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)
	router := gin.Default()
	router.Use(middlewares.HandlePanic())
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Path not found"})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to User Service"})
	})

	route := routes.NewRoute(userController, router)
	route.Run()
	router.Run(":" + port)

}
