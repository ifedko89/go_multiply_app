package utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupMongo запускает MongoDB-контейнер, подключается к нему и возвращает клиент, коллекцию и teardown-функцию
func SetupMongo(t *testing.T) (*mongo.Client, *mongo.Collection, func()) {
	t.Helper() // Помечает эту функцию как вспомогательную в отчётах о тестах

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(30 * time.Second),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Получаем хост и порт
	host, err := mongoC.Host(ctx)
	require.NoError(t, err)

	port, err := mongoC.MappedPort(ctx, "27017")
	require.NoError(t, err)

	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	// Подключаемся к Mongo
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	require.NoError(t, err)

	// Пингуем БД
	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	collection := client.Database("testdb").Collection("users")

	// teardown-функция, чтобы завершить контейнер и соединение
	tearDown := func() {
		_ = client.Disconnect(ctx)
		_ = mongoC.Terminate(ctx)
	}

	return client, collection, tearDown
}

// BeforeEach очищает коллекцию users перед каждым тестом
func BeforeEach(t *testing.T, collection *mongo.Collection) {
	t.Helper()
	err := collection.Drop(context.Background())
	require.NoError(t, err)
}
