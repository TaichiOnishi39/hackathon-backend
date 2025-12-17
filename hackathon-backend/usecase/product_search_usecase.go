package usecase

import (
	"hackathon-backend/dao"
	"hackathon-backend/model"
	"hackathon-backend/service"
)

type ProductSearchUsecase struct {
	ProductDAO     *dao.ProductDao
	StorageService *service.StorageService
}

func NewProductSearchUsecase(pDAO *dao.ProductDao, sService *service.StorageService) *ProductSearchUsecase {
	return &ProductSearchUsecase{
		ProductDAO:     pDAO,
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

	if err != nil {
		return nil, err
	}

	// 2. 画像URLの変換処理 (ファイル名 → 署名付きURL)
	for _, p := range products {
		if p.ImageURL != "" {
			// ここで署名付きURLを発行して上書きする
			signedURL, err := u.StorageService.GenerateSignedURL(p.ImageURL)
			if err == nil {
				p.ImageURL = signedURL
			}
			// エラーが出てもログに出す程度で、処理は止めない（画像なしとして扱う）
		}
	}

	return products, nil
}
