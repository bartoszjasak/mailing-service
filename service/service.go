package service

import (
	"context"
	"log"
	"sync"

	sqlc "github.com/bartoszjasak/db/sqlc/generated"

	"time"
)

type Service struct {
	db DB
}

type DB interface {
	CreateMessage(ctx context.Context, arg sqlc.CreateMessageParams) error
	DeleteByID(ctx context.Context, id int32) error
	DeleteOlderThan(ctx context.Context, insertTime time.Time) error
	GetByMailingID(ctx context.Context, mailingID int32) ([]sqlc.Message, error)
}

func New(db DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateMessage(ctx context.Context, m Message) error {
	if err := s.db.CreateMessage(ctx, sqlc.CreateMessageParams{
		Email:      m.Email,
		Title:      m.Title,
		Content:    m.Content,
		MailingID:  int32(m.MailingID),
		InsertTime: m.InsertTime,
	}); err != nil {
		return err
	}

	return nil
}

type sendMessage struct {
	MailingID int `json:"mailing_id"`
}

func (s *Service) SendMessages(ctx context.Context, mailingID int) error {
	mails, err := s.db.GetByMailingID(ctx, int32(mailingID))
	if err != nil {
		return err
	}

	for _, mail := range mails {
		log.Println("Sending mail to " + mail.Email +
			", Title: \"" + mail.Title +
			"\", Content: \"" + mail.Content +
			"\"")
		s.db.DeleteByID(ctx, mail.ID)
	}

	return nil
}

func (s *Service) DeleteMessage(ctx context.Context, id int) error {
	return s.db.DeleteByID(ctx, int32(id))
}

func (s *Service) DeleteOlderThan(ctx context.Context, t time.Duration) error {
	return s.db.DeleteOlderThan(ctx, time.Now().Add(-t))
}

// StartCleanupJob runs goroutine that periodicly deletes all messages older than 5 minutes
func (s *Service) StartCleanupJob(wg *sync.WaitGroup, shutdown <-chan struct{}) {
	wg.Add(1)
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.DeleteOlderThan(context.Background(), time.Minute*5)
			case <-shutdown:
				log.Println("Finished background job")
				ticker.Stop()
				wg.Done()
				return
			}
		}
	}()
}
