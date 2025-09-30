package controllers

import (
	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateUser(*gin.Context)
	FindUserByID(*gin.Context)
	FindAllUsers(*gin.Context)
}
