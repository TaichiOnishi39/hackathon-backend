package usecase

import (
	"hackathon-backend/dao"
	"hackathon-backend/model"
)

type SearchUserUsecase struct {
	UserDao *dao.UserDao
}

func NewSearchUserUsecase(d *dao.UserDao) *SearchUserUsecase {
	return &SearchUserUsecase{UserDao: d}
}

func (u *SearchUserUsecase) SearchUser(name string) ([]model.UserResForHTTPGet, error) {
	// ここで「名前が空でないか」などのロジックを入れることも可能
	return u.UserDao.FindByName(name)
}

func (u *SearchUserUsecase) GetUserByFirebaseUID(firebaseUID string) (*model.User, error) {
	return u.UserDao.GetUserByFirebaseUID(firebaseUID)
}
