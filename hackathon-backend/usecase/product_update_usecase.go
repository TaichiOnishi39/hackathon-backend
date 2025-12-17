package usecase

import (
	"errors"
	"hackathon-backend/dao"
	"hackathon-backend/model"
)

type ProductUpdateUsecase struct {
	ProductDAO *dao.ProductDao
	UserDAO    *dao.UserDao
}

func NewProductUpdateUsecase(pDAO *dao.ProductDao, uDAO *dao.UserDao) *ProductUpdateUsecase {
	return &ProductUpdateUsecase{
		ProductDAO: pDAO,
		UserDAO:    uDAO,
	}
}

// UpdateProduct は商品を更新します
func (u *ProductUpdateUsecase) UpdateProduct(productID, firebaseUID, name, description string, price int) (*model.Product, error) {
	// 1. Firebase UID から User ULID を特定
	user, err := u.UserDAO.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 2. 更新実行
	err = u.ProductDAO.Update(productID, user.ID, name, price, description)
	if err != nil {
		return nil, err
	}

	// 3. 更新後のデータを返却したい場合は再取得するか、入力値をそのまま返す
	// ここではシンプルに入力値を元にモデルを返します（IDなどはそのまま）
	return &model.Product{
		ID:          productID,
		Name:        name,
		Price:       price,
		Description: description,
		UserID:      user.ID,
		UserName:    user.Name,
	}, nil
}
