package server

import (
	"RentAny/internal/controller/handlers"
	"RentAny/internal/repository/postgres"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Greeting(c *gin.Context) {
	db, err := postgres.GetConnectionPool()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	userDAO, err := db.GetUserDAO()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	userID := c.GetInt("user_id")

	user, err := userDAO.GetByID(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hello, " + user.Name})
}

func initEndpoints(r *gin.Engine) {
	userAccessManager, err := handlers.NewUserAccessHandler()

	if err != nil {
		log.Fatal(err)
	}

	// Маршрут для логина (получение JWT)
	r.POST("/login", userAccessManager.Login)

	authGroup := r.Group("/auth")
	{
		authGroup.Use(userAccessManager.AuthorizationMiddleware)
		authGroup.GET("/greeting", Greeting)

		// add new protected endpoints here
	}

	// Защищённый маршрут
	//r.GET("/protected/greeting", userAccessManager.ValidateJWT, ProtectedEndpoint)

	r.POST("/sign-up", userAccessManager.Signup)
}

func Run() {
	r := gin.Default()

	initEndpoints(r)

	r.Run(":8080")
}
