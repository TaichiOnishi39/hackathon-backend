package usecase

import (
	"errors"
	"math/rand"
	"time"

	"hackathon-backend/dao"
	"hackathon-backend/model"

	"github.com/oklog/ulid/v2"
)

type ProductRegisterUsecase struct {
	ProductDAO *dao.ProductDao
	UserDAO    *dao.UserDao
}

func NewProductRegisterUsecase(pDAO *dao.ProductDao, uDAO *dao.UserDao) *ProductRegisterUsecase {
	return &ProductRegisterUsecase{
		ProductDAO: pDAO,
		UserDAO:    uDAO,
	}
}

// Execute が商品登録のメインロジックです
func (u *ProductRegisterUsecase) RegisterProduct(firebaseUID, name, description string, price int) (*model.Product, error) {
	// 1. Firebase UID から内部の User ULID を検索する
	// ※UserDAO に FindByFirebaseUID(uid string) (*model.User, error) がある前提です
	user, err := u.UserDAO.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, errors.New("ユーザーが見つかりませんでした")
	}

	// 2. 商品の ULID を生成する
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	productID := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	// 3. 保存用のモデルを作成する
	newProduct := &model.Product{
		ID:          productID,
		Name:        name,
		Price:       price,
		Description: description,
		UserID:      user.ID, // ここで内部ULIDを紐付け！
	}

	// 4. DAO に保存を依頼する
	if err := u.ProductDAO.Create(newProduct); err != nil {
		return nil, err
	}

	return newProduct, nil
}
