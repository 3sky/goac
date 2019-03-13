package main


import (

	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/gorilla/mux"
	
)


func TestAuth(t *testing.T) {
	
	var client http.Client

	
	p := &User{
		username: "3sky",
		password: "test",
	}
	
	router := mux.NewRouter()
	router.HandleFunc("/hello", use(SayHello, p.basicAuth)).Methods("GET")
	ts := httptest.NewServer(router)
    defer ts.Close()

	req1, err := http.NewRequest("GET", ts.URL + "/hello", nil)
	errorWithCreateAuthRequest(err)
	req1.SetBasicAuth(p.username,p.password)
	
	res1, err := client.Do(req1)
	errorWithCreateAuthRequest(err)

	req2, err := http.NewRequest("GET", ts.URL + "/hello", nil)
	req2.SetBasicAuth("test","test")
	errorWithCreateAuthRequest(err)
	res2, err := client.Do(req2)
	errorWithCreateAuthRequest(err)

	req3, err := http.NewRequest("GET", ts.URL + "/hello", nil)
	errorWithCreateAuthRequest(err)
	req3.SetBasicAuth("3sky","test1")
	res3, err := client.Do(req3)
	errorWithCreateAuthRequest(err)

	req4, err := http.NewRequest("GET", ts.URL + "/hello", nil)
	errorWithCreateAuthRequest(err)
	req4.SetBasicAuth("","")
	res4, err := client.Do(req4)
	errorWithCreateAuthRequest(err)

	assert.Equal(t, 200,  res1.StatusCode)
	assert.Equal(t, 401,  res2.StatusCode)
	assert.Equal(t, 401,  res3.StatusCode)
	assert.Equal(t, 401,  res4.StatusCode)


}


func errorWithCreateAuthRequest(err error) {
	if err != nil {
		fmt.Printf("Error with create new request in TestAuth: %v", err)
	}
} 