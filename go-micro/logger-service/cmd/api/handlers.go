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

	log.Println("inside alldata logs")
	data, err := app.Models.LogEntry.All()
	if err != nil {
		log.Println(err, "error reading log")
	}

	for i, val := range data {
		log.Println(i, val.Data)
		err = app.writeJSON(w, http.StatusAccepted, val.Data)
		if err != nil {
			log.Println(err, "error writing log")
		}
	}

}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	err := app.readJSON(w, r, &requestPayload)

	log.Println(requestPayload.Name)
	log.Println(requestPayload.Data)
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

	app.writeJSON(w, http.StatusAccepted, resp)
}
