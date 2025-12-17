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

// SearchProduct は商品を検索します（キーワードが空なら全件）(新着順)
func (u *ProductSearchUsecase) SearchProduct(keyword string) ([]*model.Product, error) {
	var products []*model.Product
	var err error
	// 1. DBから取得
	if keyword == "" {
		products, err = u.ProductDAO.FindAll()
	} else {
		products, err = u.ProductDAO.FindByName(keyword)
	}

	return u.processProducts(products, err)
}

func (u *ProductSearchUsecase) GetSellingProducts(firebaseUID string) ([]*model.Product, error) {
	user, err := u.UserDAO.FindByFirebaseUID(firebaseUID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	return u.processProducts(u.ProductDAO.FindByUserID(user.ID))
}

func (u *ProductSearchUsecase) GetPurchasedProducts(firebaseUID string) ([]*model.Product, error) {
	user, err := u.UserDAO.FindByFirebaseUID(firebaseUID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	return u.processProducts(u.ProductDAO.FindByBuyerID(user.ID))
}

func (u *ProductSearchUsecase) GetLikedProducts(firebaseUID string) ([]*model.Product, error) {
	user, err := u.UserDAO.FindByFirebaseUID(firebaseUID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	return u.processProducts(u.ProductDAO.FindLikedProducts(user.ID))
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
