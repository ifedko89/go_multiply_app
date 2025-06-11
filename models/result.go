package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Result представляет собой результат математической операции над двумя числами
type Result struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Number1   float64            `bson:"number1" json:"number1"`
	Number2   float64            `bson:"number2" json:"number2"`
	Result    float64            `bson:"result" json:"result"`
	Operation string             `bson:"operation" json:"operation"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
} 