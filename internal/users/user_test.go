package users

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	type args struct {
		req CreateUserReq
	}
	tests := []struct {
		name    string
		args    args
		want    User
		wantErr bool
	}{
		{"WhenInputIsCorrect",
			args{req: CreateUserReq{FirstName: "PwnLord", Password: "dAtriFPLcVJbNrZfG5hH", CountryCode: "PL"}},
			User{FirstName: "PwnLord"}, false},
		{"WhenPasswordIsEmpty",
			args{req: CreateUserReq{FirstName: "PwnLord", Password: "", CountryCode: "PL"}},
			User{FirstName: "PwnLord"}, true},
		{"WhenCountryCodeIsInvalid",
			args{req: CreateUserReq{FirstName: "PwnLord", Password: "dAtriFPLcVJbNrZfG5hH", CountryCode: "XYZ"}},
			User{FirstName: "PwnLord"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
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
