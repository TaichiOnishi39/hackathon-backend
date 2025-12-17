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

	msg, err := c.Usecase.SendMessage(firebaseUID, req.ReceiverID, req.Content)
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
