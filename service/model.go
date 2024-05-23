package service

import "time"

type Message struct {
	Email      string
	Title      string
	Content    string
	MailingID  int
	InsertTime time.Time
}
