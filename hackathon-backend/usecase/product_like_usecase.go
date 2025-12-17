package usecase

import (
	"errors"
	"hackathon-backend/dao"
)

type ProductLikeUsecase struct {
	LikeDAO *dao.LikeDao
	UserDAO *dao.UserDao
}

func NewProductLikeUsecase(lDAO *dao.LikeDao, uDAO *dao.UserDao) *ProductLikeUsecase {
	return &ProductLikeUsecase{
		LikeDAO: lDAO,
		UserDAO: uDAO,
	}
}

// ToggleLike: いいねを切り替える（していなければ追加、していれば解除）
// 戻り値 bool: 最終的に「いいね状態(true)」になったか「解除(false)」されたか
func (u *ProductLikeUsecase) ToggleLike(productID, firebaseUID string) (bool, error) {
	// 1. ユーザー特定
	user, err := u.UserDAO.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("user not found")
	}

	// 2. 現在の状態を確認
	hasLiked, err := u.LikeDAO.HasLiked(user.ID, productID)
	if err != nil {
		return false, err
	}

	// 3. 切り替え実行
	if hasLiked {
		// 解除
		if err := u.LikeDAO.RemoveLike(user.ID, productID); err != nil {
			return false, err
		}
		return false, nil // 結果: OFF
	} else {
		// 追加
		if err := u.LikeDAO.AddLike(user.ID, productID); err != nil {
			return false, err
		}
		return true, nil // 結果: ON
	}
}

// GetLikeStatus: 現在の状態を確認する
func (u *ProductLikeUsecase) GetLikeStatus(productID, firebaseUID string) (bool, error) {
	user, err := u.UserDAO.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil // ユーザーがいなければ false
	}
	return u.LikeDAO.HasLiked(user.ID, productID)
}
