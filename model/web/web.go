package web

import (
	"RentAny/model/web/handlers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Пример защищённого эндпоинта
func ProtectedEndpoint(c *gin.Context) {
	username, _ := c.Get("email")
	c.JSON(http.StatusOK, gin.H{"message": "Hello, " + username.(string)})
}

func initEndpoints(r *gin.Engine) {
	userAccessManager, err := handlers.NewUserAccessManager()

	if err != nil {
		log.Fatal(err)
	}

	// Маршрут для логина (получение JWT)
	r.POST("/login", userAccessManager.Login)

	// Защищённый маршрут
	r.GET("/protected/greeting", userAccessManager.ValidateJWT, ProtectedEndpoint)

	r.POST("/sign-up", userAccessManager.Signup)
}

func Run() {
	r := gin.Default()

	initEndpoints(r)

	r.Run(":8080")
}
