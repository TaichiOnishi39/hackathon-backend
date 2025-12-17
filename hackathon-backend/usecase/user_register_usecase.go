package usecase

import (
	"math/rand"
	"time"

	"hackathon-backend/dao"
	"hackathon-backend/model"

	"github.com/oklog/ulid/v2"
)

type RegisterUserUsecase struct {
	UserDao *dao.UserDao
}

func NewRegisterUserUsecase(d *dao.UserDao) *RegisterUserUsecase {
	return &RegisterUserUsecase{UserDao: d}
}

func (u *RegisterUserUsecase) RegisterUser(req model.CreateUserReq, firebaseUID string) (*model.UserRes, error) {
	// バリデーションロジック
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// ID生成ロジック
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	newID := ulid.MustNew(ulid.Timestamp(t), entropy)
	newIDStr := newID.String()

	// DAOを呼び出して保存
	userEntity, err := u.UserDao.CreateOrUpdate(newIDStr, firebaseUID, req.Name)
	if err != nil {
		return nil, err
	}

	return &model.UserRes{
		ID:          userEntity.ID,
		Name:        userEntity.Name,
		FirebaseUID: userEntity.FirebaseUID,
	}, nil
}
