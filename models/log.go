package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LogEntry представляет запись в журнале операций
type LogEntry struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Operation string             `bson:"operation" json:"operation"`
	Input     string             `bson:"input" json:"input"`
	Result    string             `bson:"result" json:"result"`
	UserIP    string             `bson:"user_ip" json:"user_ip"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
