package usecase

import (
	"errors"
	"hackathon-backend/dao"
)

type ProductDeleteUsecase struct {
	ProductDAO *dao.ProductDao
	UserDAO    *dao.UserDao
}

func NewProductDeleteUsecase(pDAO *dao.ProductDao, uDAO *dao.UserDao) *ProductDeleteUsecase {
	return &ProductDeleteUsecase{
		ProductDAO: pDAO,
		UserDAO:    uDAO,
	}
}

// DeleteProduct は商品を削除します
func (u *ProductDeleteUsecase) DeleteProduct(productID, firebaseUID string) error {
	// 1. Firebase UID から User ULID を特定
	user, err := u.UserDAO.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 2. 商品を削除（自分のものかどうかのチェックはDAOで行われる）
	return u.ProductDAO.Delete(productID, user.ID)
}
