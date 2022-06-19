package healthchecks

import (
	"context"
	"github.com/etherlabsio/healthcheck"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func NewService(mongoClient *mongo.Client) *Service {
	return &Service{mongoClient: mongoClient}
}

type Service struct {
	mongoClient *mongo.Client
}

func (s *Service) mongoHealthcheck(ctx context.Context) error {
	return s.mongoClient.Ping(ctx, nil)
}

func (s *Service) ConfigureHandler() http.HandlerFunc {
	return healthcheck.HandlerFunc(
		healthcheck.WithTimeout(5*time.Second),
		healthcheck.WithChecker("mongoDB", healthcheck.CheckerFunc(s.mongoHealthcheck)))
}
