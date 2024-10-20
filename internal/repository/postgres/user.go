package postgres

import (
	"RentAny/internal/types"
	"github.com/jmoiron/sqlx"
)

// UserDAO предоставляет методы для работы с пользователями
type UserDAO struct {
	db *sqlx.DB
}

// NewUserDAO создает новый экземпляр UserDAO
func newUserDAO(db *sqlx.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) AddProfilePic(user *types.UserRepository) error {
	query := `UPDATE users SET profile_pic = $1 WHERE id = $2;`

	_, err := dao.db.Exec(query, user.ProfilePic, user.ID)

	return err
}

func (dao *UserDAO) FindByEmail(email string) (*types.UserRepository, error) {
	user := &types.UserRepository{}

	query := `SELECT * FROM users WHERE email = $1`

	err := dao.db.QueryRowx(query, email).StructScan(user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *UserDAO) FindByPhone(phone string) (*types.UserRepository, error) {
	user := &types.UserRepository{}

	query := `SELECT * FROM users WHERE phone_number = $1`

	err := dao.db.QueryRowx(query, phone).StructScan(user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// Create создает нового пользователя
func (dao *UserDAO) Create(user *types.UserRepository) error {
	query := `INSERT INTO users (email, password_hash, name, surname, phone_number) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return dao.db.QueryRowx(query, user.Email, user.PasswordHash, user.Name, user.Surname, user.PhoneNumber).Scan(&user.ID)
}

func (dao *UserDAO) GetByID(id int) (*types.UserRepository, error) {
	user := &types.UserRepository{}
	query := `SELECT * FROM users WHERE id = $1`
	err := dao.db.Get(user, query, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *UserDAO) GetByItemID(itemID int) (*types.UserRepository, error) {
	user := &types.UserRepository{}
	query := `SELECT * FROM users JOIN items ON items.user_id = users.id WHERE items.id = $1`

	err := dao.db.Get(user, query, itemID)

	if err != nil {
		return nil, err
	}
	return user, nil
}
