package handlers

import (
	"RentAny/model/database"
	"RentAny/model/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type UserAccessManager struct {
	validate       *validator.Validate
	jwtKey         []byte
	connectionPool *database.Database
}

func NewUserAccessManager() (*UserAccessManager, error) {
	userAccessManager := &UserAccessManager{}
	var err error

	userAccessManager.connectionPool, err = database.GetConnectionPool()

	if err != nil {
		return nil, err
	}
	userAccessManager.jwtKey = []byte(os.Getenv("JWT_KEY"))

	userAccessManager.validate = validator.New()
	userAccessManager.validate.RegisterValidation("pass-validation", utils.ValidatePassword)
	userAccessManager.validate.RegisterValidation("phone-validation", utils.ValidatePhoneNumber)

	userAccessManager.validate.RegisterStructValidation(utils.ValidateLoginCredentials, utils.LoginCredentials{})

	return userAccessManager, nil
}

// Структура для хранения данных JWT
type claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Функция для генерации JWT токена
func (uam *UserAccessManager) generateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uam.jwtKey)
}

// Функция для валидации JWT токена
func (uam *UserAccessManager) ValidateJWT(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")

	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		c.Abort()
		return
	}

	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return uam.jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("email", claims.Email)
	c.Next()
}

// Функция для логина (создание JWT токена)
func (uam *UserAccessManager) Login(c *gin.Context) {
	var loginCreds utils.LoginCredentials

	if err := c.ShouldBindJSON(&loginCreds); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		log.Println(err)
		return
	}

	if err := uam.validate.Struct(loginCreds); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMessages []string

		for _, validationErr := range validationErrors {
			errorMessages = append(errorMessages, validationErr.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(errorMessages, ", ")})
		return
	}

	db, err := database.GetConnectionPool()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		log.Println(err)
		return
	}

	userDAO, err := db.GetUserDAO()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		log.Println(err)
		return
	}

	var user *database.User

	if loginCreds.Phone != "" {
		user, err = userDAO.FindByPhone(loginCreds.Phone)
	} else if loginCreds.Email != "" {
		user, err = userDAO.FindByEmail(loginCreds.Email)
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		return
	}

	if err := utils.CheckPassword(user.PasswordHash, loginCreds.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		return
	}

	token, err := uam.generateJWT(loginCreds.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (uam *UserAccessManager) Signup(c *gin.Context) {
	db, err := database.GetConnectionPool()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		log.Println(err)
	}

	userDAO, err := db.GetUserDAO()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		log.Println(err)
	}

	var signupCreds utils.SignupCredentials

	if err := c.ShouldBindJSON(&signupCreds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := uam.validate.Struct(signupCreds); err != nil {
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

	var user database.User
	encryptedPassword, err := utils.HashPassword(signupCreds.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		log.Println(err)
		return
	}

	user.Name = signupCreds.Name
	user.Surname = signupCreds.Surname
	user.PhoneNumber = signupCreds.Phone
	user.Email = signupCreds.Email
	user.PasswordHash = encryptedPassword

	err = userDAO.Create(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		log.Println(err)
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"You're successfully registered, ": user.Name})
	}
}
