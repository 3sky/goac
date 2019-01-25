package main


import (
	"p2go/db"
	"p2go/api"
	st "p2go/static"
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)


func TestMain(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/hello", api.SayHello).Methods("GET")
	router.HandleFunc("/api/app/{id}", api.DisplaAppByID).Methods("GET")
	router.HandleFunc("/api/app/new", api.AddNewApp).Methods("POST")
	router.HandleFunc("/api/app/{id}", api.UpdateData).Methods("PUT")
	router.HandleFunc("/api/app/{id}", api.DeleteData).Methods("DELETE")
	router.HandleFunc("/api/apps", api.DisplayAllApp).Methods("GET")
	router.HandleFunc("/", st.DisplayHtml).Methods("GET")
	ts := httptest.NewServer(router)
    defer ts.Close()

	url1 := ts.URL + "/hello"
	url2 := ts.URL + "/api/app/" + "2"
	url3 := ts.URL + "/api/app/" + "new"
	url4 := ts.URL + "/api/apps"
	url5 := ts.URL + "/"

	resp1, err := http.Get(url1)
	db.CheckErr(err)
	resp2, err := http.Get(url2)
	db.CheckErr(err)
	resp3, err := http.Get(url3)
	db.CheckErr(err)
	resp4, err := http.Get(url4)
	db.CheckErr(err)
	resp5, err := http.Get(url5)
	db.CheckErr(err)

	assert.Equal(t, 200,  resp1.StatusCode)
	assert.Equal(t, 200,  resp2.StatusCode)
	assert.Equal(t, 200,  resp3.StatusCode)
	assert.Equal(t, 200,  resp4.StatusCode)
	assert.Equal(t, 200,  resp5.StatusCode)

}
