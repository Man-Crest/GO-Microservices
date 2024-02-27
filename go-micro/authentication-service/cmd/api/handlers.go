package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	log.Println("inside auth handler func")

	requestPayload := &AuthPayload{}

	err := json.NewDecoder(r.Body).Decode(&requestPayload)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(requestPayload)

	// validate the user against the database

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Println("error in getting email in auth")
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		log.Println("password is not matching in auth")
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
