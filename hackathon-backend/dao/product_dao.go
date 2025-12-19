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

// 共通の検索条件（WHERE句とARGS）を作成するヘルパー
func (d *ProductDao) buildSearchCondition(keyword, status, targetUserID string) (string, []interface{}) {
	// u: 出品者, u2: 購入者
	query := ` FROM products p 
	           JOIN users u ON p.user_id = u.id 
	           LEFT JOIN users u2 ON p.buyer_id = u2.id 
	           WHERE 1=1 `
	var args []interface{}

	if targetUserID != "" {
		query += " AND p.user_id = ? "
		args = append(args, targetUserID)
	}
	if keyword != "" {
		query += " AND p.name LIKE ? "
		args = append(args, "%"+keyword+"%")
	}
	if status == "selling" {
		query += " AND p.buyer_id IS NULL "
	} else if status == "sold" {
		query += " AND p.buyer_id IS NOT NULL "
	}

	return query, args
}

func (d *ProductDao) Search(keyword, sortOrder, status, currentUserID, targetUserID string, limit, offset int) ([]*model.Product, error) {
	whereQuery, args := d.buildSearchCondition(keyword, status, targetUserID)

	selectQuery := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id,
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, 
			u.name, 
			COALESCE(u2.name, ''), -- ★追加: 購入者名(なければ空文字)
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count,
			EXISTS(SELECT 1 FROM likes WHERE product_id = p.id AND user_id = ?) as is_liked
	` + whereQuery

	finalArgs := append([]interface{}{currentUserID}, args...)

	switch sortOrder {
	case "price_asc":
		selectQuery += " ORDER BY p.price ASC "
	case "price_desc":
		selectQuery += " ORDER BY p.price DESC "
	case "oldest":
		selectQuery += " ORDER BY p.created_at ASC "
	case "likes":
		selectQuery += " ORDER BY like_count DESC, p.created_at DESC "
	default:
		selectQuery += " ORDER BY p.created_at DESC "
	}

	selectQuery += " LIMIT ? OFFSET ? "
	finalArgs = append(finalArgs, limit, offset)

	return d.fetchProducts(selectQuery, finalArgs...)
}

func (d *ProductDao) SearchCount(keyword, status, targetUserID string) (int, error) {
	whereQuery, args := d.buildSearchCondition(keyword, status, targetUserID)
	query := `SELECT COUNT(*) ` + whereQuery

	var count int
	err := d.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

// productIDで
func (d *ProductDao) FindByID(productID, currentUserID string) (*model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.description, p.user_id, 
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, 
			u.name, 
			COALESCE(u2.name, ''), -- ★追加
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count,
			EXISTS(SELECT 1 FROM likes WHERE product_id = p.id AND user_id = ?) as is_liked
		FROM products p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN users u2 ON p.buyer_id = u2.id -- ★追加
		WHERE p.id = ?
	`
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
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,COALESCE(u2.name, ''),
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count,
			EXISTS(SELECT 1 FROM likes WHERE product_id = p.id AND user_id = ?) as is_liked
		FROM products p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN users u2 ON p.buyer_id = u2.id
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
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,COALESCE(u2.name, ''),
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count,
			EXISTS(SELECT 1 FROM likes WHERE product_id = p.id AND user_id = ?) as is_liked
		FROM products p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN users u2 ON p.buyer_id = u2.id
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
			COALESCE(p.image_url, ''), p.created_at, p.buyer_id, u.name,COALESCE(u2.name, ''),
			(SELECT COUNT(*) FROM likes WHERE product_id = p.id) as like_count,
			EXISTS(SELECT 1 FROM likes WHERE product_id = p.id AND user_id = ?) as is_liked
		FROM products p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN users u2 ON p.buyer_id = u2.id
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

		// ★Scanに &p.BuyerName を追加
		err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.Description, &p.UserID,
			&p.ImageURL, &p.CreatedAt, &buyerID,
			&p.UserName,
			&p.BuyerName, // ★追加
			&p.LikeCount, &p.IsLiked,
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
