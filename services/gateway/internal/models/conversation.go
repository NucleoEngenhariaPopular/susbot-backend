package models

import "time"

type Message struct {
	UserID    string    `bson:"user_id"`   // ID do usuário da conversa
	Sender    string    `bson:"sender"`    // AccountSID ou BotKit
	Text      string    `bson:"text"`      // Conteúdo da mensagem
	Timestamp time.Time `bson:"timestamp"` // Hora da mensagem
}

type Conversation struct {
	ID        string     `bson:"_id,omitempty"`
	UserID    string     `bson:"user_id"`
	StartTime time.Time  `bson:"start_time"`
	EndTime   *time.Time `bson:"end_time,omitempty"` // ponteiro para permitir ser nil
	Messages  []Message  `bson:"messages"`
}

