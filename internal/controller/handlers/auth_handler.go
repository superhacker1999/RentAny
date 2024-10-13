package handlers

import (
	"RentAny/internal/controller/utils"
	"RentAny/internal/services"
	"RentAny/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strings"
)

type UserAccessHandler struct {
	validate    *validator.Validate
	authService *services.AuthService
}

func NewUserAccessHandler() (*UserAccessHandler, error) {
	userAccessManager := &UserAccessHandler{}
	var err error

	userAccessManager.authService, err = services.NewAuthService()

	if err != nil {
		return nil, err
	}

	userAccessManager.validate = validator.New()
	userAccessManager.validate.RegisterValidation("pass-validation", utils.ValidatePassword)
	userAccessManager.validate.RegisterValidation("phone-validation", utils.ValidatePhoneNumber)

	userAccessManager.validate.RegisterStructValidation(utils.ValidateLoginCredentials, types.LoginCredentials{})

	return userAccessManager, nil
}

// Validates JWT token sent by user
func (uah *UserAccessHandler) AuthorizationMiddleware(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")

	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		c.Abort()
		return
	}

	userID, status, err := uah.authService.Authenticate(tokenStr)

	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Set("user_id", userID)

	c.Next()
}

func (uah *UserAccessHandler) Login(c *gin.Context) {
	var loginCreds types.LoginCredentials

	if err := c.ShouldBindJSON(&loginCreds); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		log.Println(err)
		return
	}

	if err := uah.validate.Struct(loginCreds); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMessages []string

		for _, validationErr := range validationErrors {
			errorMessages = append(errorMessages, validationErr.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(errorMessages, ", ")})
		return
	}

	token, status, err := uah.authService.Login(loginCreds)

	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(status, gin.H{"token": token})
}

func (uah *UserAccessHandler) Signup(c *gin.Context) {
	var signupCreds types.SignupCredentials

	if err := c.ShouldBindJSON(&signupCreds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := uah.validate.Struct(signupCreds); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMessages []string

		// Собираем все ошибки в массив
		for _, err := range validationErrors {
			errorMessages = append(errorMessages, err.Error())
		}

		// Возвращаем все ошибки сразу
		c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
		return
	}

	status, err := uah.authService.Signup(signupCreds)

	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, gin.H{"result": "success"})
}
