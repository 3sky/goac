package main


import (

	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/gorilla/mux"
	
)


func TestAuth(t *testing.T) {
	
	var client http.Client
	var err, err2 error
	
	p := &User{
		username: "3sky",
		password: "test",
	}
	
	router := mux.NewRouter()
	router.HandleFunc("/hello", use(SayHello, p.basicAuth)).Methods("GET")
	ts := httptest.NewServer(router)
    defer ts.Close()

	req1, err := http.NewRequest("GET", ts.URL + "/hello", nil)
	CheckErr(err)
	req1.SetBasicAuth(p.username,p.password)
	res1, err2 := client.Do(req1)
	CheckErr(err2)

	req2, err := http.NewRequest("GET", ts.URL + "/hello", nil)
	CheckErr(err)
	req1.SetBasicAuth("test","test")
	res2, err2 := client.Do(req2)
	CheckErr(err2)

	req3, err := http.NewRequest("GET", ts.URL + "/hello", nil)
	CheckErr(err)
	req1.SetBasicAuth("3sky","test1")
	res3, err2 := client.Do(req3)
	CheckErr(err2)

	req4, err := http.NewRequest("GET", ts.URL + "/hello", nil)
	CheckErr(err)
	req1.SetBasicAuth("","")
	res4, err2 := client.Do(req4)
	CheckErr(err2)

	assert.Equal(t, 200,  res1.StatusCode)
	assert.Equal(t, 401,  res2.StatusCode)
	assert.Equal(t, 401,  res3.StatusCode)
	assert.Equal(t, 401,  res4.StatusCode)


}