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
func (u *MessageUsecase) SendMessage(senderFirebaseUID, receiverID, content, productID string) (*model.Message, error) {
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
		ProductID:  productID,
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

// GetChatList: チャット一覧（相手ごとの最新メッセージ）を取得
func (u *MessageUsecase) GetChatList(myFirebaseUID string) ([]*model.ChatListRes, error) {
	// 1. 自分を特定
	me, err := u.UserDAO.FindByFirebaseUID(myFirebaseUID)
	if err != nil {
		return nil, err
	}
	if me == nil {
		return nil, errors.New("user not found")
	}

	// 2. 全メッセージを取得（新しい順）
	allMessages, err := u.MessageDAO.FindAllByUserID(me.ID)
	if err != nil {
		return nil, err
	}

	//  未読数の集計マップを作る
	unreadCounts := make(map[string]int)
	for _, msg := range allMessages {
		// 自分が受信者 かつ まだ読んでいない(false)場合
		if msg.ReceiverID == me.ID && !msg.IsRead {
			unreadCounts[msg.SenderID]++
		}
	}

	// 3. 相手ごとに集約する
	var chatList []*model.ChatListRes
	processedPartners := make(map[string]bool) // 既に処理した相手IDを記録

	for _, msg := range allMessages {
		// 相手のIDを特定
		partnerID := msg.ReceiverID
		if msg.SenderID != me.ID {
			partnerID = msg.SenderID // 自分が受信者の場合、相手はSender
		}

		// まだリストに追加していない相手なら追加
		if !processedPartners[partnerID] {
			// 相手の情報を取得（名前を知るため）
			partner, err := u.UserDAO.FindByID(partnerID)
			if err != nil {
				// エラーでも一旦スキップして続ける
				continue
			}
			partnerName := "不明なユーザー"
			if partner != nil {
				partnerName = partner.Name
			}

			chatList = append(chatList, &model.ChatListRes{
				PartnerID:   partnerID,
				PartnerName: partnerName,
				LastMessage: msg.Content,
				LastTime:    msg.CreatedAt,
				UnreadCount: unreadCounts[partnerID],
			})

			processedPartners[partnerID] = true
		}
	}

	return chatList, nil
}

// 既読にする処理
func (u *MessageUsecase) MarkAsRead(myFirebaseUID, partnerID string) error {
	me, err := u.UserDAO.FindByFirebaseUID(myFirebaseUID)
	if err != nil || me == nil {
		return errors.New("user not found")
	}
	return u.MessageDAO.MarkAsRead(me.ID, partnerID)
}
