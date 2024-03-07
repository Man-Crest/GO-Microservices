package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"logger/logs"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	Log    LogPayload  `json:"log,omitempty"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
		log.Println("1")
		app.errorJSON(w, err)
		return
	}

	log.Println(requestPayload.Action)
	log.Println(requestPayload.Log.Data)

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.LogViaRpc(w, requestPayload.Log)
	case "show":
		app.ShowData(w)
	default:
		app.errorJSON(w, errors.New("unkown action fetched"), http.StatusUnauthorized)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	log.Println("inside authentication switch func")
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

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}
func (app *Config) ShowData(w http.ResponseWriter) {

	log.Println("inside showData function")

	jsonData := []byte{} // Empty JSON data, assuming no data to send in the request body
	logServiceURL := "http://logger-service/all"

	request, err := http.NewRequest("GET", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("error creating HTTP request: %v", err)
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("error sending HTTP request: %v", err)
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		errMsg := fmt.Errorf("unexpected status code: %d", response.StatusCode)
		// log.Printf("error response from server: %s", errMsg)
		app.errorJSON(w, errMsg)
		return
	}

	// Process response body if needed
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) LogViaRpc(w http.ResponseWriter, l LogPayload) {
	log.Println("inside log via RPC")
	client, err := rpc.Dial("tcp", "logger-service:5001")

	if err != nil {
		err := errors.Join(err, errors.New("after dial service caled"))
		app.errorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string

	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		err := errors.Join(err, errors.New("after call service called"))
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	log.Println(requestPayload)

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}
