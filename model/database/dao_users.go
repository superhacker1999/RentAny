package dao

import (
	"github.com/jmoiron/sqlx"
)

// User представляет структуру пользователя с тегами для sqlx
type User struct {
	ID           int    `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	Name         string `db:"name"`
	Surname      string `db:"surname"`
	PhoneNumber  string `db:"phone_number"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}

// UserDAO предоставляет методы для работы с пользователями
type UserDAO struct {
	db *sqlx.DB
}

// NewUserDAO создает новый экземпляр UserDAO
func newUserDAO(db *sqlx.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) FindByEmail(email string) (*User, error) {
	user := &User{}

	query := `SELECT * FROM users WHERE email = $1`

	err := dao.db.QueryRowx(query, email).StructScan(user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *UserDAO) FindByPhone(phone string) (*User, error) {
	user := &User{}

	query := `SELECT * FROM users WHERE phone_number = $1`

	err := dao.db.QueryRowx(query, phone).StructScan(user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// Create создает нового пользователя
func (dao *UserDAO) Create(user *User) error {
	query := `INSERT INTO Users (email, password_hash, name, surname, phone_number) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return dao.db.QueryRowx(query, user.Email, user.PasswordHash, user.Name, user.Surname, user.PhoneNumber).Scan(&user.ID)
}

func (dao *UserDAO) GetByID(id int) (*User, error) {
	user := &User{}
	query := `SELECT * FROM Users WHERE id = $1`
	err := dao.db.Get(user, query, id) // Используем sqlx.Get
	if err != nil {
		return nil, err
	}
	return user, nil
}
