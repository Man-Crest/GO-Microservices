package main

import (
	"log"
	"logger/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) AllData(w http.ResponseWriter, r *http.Request) {

	data, err := app.Models.LogEntry.All()
	if err != nil {
		log.Println(err, "error reading log")
	}

	// err = app.errorJSON(w, errors.New("sdgnsdjhfbgndfjk"), http.StatusBadRequest)
	// resp := "this is response string"
	resp := jsonResponse{
		Error:   false,
		Message: "showed",
		Data:    data,
	}
	err = app.writeJSON(w, http.StatusAccepted, resp)
	if err != nil {
		app.errorJSON(w, err)
		log.Println(err)
		return
	}
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		log.Println(err, "error reading log")
	}

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		log.Println(err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	err = app.writeJSON(w, http.StatusAccepted, resp)
	if err != nil {
		app.errorJSON(w, err)
		log.Println(err)
		return
	}
}
