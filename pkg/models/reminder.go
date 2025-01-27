package models

import "time"

type Reminder struct {
	AuthorID   string
	ReminderID string
	Content    string
	TargetName string
	ExpiredAt  *time.Time
}
