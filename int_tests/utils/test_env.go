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

type TestEnv struct {
	T          *testing.T
	Container  testcontainers.Container
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewMongoEnv(t *testing.T) *TestEnv {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "27017")
	require.NoError(t, err)

	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)

	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	collection := client.Database("testdb").Collection("users")

	err = ApplyMigrations(ctx, client.Database("testdb"))
	require.NoError(t, err)

	return &TestEnv{
		T:          t,
		Container:  container,
		Client:     client,
		Collection: collection,
	}
}

func (e *TestEnv) BeforeEach() {
	err := e.Collection.Drop(context.Background())
	require.NoError(e.T, err)
}

func (e *TestEnv) Close() {
	_ = e.Client.Disconnect(context.Background())
	_ = e.Container.Terminate(context.Background())
}
