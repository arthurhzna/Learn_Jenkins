package routes

import (
	"Learn_Jenkins/controllers"

	"github.com/gin-gonic/gin"
)

type routeImpl struct {
	Controller controllers.UserController
	Router     *gin.Engine
}

func NewRoute(controller controllers.UserController, router *gin.Engine) UserService {
	return &routeImpl{Controller: controller, Router: router}
}

func (r *routeImpl) Run() {
	r.Router.POST("/users", r.Controller.CreateUser)
	r.Router.GET("/users/:id", r.Controller.FindUserByID)
	r.Router.GET("/users", r.Controller.FindAllUsers)

}
