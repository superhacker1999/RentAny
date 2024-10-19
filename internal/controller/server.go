package server

import (
	"RentAny/internal/controller/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func initEndpoints(r *gin.Engine) {
	userAccessHandler, err := handlers.NewUserAccessHandler()

	if err != nil {
		log.Fatal(err)
	}

	userHandler, err := handlers.NewUserHandler()

	if err != nil {
		log.Fatal(err)
	}

	r.POST("/login", userAccessHandler.Login)
	r.POST("/sign-up", userAccessHandler.Signup)

	authGroup := r.Group("/auth")
	{
		authGroup.Use(userAccessHandler.AuthorizationMiddleware)
		authGroup.GET("/get-user-by-id", userHandler.GetUserByID)
		authGroup.GET("/get-user-by-item-id", userHandler.GetUserByID)

		// add new protected endpoints here
	}

}

func Run() {
	r := gin.Default()

	initEndpoints(r)

	r.Run(":8081")
}
