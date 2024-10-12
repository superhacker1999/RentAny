package handlers

import (
	"RentAny/internal/repository/postgres"
	"github.com/go-playground/validator/v10"
)

type ItemHandler struct {
	validate       *validator.Validate
	connectionPool *postgres.Database
}

func NewItemHandler() (*ItemHandler, error) {
	itemManager := &ItemHandler{}

	var err error

	itemManager.connectionPool, err = postgres.GetConnectionPool()

	if err != nil {
		return nil, err
	}

	itemManager.validate = validator.New()

	// TODO add validations here

	return itemManager, nil
}
