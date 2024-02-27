package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	log.Println("inside broker handler")
	_ = app.writeJSON(w, http.StatusOK, payload)
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) HandleSubmittion(w http.ResponseWriter, r *http.Request) {

	log.Println("inside handlesubmit function")
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	default:
		app.errorJSON(w, errors.New("unkown action fetched"), http.StatusUnauthorized)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	log.Println("inside authentication swtch func")
	jsonData, _ := json.Marshal(a)

	// request, err := http.NewRequest("POST", "http://localhost:8081/authenticate", bytes.NewBuffer(jsonData))
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	log.Println(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("unkown action fetched"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("unkown action fetched"), http.StatusUnauthorized)
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if jsonFromService.Error {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Authentication successfull"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)

}
