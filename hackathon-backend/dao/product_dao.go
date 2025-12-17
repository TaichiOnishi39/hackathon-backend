package dao

import (
	"database/sql"
	"hackathon-backend/model"
)

type ProductDao struct {
	db *sql.DB
}

func NewProductDAO(db *sql.DB) *ProductDao {
	return &ProductDao{db: db}
}

func (d *ProductDao) Create(product *model.Product) error {
	query := `
		INSERT INTO products (id, name, price, description, user_id) 
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := d.db.Exec(
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Description,
		product.UserID,
	)
	return err
}

func (d *ProductDao) FindAll() ([]*model.Product, error) {
	query := `
		SELECT 
			p.id, 
			p.name, 
			p.price, 
			p.description, 
			p.user_id, 
			p.created_at,
			u.name as user_name 
		FROM products p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p := &model.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.UserID, &p.CreatedAt, &p.UserName)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
