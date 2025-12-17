package usecase

import (
	"hackathon-backend/dao"
	"hackathon-backend/model"
	"hackathon-backend/service"
)

type ProductDetailUsecase struct {
	ProductDAO     *dao.ProductDao
	StorageService *service.StorageService
}

func NewProductDetailUsecase(pDAO *dao.ProductDao, sService *service.StorageService) *ProductDetailUsecase {
	return &ProductDetailUsecase{ProductDAO: pDAO, StorageService: sService}
}

func (u *ProductDetailUsecase) GetProductByID(id string) (*model.Product, error) {
	product, err := u.ProductDAO.FindByID(id)
	if err != nil {
		return nil, err
	}
	// 画像URLを変換
	if product.ImageURL != "" {
		url, err := u.StorageService.GenerateSignedURL(product.ImageURL)
		if err == nil {
			product.ImageURL = url
		}
	}
	return product, nil
}
