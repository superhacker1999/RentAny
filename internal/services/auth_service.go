package services

import (
	"RentAny/internal/controller/utils"
	"RentAny/internal/repository/minio"
	"RentAny/internal/repository/postgres"
	"RentAny/internal/types"
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// Структура для хранения данных JWT
type jwtClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

type AuthService struct {
	jwtKey         []byte
	connectionPool *postgres.Database
	s3Client       *minio.Database
}

func NewAuthService() (*AuthService, error) {
	authService := &AuthService{}
	var err error

	authService.connectionPool, err = postgres.GetConnectionPool()
	authService.jwtKey = []byte(os.Getenv("JWT_KEY"))

	if err != nil {
		return nil, err
	}

	authService.s3Client, err = minio.GetConnection()

	if err != nil {
		log.Println(err)
		// not getting minio connection is bad, but not critical
	}

	return authService, nil
}

func (a *AuthService) generateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &jwtClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtKey)
}

// returns userID, https status, error
func (a *AuthService) Authenticate(tokenStr string) (int, int, error) {
	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})

	if err != nil || !token.Valid {
		log.Println(err)
		return 0, http.StatusUnauthorized, errors.New("invalid token")
	}

	return claims.UserID, http.StatusOK, nil
}

func (a *AuthService) Login(creds types.LoginCredentials) (token string, status int, err error) {
	userDAO, err := a.connectionPool.GetUserDAO()

	if err != nil {
		log.Println(err)
		return "", http.StatusInternalServerError, errors.New("internal server error")
	}

	var user *types.UserRepository

	if creds.Phone != "" {
		user, err = userDAO.FindByPhone(creds.Phone)
	} else if creds.Email != "" {
		user, err = userDAO.FindByEmail(creds.Email)
	}

	if err != nil {
		log.Println(err)
		return "", http.StatusUnauthorized, errors.New("Invalid login or password")
	}

	if err := utils.CheckPassword(user.PasswordHash, creds.Password); err != nil {
		log.Println(err)
		return "", http.StatusUnauthorized, errors.New("Invalid login or password")
	}

	token, err = a.generateJWT(user.ID)

	if err != nil {
		log.Println(err)
		return "", http.StatusInternalServerError, errors.New("Failed to generate JWT")
	}

	return token, http.StatusOK, nil
}

func (a *AuthService) UploadProfilePic(userID int, file multipart.File, fileHeader *multipart.FileHeader) (status int, err error) {
	if fileHeader != nil && fileHeader != nil {
		fileHeader.Filename, err = utils.GenerateUniqueFileName(fileHeader.Filename)

		if err != nil {
			log.Println("Error generating file name : ", err.Error())
			return http.StatusInternalServerError, errors.New("Failed to upload profile pic")
		} else {
			err = a.s3Client.PutObject(file, fileHeader)
			if err != nil {
				log.Println("Error uploading profile pic : " + err.Error())
				return http.StatusInternalServerError, errors.New("Failed to upload profile pic")
			}
			user := &types.UserRepository{}
			user.ID = userID
			user.ProfilePic = sql.NullString{fileHeader.Filename, true}

			userDAO, err := a.connectionPool.GetUserDAO()

			if err != nil {
				log.Println(err)
				return http.StatusInternalServerError, errors.New("Failed to upload profile pic")
			}
			err = userDAO.AddProfilePic(user)

			if err != nil {
				log.Println(err)
				return http.StatusInternalServerError, errors.New("Failed to upload profile pic")
			}
			return http.StatusOK, nil
		}
	} else {
		return http.StatusBadRequest, errors.New("File missing")
	}
}

func (a *AuthService) Signup(creds types.SignupCredentials) (status int, err error) {
	userDAO, err := a.connectionPool.GetUserDAO()

	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, errors.New("internal server error")
	}

	var user types.UserRepository
	encryptedPassword, err := utils.HashPassword(creds.Password)

	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, errors.New("internal server error")
	}

	user.Name = creds.Name
	user.Surname = creds.Surname
	user.PhoneNumber = creds.Phone
	user.Email = creds.Email
	user.PasswordHash = encryptedPassword

	err = userDAO.Create(&user)

	// TODO add "already have user with such email or phone" handling

	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, errors.New("Failed to create user")
	} else {
		return http.StatusCreated, nil
	}
}
