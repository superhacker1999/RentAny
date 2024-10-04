package web

import (
	"RentAny/model/web/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Пример защищённого эндпоинта
func ProtectedEndpoint(c *gin.Context) {
	username, _ := c.Get("email")
	c.JSON(http.StatusOK, gin.H{"message": "Hello, " + username.(string)})
}

func initEndpoints(r *gin.Engine) {
	// Маршрут для логина (получение JWT)
	r.POST("/login", handlers.Login)

	// Защищённый маршрут
	r.GET("/protected/greeting", handlers.ValidateJWT, ProtectedEndpoint)

	r.POST("/sign-up", handlers.Signup)
}

func Run() {
	r := gin.Default()

	initEndpoints(r)

	r.Run(":8080")
}
