package usecase

import (
	"hackathon-backend/dao"
	"hackathon-backend/model"
	"hackathon-backend/service"
)

type ProductDetailUsecase struct {
	ProductDAO     *dao.ProductDao
	UserDAO        *dao.UserDao
	StorageService *service.StorageService
}

func NewProductDetailUsecase(pDAO *dao.ProductDao, uDAO *dao.UserDao, sService *service.StorageService) *ProductDetailUsecase {
	return &ProductDetailUsecase{ProductDAO: pDAO, UserDAO: uDAO, StorageService: sService}
}

func (u *ProductDetailUsecase) GetProductByID(id string, viewerFirebaseUID string) (*model.Product, error) {
	// 見ている人のIDを特定
	currentUserID := ""
	if viewerFirebaseUID != "" {
		user, err := u.UserDAO.FindByFirebaseUID(viewerFirebaseUID)
		if err == nil && user != nil {
			currentUserID = user.ID
		}
	}

	// DAOに渡す
	product, err := u.ProductDAO.FindByID(id, currentUserID)
	if err != nil {
		return nil, err
	}

	// 画像URL変換
	if product.ImageURL != "" {
		url, err := u.StorageService.GenerateSignedURL(product.ImageURL)
		if err == nil {
			product.ImageURL = url
		}
	}
	return product, nil
}
