package main


import (

	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
	"encoding/json"
	"bytes"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)


func TestApp(t *testing.T) {
	

	P := &AppStatusStruct{
		UPDATE_BY: "test3",
	}

	a := createTestDBConnection()

	router := mux.NewRouter()
	router.HandleFunc("/hello", SayHello).Methods("GET")
	router.HandleFunc("/api/app/{id}", a.DisplayAppByID).Methods("GET")
	router.HandleFunc("/api/app/new", a.AddNewApp).Methods("POST")
	router.HandleFunc("/api/app/{id}", a.UpdateData).Methods("PUT")
	router.HandleFunc("/api/app/{id}", a.DeleteData).Methods("DELETE")
	router.HandleFunc("/api/apps", a.DisplayAllApp).Methods("GET")
	router.HandleFunc("/", a.DisplayHtml).Methods("GET")
	ts := httptest.NewServer(router)
    defer ts.Close()

	url1 := ts.URL + "/hello"
	url2 := ts.URL + "/api/app/" + "10"
	url3 := ts.URL + "/api/apps"
	url4 := ts.URL + "/"

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


	payload, err := json.Marshal(P)
	if err != nil {
		fmt.Printf("Error while marshall in TestAddNewApp: %v", err)
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:5000/api/app/new", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Error while make new POST Request: %v", err)
	}	

	resp5 := httptest.NewRecorder()

	a.AddNewApp(resp5, req)

	assert.Equal(t, 200,  resp1.StatusCode)
	assert.Equal(t, 200,  resp2.StatusCode)
	assert.Equal(t, 200,  resp3.StatusCode)
	assert.Equal(t, 200,  resp4.StatusCode)
	assert.Equal(t, 200,  resp5.Code)

}
