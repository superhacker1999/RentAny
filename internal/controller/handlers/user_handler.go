package handlers

import (
	"RentAny/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type UserHandler struct {
	validate    *validator.Validate
	userService *services.UserService
}

func NewUserHandler() (*UserHandler, error) {
	userHandler := &UserHandler{}

	userHandler.validate = validator.New()
	// TODO add future validations here

	var err error

	userHandler.userService, err = services.NewUserService()

	if err != nil {
		return nil, err
	}

	return userHandler, nil
}

func (u *UserHandler) GetUserByID(c *gin.Context) {
	type userRequest struct {
		ID int `json:"id" binding:"required"`
	}

	ur := userRequest{}

	if err := c.ShouldBindJSON(&ur); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, status, err := u.userService.GetUserByID(ur.ID)

	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(status, gin.H{"user": user})

}

func (u *UserHandler) GetUserByItem(c *gin.Context) {
	type userRequest struct {
		ItemID int `json:"item_id" binding:"required"`
	}

	ur := userRequest{}

	if err := c.ShouldBindJSON(&ur); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, status, err := u.userService.GetUserByItemID(ur.ItemID)

	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(status, gin.H{"user": user})
}
