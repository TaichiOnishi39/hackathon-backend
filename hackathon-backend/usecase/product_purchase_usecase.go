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

	product, err := u.ProductDAO.FindByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}
	if product.UserID == user.ID {
		return errors.New("cannot purchase your own product")
	}
	if product.BuyerID != "" {
		return errors.New("product is already sold out")
	}
	return u.ProductDAO.UpdateBuyerID(productID, user.ID)
}
