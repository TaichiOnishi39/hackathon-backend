package usecase

import (
	"hackathon-backend/dao"
	"hackathon-backend/model"
)

type ProductSearchUsecase struct {
	ProductDAO *dao.ProductDao
}

func NewProductSearchUsecase(pDAO *dao.ProductDao) *ProductSearchUsecase {
	return &ProductSearchUsecase{
		ProductDAO: pDAO,
	}
}

// SearchProducts は商品を検索します（キーワードが空なら全件）(新着順)
func (u *ProductSearchUsecase) SearchProducts(keyword string) ([]*model.Product, error) {
	if keyword == "" {
		// キーワードがないなら全件取得
		return u.ProductDAO.FindAll()
	}
	// キーワードがあるなら検索
	return u.ProductDAO.SearchByName(keyword)
}
