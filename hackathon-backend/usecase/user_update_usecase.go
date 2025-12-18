package usecase

import (
	"errors"
	"hackathon-backend/dao"
	"hackathon-backend/model"
)

type UserUpdateUsecase struct {
	UserDAO *dao.UserDao
}

func NewUserUpdateUsecase(uDAO *dao.UserDao) *UserUpdateUsecase {
	return &UserUpdateUsecase{UserDAO: uDAO}
}

func (u *UserUpdateUsecase) UpdateUser(firebaseUID, name, bio string) (*model.User, error) {
	// 1. ユーザー特定
	user, err := u.UserDAO.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 2. 値を書き換え
	user.Name = name
	user.Bio = bio

	// 3. 保存
	if err := u.UserDAO.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
