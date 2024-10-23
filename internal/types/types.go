package types

import "database/sql"

type ItemRepository struct {
	ID        int    `db:"id"`
	UserID    int    `db:"user_id"`
	Title     string `db:"title"`
	Desc      string `db:"description"`
	Price     int    `db:"price_per_hour"`
	Category  string `db:"category"`
	Available bool   `db:"available"`
	Location  string `db:"location"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type UserRepository struct {
	ID           int            `db:"id"`
	Email        string         `db:"email"`
	PasswordHash string         `db:"password_hash"`
	Name         string         `db:"name"`
	Surname      string         `db:"surname"`
	PhoneNumber  string         `db:"phone_number"`
	ProfilePic   sql.NullString `db:"profile_pic"`
	CreatedAt    string         `db:"created_at"`
	UpdatedAt    string         `db:"updated_at"`
}

type UserDTO struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	PhotoURL    string `json:"photo_url"`
	PhoneNumber string `json:"phone_number"`
}

type LoginCredentials struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Password string `json:"password" binding:"required"`
}

type SignupCredentials struct {
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
	Phone    string `json:"phone" validate:"phone-validation"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" validate:"pass-validation"`
}

func UserRepoToUserDTO(user *UserRepository) *UserDTO {
	userDTO := &UserDTO{
		Name:        user.Name,
		Surname:     user.Surname,
		PhotoURL:    user.ProfilePic.String,
		PhoneNumber: user.PhoneNumber,
	}
	return userDTO
}
