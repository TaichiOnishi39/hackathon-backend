package dao

import (
	"database/sql"
	"fmt"
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
		INSERT INTO products (id, name, price, description, user_id, image_url) 
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := d.db.Exec(
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Description,
		product.UserID,
		product.ImageURL,
	)
	return err
}

func (d *ProductDao) FindAll() ([]*model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id, 
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count
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
		var buyerID sql.NullString // NULL対策
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.UserID, &p.ImageURL, &p.CreatedAt, &buyerID, &p.UserName, &p.LikeCount)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
		if buyerID.Valid {
			p.BuyerID = buyerID.String
		}
	}

	return products, nil
}

func (d *ProductDao) FindByName(keyword string) ([]*model.Product, error) {
	// ユーザー名も取得したいのでJOINします
	// WHERE p.name LIKE ? を追加
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id,
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count
		FROM products p
		JOIN users u ON p.user_id = u.id
		WHERE p.name LIKE ?
		ORDER BY p.created_at DESC
	`
	// %keyword% の形にして部分一致にする
	likeQuery := "%" + keyword + "%"

	rows, err := d.db.Query(query, likeQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p := &model.Product{}
		var buyerID sql.NullString
		err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.Description, &p.UserID, &p.ImageURL, &p.CreatedAt, &buyerID, &p.UserName, &p.LikeCount,
		)
		if err != nil {
			return nil, err
		}
		if buyerID.Valid {
			p.BuyerID = buyerID.String
		}
		products = append(products, p)
	}
	return products, nil
}

// productIDで
func (d *ProductDao) FindByID(productID string) (*model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id, 
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count
		FROM products p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`
	// buyer_id は NULL の可能性があるので sql.NullString で受け取るのが安全ですが
	// 今回はポインタか、あるいは NULL なら空文字にするなど工夫します。
	// シンプルに Scan で *string に入れると NULL 対応できます。
	var buyerID *string

	p := &model.Product{}
	err := d.db.QueryRow(query, productID).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Description,
		&p.UserID,
		&p.ImageURL,
		&p.CreatedAt,
		&buyerID, // NULL許容
		&p.UserName,
		&p.LikeCount,
	)

	if err != nil {
		return nil, err
	}

	// *string から string へ変換（NULLなら空文字）
	if buyerID != nil {
		p.BuyerID = *buyerID
	}

	return p, nil
}

// FindByUserID: 自分が【出品】した商品を取得
func (d *ProductDao) FindByUserID(userID string) ([]*model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id, 
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count
		FROM products p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
	`
	return d.fetchProducts(query, userID)
}

// FindByBuyerID: 自分が【購入】した商品を取得
func (d *ProductDao) FindByBuyerID(buyerID string) ([]*model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id, 
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count
		FROM products p
		JOIN users u ON p.user_id = u.id
		WHERE p.buyer_id = ?
		ORDER BY p.created_at DESC
	`
	return d.fetchProducts(query, buyerID)
}

// FindLikedProducts: 自分が【いいね】した商品を取得
func (d *ProductDao) FindLikedProducts(userID string) ([]*model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id, 
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count
		FROM products p
		JOIN users u ON p.user_id = u.id
		JOIN likes l ON p.id = l.product_id
		WHERE l.user_id = ?
		ORDER BY p.created_at DESC
	`
	return d.fetchProducts(query, userID)
}

// 共通処理
func (d *ProductDao) fetchProducts(query string, args ...interface{}) ([]*model.Product, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p := &model.Product{}
		var buyerID sql.NullString
		err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.Description, &p.UserID,
			&p.ImageURL, &p.CreatedAt, &buyerID, &p.UserName, &p.LikeCount,
		)
		if err != nil {
			return nil, err
		}
		if buyerID.Valid {
			p.BuyerID = buyerID.String
		}
		products = append(products, p)
	}
	return products, nil
}

// UpdateBuyerID は購入処理です（既に売れていないかチェックも含みます）
func (d *ProductDao) UpdateBuyerID(productID string, buyerID string) error {
	// buyer_id が NULL の場合のみ更新する（＝早い者勝ち）
	query := `
		UPDATE products 
		SET buyer_id = ? 
		WHERE id = ? AND buyer_id IS NULL
	`
	result, err := d.db.Exec(query, buyerID, productID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		// 更新されなかった＝既に売り切れ or 商品がない
		return fmt.Errorf("sold out or not found")
	}

	return nil
}

func (d *ProductDao) Delete(productID string, userID string) error {
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

func (d *ProductDao) Update(productID string, userID string, name string, price int, description string) error {
	// user_id も条件に入れることで、他人の商品を更新できないようにする
	query := `
		UPDATE products 
		SET name = ?, price = ?, description = ? 
		WHERE id = ? AND user_id = ?
	`
	result, err := d.db.Exec(query, name, price, description, productID, userID)
	if err != nil {
		return err
	}

	// 更新対象があったか確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // 対象なし（権限なし含む）
	}

	return nil
}
