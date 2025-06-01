package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Result представляет собой результат умножения двух чисел
type Result struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Number1   float64            `bson:"number1" json:"number1"`
	Number2   float64            `bson:"number2" json:"number2"`
	Product   float64            `bson:"product" json:"product"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
} 