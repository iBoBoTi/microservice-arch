package main

import (
	"net/http"

	"github.com/iBoBoTi/microservice-arch/logger-service/data"
)

type JsonPayload struct {
	Name string
	Data string
}

func (app *Config) WriteLog(rw http.ResponseWriter, req *http.Request) {
	//read json
	var reqPayload JsonPayload

	_ = app.readJSON(rw, req, &reqPayload)

	//insert data
	event := data.LogEntry{
		Name: reqPayload.Name,
		Data: reqPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(rw, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(rw, http.StatusAccepted, resp)

}
