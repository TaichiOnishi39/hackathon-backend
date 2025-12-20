package controller

import (
	"encoding/json"
	"hackathon-backend/model"
	"hackathon-backend/usecase"
	"net/http"

	"firebase.google.com/go/auth"
)

type MessageController struct {
	BaseController
	Usecase *usecase.MessageUsecase
}

func NewMessageController(u *usecase.MessageUsecase, auth *auth.Client) *MessageController {
	return &MessageController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

// HandleSendMessage: POST /messages
func (c *MessageController) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	var req model.SendMessageReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondError(w, http.StatusBadRequest, err)
		return
	}

	if req.Content == "" || req.ReceiverID == "" {
		c.respondError(w, http.StatusBadRequest, nil) // エラー詳細を入れるとなお良し
		return
	}

	msg, err := c.Usecase.SendMessage(firebaseUID, req.ReceiverID, req.Content, req.ProductID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, msg)
}

// HandleGetChat: GET /messages?user_id=相手のID
func (c *MessageController) HandleGetChat(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	otherUserID := r.URL.Query().Get("user_id")
	if otherUserID == "" {
		c.respondError(w, http.StatusBadRequest, nil)
		return
	}

	msgs, err := c.Usecase.GetChatHistory(firebaseUID, otherUserID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, msgs)
}

// HandleGetChatList: GET /messages/list
func (c *MessageController) HandleGetChatList(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	chatList, err := c.Usecase.GetChatList(firebaseUID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, chatList)
}

// POST /messages/read?partner_id=相手のID
func (c *MessageController) HandleMarkAsRead(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	partnerID := r.URL.Query().Get("partner_id")
	if partnerID == "" {
		c.respondError(w, http.StatusBadRequest, nil)
		return
	}

	if err := c.Usecase.MarkAsRead(firebaseUID, partnerID); err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (c *MessageController) HandleUnsendMessage(w http.ResponseWriter, r *http.Request) {
	_, err := c.verifyToken(r) // 認証確認
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	messageID := r.PathValue("id")
	if messageID == "" {
		c.respondError(w, http.StatusBadRequest, nil)
		return
	}

	if err := c.Usecase.UnsendMessage(messageID); err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, map[string]string{"status": "unsent"})
}

func (c *MessageController) HandleDeleteMessage(w http.ResponseWriter, r *http.Request) {
	_, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	messageID := r.PathValue("id")
	if messageID == "" {
		c.respondError(w, http.StatusBadRequest, nil)
		return
	}

	if err := c.Usecase.DeleteMessage(messageID); err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
