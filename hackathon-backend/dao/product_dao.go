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

func (d *ProductDao) DeleteProduct(productID string, userID string) error {
	// user_id も条件に入れることで、他人の商品を消せないようにする
	query := "DELETE FROM products WHERE id = ? AND user_id = ?"
	result, err := d.db.Exec(query, productID, userID)
	if err != nil {
		return err
	}

	// 実際に消えたか確認（該当なし＝他人の商品 or 存在しない）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// 削除対象が見つからなかった（権限なし含む）場合はエラーを返すなどの設計も可能ですが、
		// ここではエラーなしとして返すか、専用エラーを返すか決められます。
		// 今回は「権限がないか商品がない」ことがわかるようにエラーを返しましょう。
		return sql.ErrNoRows
	}

	return nil
}
