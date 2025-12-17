package usecase

import (
	"errors"
	"hackathon-backend/dao"
)

type ProductPurchaseUsecase struct {
	ProductDAO *dao.ProductDao
	UserDAO    *dao.UserDao
}

func NewProductPurchaseUsecase(pDAO *dao.ProductDao, uDAO *dao.UserDao) *ProductPurchaseUsecase {
	return &ProductPurchaseUsecase{ProductDAO: pDAO, UserDAO: uDAO}
}

func (u *ProductPurchaseUsecase) PurchaseProduct(productID, firebaseUID string) error {
	user, err := u.UserDAO.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	return u.ProductDAO.UpdateBuyerID(productID, user.ID)
}
