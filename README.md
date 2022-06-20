# Users Service

This service exposes simple CRUD for Users management using Go language and MongoDB as persistence layer.

Code is split according to concerns i.e. client request logic is separated from domain services `pkg users` and especially from domain object `User`.
Given trivial example I found no need to introduce more complex architecture than simple layered architecture.

I've used some concepts from DDD like Value Objects in order to encapsulate logic necessary for data validations.

I didn't bother with implementation of `gRPC` as i have not worked with this technology and it will increase amount of time spent on this simple service. 
Code separation allows one to extend service to support additional transport methods like `gRPC`, `websockets`.    

## Events system

According to requirements there I've implemented `Event Bus` concept which is used to notify other parts about domain events happening in this service.
Implementation is limited to simple log of passed event but can be easily extended with support for external bus i.e. `RabbitMQ` or `Kafka`.

## Updates

I am not a fan of single `Users.Update` method or handler which will handle bulk of update operations. This approach leads to 
lots of conditionals and branching in the code which makes it hard to read, test and reason about. I am using small self-descriptive domain operations
i.e. `Users.ChangeEmail`. Please see `users.Service` for implementation of email change.

## Health checks

Service is exposing health checks under `/health` HTTP endpoint. 

## Tests
Tests coverage is limited to `User` domain object which is a core concept used in this service. Most of the code could be tested using integration tests which were not implemented because of time constraints.

Open this project in the terminal on computer with Go `1.18` and run `go test .\***`

Additionally, each build of Docker container (attached Dockerfile) will trigger `go test` which ensures that developer did not introduce unwanted side effects which break application.

## Run
I've attached `Dockerfile` and `docker-compose.yaml` so one can simply run in the terminal:

Environment variables:

`USERS_MONGO_DB=users_service`
`USERS_MONGO_URL=mongodb://root:passw0rd@localhost/?ssl=false&authSource=admin`
`USERS_LISTEN_PORT=3000`

Application by default listens on port `3000`. It can be changed with env variable `USERS_LISTEN_PORT`.

In order to run application expose `ENV` variables listed above and run `go run main.go` or use `docker-compose up`.

## Continuous Integration
There is a configuration set to run tests automatically on each push to `master` branch on Github. Please refer to `https://github.com/adrianbanasiak/go-users-service` -> Actions in order to check the status.

# API schema
I've attached OpenAPI/Swagger specification for this API under `infrastructure/rest/swagger.json`.
One can visit `https://editor.swagger.io` and import this file or simply use code editor with Swagger support built-in.


# Possible improvements
* extended policies to validate data i.e. anti-abuse for `Nickname`
* error mapping to `http.StatusCode`
* increase amount of value objects - `User.Email` could benefit from this one
* extended tests coverage
* contract tests on HTTP handlers level
* passing `RequestID` in the `context` for logging and debugging purposes
* extract password encryption from `Password` value object
* expose application metrics through `/metrics` endpoint
* make some use of `Context` - introduce timers to limit amount of time waited for some action to happen i.e. `database query`.
    If database has not responded for simple `FindByID` in few hundred `ms` we can assume that something bad is going on.
* measure runtime of each action
* use TLS communication with MongoDB
* use soft delete instead of temporal removal of data from the database
* better API specification - describe request / response in detail