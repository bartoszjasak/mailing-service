package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bartoszjasak/service"
	"github.com/go-chi/chi/v5"
)

type MailingService interface {
	CreateMessage(ctx context.Context, newMessage service.Message) error
	SendMessages(ctx context.Context, mailingID int) error
	DeleteMessage(ctx context.Context, id int) error
}

type API struct {
	mailing MailingService
}

func New(srv MailingService) *API {
	return &API{
		mailing: srv,
	}
}

type newMessage struct {
	Email      string    `json:"email"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	MailingID  int       `json:"mailing_id"`
	InsertTime time.Time `json:"insert_time"`
}

func (a *API) postMessage(w http.ResponseWriter, r *http.Request) {
	var m newMessage

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.mailing.CreateMessage(r.Context(), service.Message{
		Email:      m.Email,
		Title:      m.Title,
		Content:    m.Content,
		MailingID:  m.MailingID,
		InsertTime: m.InsertTime,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

type sendMessage struct {
	MailingID int `json:"mailing_id"`
}

func (a *API) sendMessages(w http.ResponseWriter, r *http.Request) {
	var sendMessage sendMessage

	err := json.NewDecoder(r.Body).Decode(&sendMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.mailing.SendMessages(r.Context(), sendMessage.MailingID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) deleteMessage(w http.ResponseWriter, r *http.Request) {
	pathParam := chi.URLParam(r, "id")

	if pathParam == "" {
		http.Error(w, "Message ID should not be empty", http.StatusBadRequest)
	}

	id, err := strconv.Atoi(pathParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := a.mailing.DeleteMessage(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) NewHandler(r *chi.Mux) http.Handler {
	r.Post("/api/messages", a.postMessage)
	r.Post("/api/messages/send", a.sendMessages)
	r.Delete("/api/messages/{id}", a.deleteMessage)

	return r
}
