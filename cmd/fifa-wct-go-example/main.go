package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NeowayLabs/fifa-wct-go-example/internal/application"
	"github.com/NeowayLabs/fifa-wct-go-example/internal/infrastructure/repository/mongo"
	"github.com/NeowayLabs/fifa-wct-go-example/internal/infrastructure/serve/rest"
)

const (
	envVarMongoURL          = "MONGO_URL"
	envVarMongoDatabaseName = "MONGO_DATABASE_NAME"
	envVarHTTPServerHost    = "HTTP_SERVER_HOST"
	envVarHTTPServerPort    = "HTTP_SERVER_PORT"

	defaultShutdownTimeout   = 60 * time.Second
	defaultMongoDatabaseName = "fifa-wct-go-example"
	defaultHTTPServerHost    = "0.0.0.0"
	defaultHTTPServerPort    = "80"
)

var (
	ErrEnvIsRequired     = errors.New("env is required")
	version, build, date string
)

func main() {
	log := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	log.Printf("fifa-wct-go-example starting - build:%s; date:%s; version:%s", build, date, version)

	// Check required variables
	if err := checkRequiredEnvs(envVarMongoURL); err != nil {
		log.Fatalf("error missing variables: %v", err)
		return
	}

	// Repositories
	ctx := context.Background()

	mongoClient, err := mongo.Connect(ctx, getMongoURL())
	if err != nil {
		log.Fatalf("error connecting database: %v", err)
		return
	}

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("error disconnecting database: %v", err)
		}
	}()

	mongoDataBase := mongo.NewDatabase(mongoClient, getMongoDatabaseName())

	teamRepository, err := mongo.NewTeamRepository(mongoDataBase, log)
	if err != nil {
		log.Fatalf("error on creating team repository instance: %v", err)
	}

	// Domain services

	// Application services
	teamService := application.NewTeamService(teamRepository)

	// Server HTTP
	handler := rest.NewHandler(teamService, log)
	server := rest.NewServer(handler, getHTTPServerHost(), getHTTPServerPort(), log)
	server.ListenAndServe()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan

	if err := server.Shutdown(defaultShutdownTimeout); err != nil {
		log.Fatalf("error on shutdown HTTP server: %v", err)
		return
	}
}

func getMongoURL() string {
	return os.Getenv(envVarMongoURL)
}

func getMongoDatabaseName() string {
	return getEnvOrDefaultValue(envVarMongoDatabaseName, defaultMongoDatabaseName)
}

func getHTTPServerHost() string {
	return getEnvOrDefaultValue(envVarHTTPServerHost, defaultHTTPServerHost)
}

func getHTTPServerPort() string {
	return getEnvOrDefaultValue(envVarHTTPServerPort, defaultHTTPServerPort)
}

func getEnvOrDefaultValue(envVar string, defaultValue string) string {
	v := os.Getenv(envVar)
	if v == "" && len(defaultValue) > 0 {
		return defaultValue
	}

	return v
}

func checkRequiredEnvs(envVarArgs ...string) error {
	for _, envVar := range envVarArgs {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("environment variable '%s': %w", envVar, ErrEnvIsRequired)
		}
	}

	return nil
}
