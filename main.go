package main

import (
	"fmt"
	"github.com/adrianbanasiak/go-users-service/infrastructure/healthchecks"
	"github.com/adrianbanasiak/go-users-service/infrastructure/mongodb"
	"github.com/adrianbanasiak/go-users-service/infrastructure/rest"
	"github.com/adrianbanasiak/go-users-service/internal/events"
	"github.com/adrianbanasiak/go-users-service/internal/shared"
	"github.com/adrianbanasiak/go-users-service/internal/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"strconv"
)

func main() {
	log, err := shared.NewLogger()
	if err != nil {
		panic(fmt.Sprintf("failed to configure logger: %s", err))
	}

	c := connectMongo(log)
	mongoDBName, ok := os.LookupEnv("USERS_MONGO_DB")
	if !ok {
		log.Fatal("USERS_MONGO_DB environment variable is missing")
	}

	db := c.Database(mongoDBName)

	bus := events.NewInternalEventBus(log)
	usersRepository := users.NewMongoRepository(log, db)
	usersService := users.NewService(usersRepository, log, bus)

	healthchecksService := healthchecks.NewService(c)

	port := 3000
	if listenPort, ok := os.LookupEnv("USERS_LISTEN_PORT"); ok {
		port, err = strconv.Atoi(listenPort)
		if err != nil {
			panic("Invalid value for USERS_LISTEN_PORT env variable. Want integer value")
		}
	}

	restServer := rest.NewServer(rest.Dependencies{
		Log:                 log,
		UsersService:        usersService,
		HealthchecksService: healthchecksService,
	}, port)

	err = restServer.Start()
	if err != nil {
		log.Fatal("failed to start REST API server")
	}
}

func connectMongo(log *zap.SugaredLogger) *mongo.Client {
	mongoURL, ok := os.LookupEnv("USERS_MONGO_URL")
	if !ok {
		log.Fatal("USERS_MONGO_URL environment variable is missing")
	}

	c, err := mongodb.Init(mongoURL)
	if err != nil {
		log.Fatalw("failed to establish connection with MongoDB", "error", err)
	}

	log.Info("established connection with MongoDB")

	return c
}
