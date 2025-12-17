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

// GetAllProducts は全商品を新着順に取得します
func (u *ProductSearchUsecase) GetAllProducts() ([]*model.Product, error) {
	// 将来的にはここで「価格フィルタ」や「ワード検索」のロジックを追加できます
	return u.ProductDAO.FindAll()
}
