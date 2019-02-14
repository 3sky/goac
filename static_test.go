package main

import (
	"testing"
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"net/http/httptest"
	

	"github.com/stretchr/testify/assert"
	"github.com/gorilla/mux"
)


func TestDisplayHtml(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	var bodyString string
	err := a.MakeMigration()
	if err != nil {
		fmt.Printf("Error with migration in TestDisplayHtml: %v", err)
	}

	err = a.InsertToDB("Test_run_app_1", "1", "UnitTest_1", "dev", "hotfix1")
	if err != nil {
		fmt.Printf("Error while insert data in TestDisplayHtml: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", a.DisplayHtml).Methods("GET")
	
    ts := httptest.NewServer(r)
	defer ts.Close()
	
	url := ts.URL + "/"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Error while make GET Request: %v", err)
	}


	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		bodyString = string(bodyBytes)
	}

	assert.Contains(t, bodyString, "Hello There!" )
	assert.Contains(t, bodyString, "Test_run_app_1" )
	assert.Contains(t, bodyString, "<td>1 </td>" )
	assert.Contains(t, bodyString, "UnitTest_1" )

	err = os.Remove("TestDB.db")
	if err != nil {
		fmt.Printf("Error while trying remove TestDB: %v", err)
	}
		
}