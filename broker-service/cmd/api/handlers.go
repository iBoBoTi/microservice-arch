package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	authAction = "auth"
)

type RequestPayload struct{
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request){
	payload:= jsonResponse{
		Error: false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request){

	var reqPayload RequestPayload
	
	if err := app.readJSON(w, r, &reqPayload); err != nil{
		app.errorJSON(w, err)
		return
	}

	switch reqPayload.Action {
	case authAction:
		app.authenticate(w, reqPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload){
	// create some json to send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	
	req, err := http.NewRequest("POST","http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get the correct status code
	if response.StatusCode == http.StatusUnauthorized{
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted{
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// read response.Body
	var jsonFromService jsonResponse

	if err := json.NewDecoder(response.Body).Decode(&jsonFromService); err != nil{
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error{
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message =  "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}