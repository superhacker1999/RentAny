package postgres

import (
	"RentAny/internal/types"
	"github.com/jmoiron/sqlx"
)

type ItemsDAO struct {
	db *sqlx.DB
}

func NewItemsDAO(db *sqlx.DB) *ItemsDAO {
	return &ItemsDAO{db: db}
}

func (dao *ItemsDAO) Create(item *types.ItemRepository) error {
	query := `INSERT INTO Items (user_id, title, description, price_per_hour, category, available, location)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	return dao.db.QueryRowx(query, item.UserID, item.Title, item.Desc, item.Price, item.Category, item.Available, item.Location).Scan(&item.ID)
}

func (dao *ItemsDAO) Delete(id int) error {
	query := `DELETE FROM Items WHERE id = $1`
	_, err := dao.db.Exec(query, id)
	return err
}

func (dao *ItemsDAO) Get(id int) (*types.ItemRepository, error) {
	user := &types.ItemRepository{}
	query := `SELECT * FROM Items WHERE id = $1`

	err := dao.db.Get(user, query, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *ItemsDAO) getAllItemsByUserID(userID, offset, limit int) ([]*types.ItemRepository, error) {
	query := `SELECT * FROM Items WHERE user_id = $1
			  LIMIT $2 OFFSET $3`
	rows, err := dao.db.Queryx(query, userID, offset, limit)
	if err != nil {
		return nil, err
	}
	items := []*types.ItemRepository{}
	for rows.Next() {
		item := &types.ItemRepository{}
		if err := rows.StructScan(item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (dao *ItemsDAO) Update(id int, item *types.ItemRepository) error {
	query := `UPDATE Items 
			  SET user_id = $1, title = $2, description = $3, price_per_hour = $4, 
			  category = $5, available = $6, location = $7, updated_at = NOW()
			  WHERE id = $8`

	_, err := dao.db.Exec(query, item.UserID, item.Title, item.Desc, item.Price, item.Category, item.Available, item.Location, id)
	return err
}
