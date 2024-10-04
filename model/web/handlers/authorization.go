package handlers

import (
	"RentAny/model/dao"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"regexp"
	"time"
)

var validate *validator.Validate

// Секретный ключ для подписи JWT
var jwtKey = []byte("my_secret_key")

// Структура для хранения данных JWT
type claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" validate:"pass-validation"`
}

// Функция для генерации JWT токена
func generateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	rules := [4]string{"([a-z])+", "([A-Z])+", "([0-9])+", "([!@#$%^&*.?-])+"}

	for _, rule := range rules {
		if !regexp.MustCompile(rule).MatchString(password) {
			return false
		}
	}

	return true
}

// Функция для валидации JWT токена
func ValidateJWT(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")

	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		c.Abort()
		return
	}

	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
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
func Login(c *gin.Context) {
	var creds credentials

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return

	}

	connStr := "postgres://postgres:1234@localhost:5432/mydb12?sslmode=disable"

	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userDAO := dao.NewUserDAO(db)

	user, err := userDAO.FindByEmail(creds.Email)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
	}

	if err := checkPassword(user.PasswordHash, creds.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := generateJWT(creds.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func checkPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func Signup(c *gin.Context) {
	connStr := "postgres://postgres:1234@localhost:5432/mydb12?sslmode=disable"

	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userDAO := dao.NewUserDAO(db)

	var creds credentials
	validate = validator.New()
	validate.RegisterValidation("pass-validation", validatePassword)

	if err := c.ShouldBindJSON(&creds); err != nil /*|| creds.Email != "user" || creds.Password != "password"*/ {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := validate.Struct(creds); err != nil {
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

	var user dao.User
	encryptedPassword, err := hashPassword(creds.Password)

	user.Email = creds.Email
	user.PasswordHash = encryptedPassword
	err = userDAO.Create(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user : " + err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"Now I know you, ": user.Email})
	}
}
