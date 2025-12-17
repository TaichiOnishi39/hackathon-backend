package usecase

import (
	"errors"
	"hackathon-backend/dao"
	"hackathon-backend/model"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

type MessageUsecase struct {
	MessageDAO *dao.MessageDao
	UserDAO    *dao.UserDao
}

func NewMessageUsecase(mDAO *dao.MessageDao, uDAO *dao.UserDao) *MessageUsecase {
	return &MessageUsecase{
		MessageDAO: mDAO,
		UserDAO:    uDAO,
	}
}

// SendMessage: メッセージを送信
func (u *MessageUsecase) SendMessage(senderFirebaseUID, receiverID, content string) (*model.Message, error) {
	// 1. 送信者を特定
	sender, err := u.UserDAO.FindByFirebaseUID(senderFirebaseUID)
	if err != nil {
		return nil, err
	}
	if sender == nil {
		return nil, errors.New("sender not found")
	}

	// 2. メッセージID生成
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	msgID := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	msg := &model.Message{
		ID:         msgID,
		SenderID:   sender.ID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  t,
	}

	// 3. 保存
	if err := u.MessageDAO.Create(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

// GetChatHistory: 特定の相手とのチャット履歴を取得
func (u *MessageUsecase) GetChatHistory(myFirebaseUID, otherUserID string) ([]*model.Message, error) {
	// 1. 自分を特定
	me, err := u.UserDAO.FindByFirebaseUID(myFirebaseUID)
	if err != nil {
		return nil, err
	}
	if me == nil {
		return nil, errors.New("user not found")
	}

	// 2. 履歴取得
	return u.MessageDAO.GetMessagesBetween(me.ID, otherUserID)
}
