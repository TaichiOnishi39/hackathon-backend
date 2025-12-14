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

func (u *RegisterUserUsecase) RegisterUser(req model.UserReqForHTTPPost) (string, error) {
	// バリデーションロジック
	if err := req.Validate(); err != nil {
		return "", err
	}
	// ID生成ロジック
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	newID := ulid.MustNew(ulid.Timestamp(t), entropy)
	newIDStr := newID.String()

	// DAOを呼び出して保存
	err := u.UserDao.Insert(newIDStr, req.Name, req.Age)
	if err != nil {
		return "", err
	}

	return newIDStr, nil
}
