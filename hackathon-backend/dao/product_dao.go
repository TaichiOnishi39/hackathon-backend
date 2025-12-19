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

func (d *ProductDao) Search(keyword string, sortOrder string, status string, currentUserID string, targetUserID string) ([]*model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id,
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count,
			EXISTS(SELECT 1 FROM likes WHERE product_id = p.id AND user_id = ?) as is_liked
		FROM products p
		JOIN users u ON p.user_id = u.id
		WHERE 1=1 
	`

	var args []interface{}
	args = append(args, currentUserID)

	// ★追加: 特定ユーザーの商品に絞る場合
	if targetUserID != "" {
		query += " AND p.user_id = ? "
		args = append(args, targetUserID)
	}

	// キーワード検索
	if keyword != "" {
		query += " AND p.name LIKE ? "
		args = append(args, "%"+keyword+"%")
	}

	// ステータス絞り込み (Selling / Sold)
	if status == "selling" {
		query += " AND p.buyer_id IS NULL "
	} else if status == "sold" {
		query += " AND p.buyer_id IS NOT NULL "
	}

	// ソート順
	switch sortOrder {
	case "price_asc":
		query += " ORDER BY p.price ASC "
	case "price_desc":
		query += " ORDER BY p.price DESC "
	case "oldest":
		query += " ORDER BY p.created_at ASC "
	case "likes":
		query += " ORDER BY like_count DESC, p.created_at DESC "
	default:
		query += " ORDER BY p.created_at DESC "
	}

	return d.fetchProducts(query, args...)
}

// productIDで
func (d *ProductDao) FindByID(productID, currentUserID string) (*model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id, 
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count,
			EXISTS(SELECT 1 FROM likes WHERE product_id = p.id AND user_id = ?) as is_liked
		FROM products p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`
	var products []*model.Product
	// fetchProductsを再利用（1件だけどリストとして取得）
	products, err := d.fetchProducts(query, currentUserID, productID)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, sql.ErrNoRows
	}
	return products[0], nil
}

// FindByUserID: 特定のユーザーが出品した商品
func (d *ProductDao) FindByUserID(targetUserID, currentUserID string) ([]*model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id, 
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count,
			EXISTS(SELECT 1 FROM likes WHERE product_id = p.id AND user_id = ?) as is_liked
		FROM products p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
	`
	return d.fetchProducts(query, currentUserID, targetUserID)
}

// FindByBuyerID:特定のユーザーが購入した商品
func (d *ProductDao) FindByBuyerID(targetBuyerID, currentUserID string) ([]*model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id, 
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count,
			EXISTS(SELECT 1 FROM likes WHERE product_id = p.id AND user_id = ?) as is_liked
		FROM products p
		JOIN users u ON p.user_id = u.id
		WHERE p.buyer_id = ?
		ORDER BY p.created_at DESC
	`
	return d.fetchProducts(query, currentUserID, targetBuyerID)
}

// FindLikedProducts: 特定のユーザーがいいねした商品
func (d *ProductDao) FindLikedProducts(targetUserID, currentUserID string) ([]*model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id, 
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count,
			EXISTS(SELECT 1 FROM likes WHERE product_id = p.id AND user_id = ?) as is_liked
		FROM products p
		JOIN users u ON p.user_id = u.id
		JOIN likes l ON p.id = l.product_id
		WHERE l.user_id = ?
		ORDER BY p.created_at DESC
	`
	return d.fetchProducts(query, currentUserID, targetUserID)
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
			&p.ImageURL, &p.CreatedAt, &buyerID, &p.UserName, &p.LikeCount, &p.IsLiked,
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
		// 本当に存在しない（または権限がない）のか、値が変わらなかっただけなのかを確認
		var exists int
		checkQuery := "SELECT 1 FROM products WHERE id = ? AND user_id = ?"
		err := d.db.QueryRow(checkQuery, productID, userID).Scan(&exists)

		if err == sql.ErrNoRows {
			return sql.ErrNoRows // 本当になかった
		} else if err != nil {
			return err // その他のエラー
		}
		// エラーがなければ「商品はあるけど変更なし」なので成功(nil)を返す
		return nil
	}

	return nil
}
