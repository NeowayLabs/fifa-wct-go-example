//go:build integration

package mongo_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	internal_mongo "github.com/NeowayLabs/fifa-wct-go-example/internal/infrastructure/repository/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mongoClient *mongo.Client
	testLog     *log.Logger
)

func init() {
	testLog = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	mongoURL := os.Getenv("MONGO_URL")

	client, err := internal_mongo.Connect(context.Background(), mongoURL)
	if err != nil {
		log.Fatalf("error on MongoDB connection: %v", err)
	}
	mongoClient = client
}

func newDatabaseForTest() *mongo.Database {
	return mongoClient.Database(primitive.NewObjectID().Hex())
}

func TestConnect_Success(t *testing.T) {
	mongoURL := os.Getenv("MONGO_URL")
	mongoClient, err := internal_mongo.Connect(context.Background(), mongoURL)
	assert.NoError(t, err)
	assert.NotNil(t, mongoClient)

	if mongoClient != nil {
		database := internal_mongo.NewDatabase(mongoClient, "fifa-wct-go-example-test")
		assert.NotNil(t, database)
	}
}

func TestConnect_EmptyURL(t *testing.T) {
	mongoClient, err := internal_mongo.Connect(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, mongoClient)
}

func TestConnect_InvalidURL(t *testing.T) {
	mongoClient, err := internal_mongo.Connect(context.Background(), "192.168.0.1")
	assert.Error(t, err)
	assert.Nil(t, mongoClient)
}

func TestIsDuplicateKeyError(t *testing.T) {
	assert.True(t, internal_mongo.IsDuplicateKeyError(fmt.Errorf("duplicate key")))
	assert.False(t, internal_mongo.IsDuplicateKeyError(fmt.Errorf("error")))
}

func TestIsNotFoundError(t *testing.T) {
	assert.True(t, internal_mongo.IsNotFoundError(fmt.Errorf("no documents in result")))
	assert.False(t, internal_mongo.IsNotFoundError(fmt.Errorf("error")))
}
