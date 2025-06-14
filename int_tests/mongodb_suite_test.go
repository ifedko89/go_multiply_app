package int_tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoCollection *mongo.Collection
var mongoContainer testcontainers.Container

func TestMain(m *testing.M) {
	ctx := context.Background()

	// === SETUP ===
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to start Mongo container: %w", err))
	}
	mongoContainer = container

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "27017")

	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	mongoClient = client
	mongoCollection = client.Database("testdb").Collection("users")

	// === RUN TESTS ===
	code := m.Run()

	// === TEARDOWN ===
	_ = client.Disconnect(ctx)
	_ = container.Terminate(ctx)

	os.Exit(code)
}
