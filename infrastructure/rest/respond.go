package rest

import (
	"encoding/json"
	"net/http"
	"time"
)

func RespondSuccess(w http.ResponseWriter, result any, statusCode int) error {
	resp := Response{
		Data:       result,
		Successful: true,
		SentAt:     time.Now(),
	}

	res, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(statusCode)
	}

	_, err = w.Write(res)
	return err
}

func RespondError(w http.ResponseWriter, errors []error, statusCode int) error {
	respErrors := make([]string, 0)

	for _, err := range errors {
		respErrors = append(respErrors, err.Error())
	}

	resp := Response{
		Errors: respErrors,
		SentAt: time.Now().UTC(),
	}

	res, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(statusCode)
	}

	_, err = w.Write(res)
	return err
}
