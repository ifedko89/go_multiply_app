package int_tests

import (
	"context"
	"github.com/igor-fedko/go_multiply_app/int_tests/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestMongoDBContainer(t *testing.T) {
	utils.BeforeEach(t, mongoCollection) // очистка перед тестом
	// Вставка документа
	user := bson.M{"name": "Игорь"}
	res, err := mongoCollection.InsertOne(context.Background(), user)
	assert.NoError(t, err)

	// Чтение документа
	var result bson.M
	err = mongoCollection.FindOne(context.Background(), bson.M{"_id": res.InsertedID}).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "Игорь", result["name"])
	defer mongoCollection.Drop(context.Background()) // удаляет коллекцию

}

func TestMongoDB_FindNonExistent(t *testing.T) {
	utils.BeforeEach(t, mongoCollection) // очистка перед тестом
	// Чтение несуществующего документа
	var result bson.M
	err := mongoCollection.FindOne(context.Background(), bson.M{"_id": "fake_id"}).Decode(&result)
	assert.Error(t, err)                             // ожидаем ошибку, так как документа нет
	defer mongoCollection.Drop(context.Background()) // удаляет коллекцию

}
