package value_objects

import (
	"reflect"
	"testing"
)

func TestCountryCodeFromString(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    CountryCode
		wantErr bool
	}{
		{"WhenValidCountryCodeProvided", args{in: "PL"}, CountryCode{value: "PL"}, false},
		{"WhenValidButLowerCasedCountryCodeProvided", args{in: "pl"}, CountryCode{value: "PL"}, false},
		{"WhenInvalidCountryCode", args{in: "PLxyz"}, CountryCode{}, true},
		{"WhenEmptyStringAsInput", args{in: ""}, CountryCode{}, true},
		{"WhenNumbersAsInput", args{in: "12"}, CountryCode{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CountryCodeFromString(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountryCodeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CountryCodeFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
