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

func (u *RegisterUserUsecase) RegisterUser(req model.UserReqForHTTPPost, firebaseUID string) (*model.UserResForHTTPPost, error) {
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
	userEntity, err := u.UserDao.RegisterUser(newIDStr, firebaseUID, req.Name)
	if err != nil {
		return nil, err
	}

	return &model.UserResForHTTPPost{
		Id:          userEntity.Id,
		Name:        userEntity.Name,
		FirebaseUID: userEntity.FirebaseUID,
	}, nil
}
