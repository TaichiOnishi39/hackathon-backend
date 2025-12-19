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

func (u *SearchUserUsecase) GetUserByFirebaseUID(firebaseUID string) (*model.User, error) {
	return u.UserDao.FindByFirebaseUID(firebaseUID)
}

func (u *SearchUserUsecase) GetUserByID(id string) (*model.User, error) {
	return u.UserDao.FindByID(id)
}
