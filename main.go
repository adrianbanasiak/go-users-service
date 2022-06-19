package main

import (
	"fmt"
	"github.com/adrianbanasiak/go-users-service/infrastructure/rest"
	"github.com/adrianbanasiak/go-users-service/internal/events"
	"github.com/adrianbanasiak/go-users-service/internal/shared"
	"github.com/adrianbanasiak/go-users-service/internal/users"
	"os"
	"strconv"
)

func main() {
	log, err := shared.NewLogger()
	if err != nil {
		panic(fmt.Sprintf("failed to configure logger: %s", err))
	}

	port := 3000
	if listenPort, ok := os.LookupEnv("USERS_LISTEN_PORT"); ok {
		port, err = strconv.Atoi(listenPort)
		if err != nil {
			panic("Invalid value for USERS_LISTEN_PORT env variable. Want integer value")
		}
	}

	bus := events.NewInternalEventBus(log)
	usersRepository := users.NewInMemoryRepository(nil)
	usersService := users.NewService(usersRepository, log, bus)

	restServer := rest.NewServer(rest.Dependencies{
		Log:          log,
		UsersService: usersService,
	}, port)

	err = restServer.Start()
	if err != nil {
		log.Fatal("failed to start REST API server")
	}
}
