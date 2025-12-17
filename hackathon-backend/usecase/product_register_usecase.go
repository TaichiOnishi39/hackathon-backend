package usecase

import (
	"errors"
	"io"
	"math/rand"
	"time"

	"hackathon-backend/dao"
	"hackathon-backend/model"
	"hackathon-backend/service"

	"github.com/oklog/ulid/v2"
)

type ProductRegisterUsecase struct {
	ProductDAO     *dao.ProductDao
	UserDAO        *dao.UserDao
	StorageService *service.StorageService
}

func NewProductRegisterUsecase(pDAO *dao.ProductDao, uDAO *dao.UserDao, sService *service.StorageService) *ProductRegisterUsecase {
	return &ProductRegisterUsecase{
		ProductDAO:     pDAO,
		UserDAO:        uDAO,
		StorageService: sService,
	}
}

// UpdateProduct が商品登録のメインロジックです
func (u *ProductRegisterUsecase) RegisterProduct(firebaseUID, name, description string, price int, imageFile io.Reader, imageFilename string) (*model.Product, error) {
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

	// 3. 画像アップロード (ファイルがある場合のみ)
	var storedImageName string
	if imageFile != nil && imageFilename != "" {
		// ファイル名が重複しないようにIDをプレフィックスにつける
		// 例: products/01HXYZ..._cat.jpg
		uploadPath := "products/" + productID + "_" + imageFilename

		path, err := u.StorageService.Upload(nil, imageFile, uploadPath) // contextは一旦nil
		if err != nil {
			return nil, err
		}
		storedImageName = path
	}

	//  保存用のモデルを作成する
	newProduct := &model.Product{
		ID:          productID,
		Name:        name,
		Price:       price,
		Description: description,
		UserID:      user.ID, // ここで内部ULIDを紐付け！
		ImageURL:    storedImageName,
	}

	// 4. DAO に保存を依頼する
	if err := u.ProductDAO.Create(newProduct); err != nil {
		return nil, err
	}

	return newProduct, nil
}
