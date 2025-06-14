package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ApplyMigrations применяет схему и индексы для всех коллекций
func ApplyMigrations(ctx context.Context, db *mongo.Database) error {
	// Очистим коллекции, если они уже есть
	_ = db.Collection("results").Drop(ctx)
	_ = db.Collection("logs").Drop(ctx)

	// === Коллекция results ===
	resultsValidator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"number1", "number2", "result", "operation", "created_at"},
			"properties": bson.M{
				"number1":   bson.M{"bsonType": "double"},
				"number2":   bson.M{"bsonType": "double"},
				"result":    bson.M{"bsonType": "double"},
				"operation": bson.M{"bsonType": "string"},
				"created_at": bson.M{
					"bsonType":    "date",
					"description": "timestamp of result creation",
				},
			},
		},
	}

	if err := db.CreateCollection(ctx, "results", options.CreateCollection().SetValidator(resultsValidator)); err != nil {
		return err
	}

	_, err := db.Collection("results").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "operation", Value: 1},
			{Key: "created_at", Value: -1},
		},
	})
	if err != nil {
		return err
	}

	// === Коллекция logs ===
	logsValidator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"operation", "input", "result", "user_ip", "timestamp"},
			"properties": bson.M{
				"operation": bson.M{"bsonType": "string"},
				"input":     bson.M{"bsonType": "string"},
				"result":    bson.M{"bsonType": "string"},
				"user_ip":   bson.M{"bsonType": "string"},
				"timestamp": bson.M{
					"bsonType":    "date",
					"description": "timestamp of the log entry",
				},
			},
		},
	}

	if err := db.CreateCollection(ctx, "logs", options.CreateCollection().SetValidator(logsValidator)); err != nil {
		return err
	}

	_, err = db.Collection("logs").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "operation", Value: 1},
			{Key: "timestamp", Value: -1},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
