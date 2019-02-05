package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	router := mux.NewRouter()
	router.HandleFunc("/hello", SayHello).Methods("GET")
	router.HandleFunc("/api/app/{id}", a.DisplaAppByID).Methods("GET")
	router.HandleFunc("/api/app/new", a.AddNewApp).Methods("POST")
	router.HandleFunc("/api/app/{id}", a.UpdateData).Methods("PUT")
	router.HandleFunc("/api/app/{id}", a.DeleteData).Methods("DELETE")
	router.HandleFunc("/api/apps", a.DisplayAllApp).Methods("GET")
	router.HandleFunc("/", a.DisplayHTML).Methods("GET")
	ts := httptest.NewServer(router)
	defer ts.Close()

	url1 := ts.URL + "/hello"
	url2 := ts.URL + "/api/app/" + "10"
	url3 := ts.URL + "/api/app/" + "new"
	url4 := ts.URL + "/api/apps"
	url5 := ts.URL + "/"

	resp1, err := http.Get(url1)
	CheckErr(err)
	resp2, err := http.Get(url2)
	CheckErr(err)
	resp3, err := http.Get(url3)
	CheckErr(err)
	resp4, err := http.Get(url4)
	CheckErr(err)
	resp5, err := http.Get(url5)
	CheckErr(err)

	assert.Equal(t, 200, resp1.StatusCode)
	assert.Equal(t, 200, resp2.StatusCode)
	assert.Equal(t, 200, resp3.StatusCode)
	assert.Equal(t, 200, resp4.StatusCode)
	assert.Equal(t, 200, resp5.StatusCode)

}
