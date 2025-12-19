package usecase

import (
	"errors"
	"hackathon-backend/dao"
	"hackathon-backend/model"
	"hackathon-backend/service"
)

type ProductSearchUsecase struct {
	ProductDAO     *dao.ProductDao
	UserDAO        *dao.UserDao
	StorageService *service.StorageService
}

func NewProductSearchUsecase(pDAO *dao.ProductDao, uDAO *dao.UserDao, sService *service.StorageService) *ProductSearchUsecase {
	return &ProductSearchUsecase{
		ProductDAO:     pDAO,
		UserDAO:        uDAO,
		StorageService: sService,
	}
}

// 内部用ヘルパー: FirebaseUID から 内部UserID を取得する (未ログインなら空文字)
func (u *ProductSearchUsecase) getInternalUserID(firebaseUID string) string {
	if firebaseUID == "" {
		return ""
	}
	user, err := u.UserDAO.FindByFirebaseUID(firebaseUID)
	if err != nil || user == nil {
		return ""
	}
	return user.ID
}

// SearchProduct: 商品検索
func (u *ProductSearchUsecase) SearchProduct(keyword, sortOrder, status, viewerFirebaseUID string, page, limit int) (*model.ProductPage, error) {
	currentUserID := u.getInternalUserID(viewerFirebaseUID)

	// ページ番号の補正
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	// 1. データ取得
	products, err := u.ProductDAO.Search(keyword, sortOrder, status, currentUserID, "", limit, offset)
	if err != nil {
		return nil, err
	}
	// 画像URL処理
	products, _ = u.processProducts(products, nil)

	// 2. 件数取得
	total, err := u.ProductDAO.SearchCount(keyword, status, "")
	if err != nil {
		return nil, err
	}

	return &model.ProductPage{Products: products, Total: total}, nil
}

// GetProductsByUserID: 特定のユーザーの商品一覧
func (u *ProductSearchUsecase) GetProductsByUserID(targetUserID, sortOrder, status, viewerFirebaseUID string, page, limit int) (*model.ProductPage, error) {
	currentUserID := u.getInternalUserID(viewerFirebaseUID)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	products, err := u.ProductDAO.Search("", sortOrder, status, currentUserID, targetUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	products, _ = u.processProducts(products, nil)

	total, err := u.ProductDAO.SearchCount("", status, targetUserID)
	if err != nil {
		return nil, err
	}

	return &model.ProductPage{Products: products, Total: total}, nil
}

// GetSellingProducts: 出品している商品 (targetFirebaseUID: 出品者, viewerFirebaseUID: 見ている人)
func (u *ProductSearchUsecase) GetSellingProducts(targetFirebaseUID string, viewerFirebaseUID string) ([]*model.Product, error) {
	// 出品者を特定
	targetUser, err := u.UserDAO.FindByFirebaseUID(targetFirebaseUID)
	if err != nil || targetUser == nil {
		return nil, errors.New("user not found")
	}

	// 見ている人を特定
	currentUserID := u.getInternalUserID(viewerFirebaseUID)

	return u.processProducts(u.ProductDAO.FindByUserID(targetUser.ID, currentUserID))
}

// GetPurchasedProducts: 購入した商品 (targetFirebaseUID: 購入者, viewerFirebaseUID: 見ている人)
func (u *ProductSearchUsecase) GetPurchasedProducts(targetFirebaseUID string, viewerFirebaseUID string) ([]*model.Product, error) {
	targetUser, err := u.UserDAO.FindByFirebaseUID(targetFirebaseUID)
	if err != nil || targetUser == nil {
		return nil, errors.New("user not found")
	}

	currentUserID := u.getInternalUserID(viewerFirebaseUID)

	return u.processProducts(u.ProductDAO.FindByBuyerID(targetUser.ID, currentUserID))
}

// GetLikedProducts: いいねした商品 (targetFirebaseUID: いいねした人, viewerFirebaseUID: 見ている人)
func (u *ProductSearchUsecase) GetLikedProducts(targetFirebaseUID string, viewerFirebaseUID string) ([]*model.Product, error) {
	targetUser, err := u.UserDAO.FindByFirebaseUID(targetFirebaseUID)
	if err != nil || targetUser == nil {
		return nil, errors.New("user not found")
	}

	currentUserID := u.getInternalUserID(viewerFirebaseUID)

	return u.processProducts(u.ProductDAO.FindLikedProducts(targetUser.ID, currentUserID))
}

// 共通処理: DBから取った商品の画像URLを変換して返す
func (u *ProductSearchUsecase) processProducts(products []*model.Product, err error) ([]*model.Product, error) {
	if err != nil {
		return nil, err
	}
	for _, p := range products {
		if p.ImageURL != "" {
			signedURL, err := u.StorageService.GenerateSignedURL(p.ImageURL)
			if err == nil {
				p.ImageURL = signedURL
			}
		}
	}
	return products, nil
}
