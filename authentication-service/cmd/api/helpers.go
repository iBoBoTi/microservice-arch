package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

func (app *Config) readJSON (w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes :=  1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(data); err != nil {
		return fmt.Errorf("error decoding json body %s", err)
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}
	
	return nil
}

func (app *Config) writeJSON (w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data to json: %s", err)
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return fmt.Errorf("error writing response json %s", err)
	}

	return nil
}

func (app *Config) errorJSON (w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message =  err.Error()

	return app.writeJSON(w, statusCode, payload)
}