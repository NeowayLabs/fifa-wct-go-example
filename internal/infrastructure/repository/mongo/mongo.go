package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrMongoURLIsRequired = errors.New("mongoURL is required")
)

func Connect(ctx context.Context, mongoURL string) (*mongo.Client, error) {
	if strings.TrimSpace(mongoURL) == "" {
		return nil, ErrMongoURLIsRequired
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, fmt.Errorf("error on MongoDB connection: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error on MongoDB connection: %w", err)
	}

	return client, nil
}

func NewDatabase(client *mongo.Client, databaseName string) *mongo.Database {
	return client.Database(databaseName)
}

func IsDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key")
}

func IsNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "no documents in result")
}
