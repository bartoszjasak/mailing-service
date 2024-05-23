package service

import (
	"context"
	"errors"
	"testing"
	"time"

	sqlc "github.com/bartoszjasak/db/sqlc/generated"
	"github.com/stretchr/testify/require"
)

type MockDB struct {
	createMessageFunc        func(ctx context.Context, arg sqlc.CreateMessageParams) error
	createMessageFuncCounter int
	deleteByIDFunc           func(ctx context.Context, id int32) error
	deleteOlderThanFunc      func(ctx context.Context, insertTime time.Time) error
	getByMailingIDFunc       func(ctx context.Context, mailingID int32) ([]sqlc.Message, error)
}

func (m *MockDB) CreateMessage(ctx context.Context, arg sqlc.CreateMessageParams) error {
	return m.createMessageFunc(ctx, arg)
}

func (m *MockDB) DeleteByID(ctx context.Context, id int32) error {
	return m.deleteByIDFunc(ctx, id)
}

func (m *MockDB) DeleteOlderThan(ctx context.Context, insertTime time.Time) error {
	return m.deleteOlderThanFunc(ctx, insertTime)
}

func (m *MockDB) GetByMailingID(ctx context.Context, mailingID int32) ([]sqlc.Message, error) {
	return m.getByMailingIDFunc(ctx, mailingID)
}

func NewMockDB() *MockDB {
	return &MockDB{
		createMessageFunc:   func(ctx context.Context, arg sqlc.CreateMessageParams) error { return nil },
		deleteByIDFunc:      func(ctx context.Context, id int32) error { return nil },
		deleteOlderThanFunc: func(ctx context.Context, insertTime time.Time) error { return nil },
		getByMailingIDFunc: func(ctx context.Context, mailingID int32) ([]sqlc.Message, error) {
			return []sqlc.Message{}, nil
		},
	}
}

func TestService_CreateMessage(t *testing.T) {
	t.Run("New message should be added to database", func(t *testing.T) {
		db := NewMockDB()
		service := New(db)

		testMessage := Message{
			Email:      "test@test.com",
			Title:      "Test title",
			Content:    "Test content",
			MailingID:  1,
			InsertTime: time.Now(),
		}

		db.createMessageFunc = func(ctx context.Context, arg sqlc.CreateMessageParams) error {
			db.createMessageFuncCounter++

			return nil
		}

		err := service.CreateMessage(context.Background(), testMessage)

		require.NoError(t, err)
		require.Equal(t, 1, db.createMessageFuncCounter)
	})

	t.Run("CreateMessage should return error when inserting message to database failed", func(t *testing.T) {
		db := NewMockDB()
		service := New(db)

		testMessage := Message{
			Email:      "test@test.com",
			Title:      "Test title",
			Content:    "Test content",
			MailingID:  1,
			InsertTime: time.Now(),
		}

		db.createMessageFunc = func(ctx context.Context, arg sqlc.CreateMessageParams) error {
			return errors.New("Write to DB failed.")
		}

		err := service.CreateMessage(context.Background(), testMessage)

		require.Error(t, err)
	})
}
