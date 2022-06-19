package users

import (
	"context"
	"errors"
	"github.com/adrianbanasiak/go-users-service/internal/events"
	"github.com/adrianbanasiak/go-users-service/internal/shared"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestService_CreateUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req CreateUserReq
	}
	tests := []struct {
		name           string
		repository     *InMemoryRepository
		eventbus       *events.TestsEventBus
		args           args
		want           User
		wantErr        bool
		wantRepoCalled bool
		wantBusCalled  bool
	}{
		{"WhenValidRequestIsProvided",
			NewInMemoryRepository(nil),
			events.NewTestsEventBus(nil),
			args{
				ctx: context.Background(),
				req: CreateUserReq{Password: "QIPbiN5zTCGMrz6uyo31", CountryCode: "PL"},
			},
			User{}, false, true, true},
		{"WhenEmitFails",
			NewInMemoryRepository(nil),
			events.NewTestsEventBus(errors.New("failed")),
			args{
				ctx: context.Background(),
				req: CreateUserReq{Password: "QIPbiN5zTCGMrz6uyo31", CountryCode: "PL"},
			},
			User{}, true, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &Service{
				repository: tt.repository,
				log:        shared.NewTestsLogger(),
				bus:        tt.eventbus,
			}

			got, err := s.CreateUser(tt.args.ctx, tt.args.req)

			if tt.wantRepoCalled != tt.repository.called {
				t.Errorf("CreateUser() got = %v, wantRepositoryCalled %v", tt.repository.Called(), tt.wantRepoCalled)
				return
			}

			if tt.wantBusCalled != tt.eventbus.Called() {
				t.Errorf("CreateUser() got = %v, wantBusCalled %v", tt.eventbus.Called(), tt.wantBusCalled)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				gotEq := User{NickName: got.NickName, FirstName: got.FirstName, LastName: got.LastName, Email: got.Email}
				if !reflect.DeepEqual(gotEq, tt.want) {
					t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
				}

				if got.ID == uuid.Nil {
					t.Errorf("CreateUser() got = %v, want non nil UUID", got)
				}

				if !got.CreatedAt.Before(time.Now()) || got.CreatedAt.IsZero() {
					t.Errorf("CreateUser() got = %v, want valid createdAt", got)
				}
			}
		})
	}
}
