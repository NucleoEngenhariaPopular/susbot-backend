package models

import "time"

type Message struct {
	UserID    string    `json:"user_id" bson:"user_id"`
	Sender    string    `json:"sender" bson:"sender"`
	Text      string    `json:"text" bson:"text"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type Conversation struct {
	ID        string     `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string     `json:"user_id" bson:"user_id"`
	StartTime time.Time  `json:"start_time" bson:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty" bson:"end_time,omitempty"`
	Messages  []Message  `json:"messages" bson:"messages"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

