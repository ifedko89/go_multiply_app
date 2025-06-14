package int_tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestMongoDB_InsertValidResult(t *testing.T) {
	env.BeforeEach()

	doc := bson.M{
		"number1":    10.5,
		"number2":    5.0,
		"result":     15.5,
		"operation":  "add",
		"created_at": time.Now(),
	}

	_, err := env.Client.Database("testdb").Collection("results").InsertOne(context.Background(), doc)
	assert.NoError(t, err)
}

func TestMongoDB_InsertInvalidResult_MissingField(t *testing.T) {
	env.BeforeEach()

	doc := bson.M{
		"number1": 10.5,
		"number2": 5.0,
		// result пропущен
		"operation":  "add",
		"created_at": time.Now(),
	}

	_, err := env.Client.Database("testdb").Collection("results").InsertOne(context.Background(), doc)
	assert.Error(t, err)
}

func TestMongoDB_InsertInvalidResult_WrongType(t *testing.T) {
	env.BeforeEach()

	doc := bson.M{
		"number1":    "wrong", // должно быть float64
		"number2":    5.0,
		"result":     15.5,
		"operation":  "add",
		"created_at": time.Now(),
	}

	_, err := env.Client.Database("testdb").Collection("results").InsertOne(context.Background(), doc)
	assert.Error(t, err)
}

func TestMongoDB_InsertDuplicateLogIndex(t *testing.T) {
	env.BeforeEach()

	log := bson.M{
		"operation": "divide",
		"input":     "10/2",
		"result":    "5",
		"user_ip":   "192.168.0.1",
		"timestamp": time.Now(),
	}

	coll := env.Client.Database("testdb").Collection("logs")

	_, err := coll.InsertOne(context.Background(), log)
	assert.NoError(t, err)

	// Повтор с тем же operation + timestamp — конфликт, если бы индекс был уникальным
	_, err = coll.InsertOne(context.Background(), log)
	// ❗ assert.NoError(t, err) — если индекс не уникальный
	// ✅ assert.Error(t, err) — если индекс уникальный (в данном случае, у нас НЕ уникальный)
	assert.NoError(t, err)
}

func TestMongoDB_LogInsert_ValidationFails(t *testing.T) {
	env.BeforeEach()

	doc := bson.M{
		"operation": "subtract",
		// input и result отсутствуют
		"user_ip":   "192.168.0.1",
		"timestamp": time.Now(),
	}

	_, err := env.Client.Database("testdb").Collection("logs").InsertOne(context.Background(), doc)
	assert.Error(t, err)
}
