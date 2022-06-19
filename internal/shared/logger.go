package shared

import (
	"go.uber.org/zap"
	"log"
)

type Logger interface {
	Info(...any)
	Infow(string, ...any)
	Errorw(string, ...any)
}

type TestsLogger struct{}

func NewTestsLogger() *TestsLogger {
	return &TestsLogger{}
}

func (t TestsLogger) Info(...any) {
}

func (t TestsLogger) Infow(string, ...any) {
}

func (t TestsLogger) Errorw(in string, args ...any) {
	log.Println(in, args)
}

func NewLogger() (*zap.SugaredLogger, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return l.Sugar(), nil
}
