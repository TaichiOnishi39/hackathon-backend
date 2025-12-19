package usecase

import (
	"context"
	"errors"
	"hackathon-backend/dao"
	"hackathon-backend/model"
	"hackathon-backend/service"
	"mime/multipart"
)

type UserUpdateUsecase struct {
	UserDAO        *dao.UserDao
	StorageService *service.StorageService
}

func NewUserUpdateUsecase(uDAO *dao.UserDao, sService *service.StorageService) *UserUpdateUsecase {
	return &UserUpdateUsecase{UserDAO: uDAO, StorageService: sService}
}

func (u *UserUpdateUsecase) UpdateUser(ctx context.Context, firebaseUID, name, bio string, imageFile *multipart.FileHeader) (*model.User, error) {
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

	if imageFile != nil {
		file, err := imageFile.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// StorageServiceを使ってアップロード
		imageURL, err := u.StorageService.Upload(ctx, file, imageFile.Filename)
		if err != nil {
			return nil, err
		}
		user.ImageURL = imageURL
	}

	// 3. 保存
	if err := u.UserDAO.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
