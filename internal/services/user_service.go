package services

import (
	"RentAny/internal/repository/postgres"
	"RentAny/internal/types"
	"database/sql"
	"errors"
	"log"
	"net/http"
)

type UserService struct {
	connectionPool *postgres.Database
}

func NewUserService() (*UserService, error) {
	userService := &UserService{}
	var err error

	userService.connectionPool, err = postgres.GetConnectionPool()

	if err != nil {
		return nil, err
	}
	return userService, nil
}

func (u *UserService) GetUserByID(id int) (*types.UserDTO, int, error) {
	daoUsers, err := u.connectionPool.GetUserDAO()

	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError, errors.New("internal server error")
	}

	user, err := daoUsers.GetByID(id)

	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, http.StatusNotFound, errors.New("user not found")
		}
		return nil, http.StatusInternalServerError, errors.New("couldn't get user due to internal error")
	}

	return types.UserRepoToUserDTO(user), http.StatusOK, nil
}

func (u *UserService) GetUserByItemID(itemID int) (*types.UserDTO, int, error) {
	daoUsers, err := u.connectionPool.GetUserDAO()

	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError, errors.New("internal server error")
	}

	user, err := daoUsers.GetByItemID(itemID)

	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, http.StatusNotFound, errors.New("user not found")
		}
		return nil, http.StatusInternalServerError, errors.New("couldn't get user due to internal error")
	}
	return types.UserRepoToUserDTO(user), http.StatusOK, nil
}
