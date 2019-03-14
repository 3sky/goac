package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {

	P := &AppStatusStruct{
		UpdateBy: "test3",
	}

	a := createTestDBConnection()

	router := mux.NewRouter()
	router.HandleFunc("/hello", SayHello).Methods("GET")
	router.HandleFunc("/api/app/{id}", a.DisplayAppByID).Methods("GET")
	router.HandleFunc("/api/app/new", a.AddNewApp).Methods("POST")
	router.HandleFunc("/api/app/{id}", a.UpdateData).Methods("PUT")
	router.HandleFunc("/api/app/{id}", a.DeleteData).Methods("DELETE")
	router.HandleFunc("/api/apps", a.DisplayAllApp).Methods("GET")
	router.HandleFunc("/dev", a.DisplayHTMLDev).Methods("GET")
	router.HandleFunc("/stg", a.DisplayHTMLStg).Methods("GET")

	ts := httptest.NewServer(router)
	defer ts.Close()

	url1 := ts.URL + "/hello"
	url2 := ts.URL + "/api/app/" + "10"
	url3 := ts.URL + "/api/apps"
	url4 := ts.URL + "/dev"
	url5 := ts.URL + "/stg"

	resp1, err := http.Get(url1)
	if err != nil {
		fmt.Printf("Error while make new GET Request: %v", err)
	}

	resp2, err := http.Get(url2)
	if err != nil {
		fmt.Printf("Error while make new GET Request: %v", err)
	}

	resp3, err := http.Get(url3)
	if err != nil {
		fmt.Printf("Error while make new GET Request: %v", err)
	}

	resp4, err := http.Get(url4)
	if err != nil {
		fmt.Printf("Error while make new GET Request: %v", err)
	}

	resp5, err := http.Get(url5)
	if err != nil {
		fmt.Printf("Error while make new GET Request: %v", err)
	}

	payload, err := json.Marshal(P)
	if err != nil {
		fmt.Printf("Error while marshall in TestAddNewApp: %v", err)
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:5000/api/app/new", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Error while make new POST Request: %v", err)
	}

	resp6 := httptest.NewRecorder()

	a.AddNewApp(resp6, req)

	assert.Equal(t, 200, resp1.StatusCode)
	assert.Equal(t, 200, resp2.StatusCode)
	assert.Equal(t, 200, resp3.StatusCode)
	assert.Equal(t, 200, resp4.StatusCode)
	assert.Equal(t, 200, resp5.StatusCode)
	assert.Equal(t, 200, resp6.Code)

}
