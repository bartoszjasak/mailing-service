package service

import (
	"database/sql"
	"encoding/json"
	"log"
	sqlc "main/db/sqlc/generated"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Service struct {
	DB *sqlc.Queries
}

func New(db *sql.DB) *Service {
	return &Service{
		DB: sqlc.New(db),
	}
}

type newMessage struct {
	Email      string    `json:"email"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	MailingID  int       `json:"mailing_id"`
	InsertTime time.Time `json:"inser_time"`
}

func (s *Service) PostMessage(w http.ResponseWriter, r *http.Request) {
	var newMessage newMessage

	err := json.NewDecoder(r.Body).Decode(&newMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.DB.CreateMessage(r.Context(), sqlc.CreateMessageParams{
		Email:      newMessage.Email,
		Title:      newMessage.Title,
		Content:    newMessage.Content,
		MailingID:  int32(newMessage.MailingID),
		InsertTime: newMessage.InsertTime,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type sendMessage struct {
	MailingID int `json:"mailing_id"`
}

func (s *Service) SendMessages(w http.ResponseWriter, r *http.Request) {
	var sendMessage sendMessage

	err := json.NewDecoder(r.Body).Decode(&sendMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mails, err := s.DB.GetByMailingID(r.Context(), int32(sendMessage.MailingID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, mail := range mails {
		log.Println("Sending mail to: " + mail.Email)
		s.DB.DeleteByID(r.Context(), mail.ID)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	pathParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(pathParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = s.DB.DeleteByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
