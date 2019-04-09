package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestDisplayHTMLDev(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	var bodyString string
	err := a.MakeMigration()
	if err != nil {
		fmt.Printf("Error with migration in TestDisplayHtml: %v", err)
	}

	err = a.InsertToDB("Test_run_app_1", "1", "UnitTest_1", "dev", "hotfix1", "")
	if err != nil {
		fmt.Printf("Error while insert data in TestDisplayHtml: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/dev", a.DisplayHTMLDev).Methods("GET")

	ts := httptest.NewServer(r)
	defer ts.Close()

	url := ts.URL + "/dev"

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

	assert.Contains(t, bodyString, "Hello There!")
	assert.Contains(t, bodyString, "Test_run_app_1")
	assert.Contains(t, bodyString, "<td>1</td>")
	assert.Contains(t, bodyString, "<td>dev</td>")
	assert.Contains(t, bodyString, "UnitTest_1")

	err = os.Remove("TestDB.db")
	if err != nil {
		fmt.Printf("Error while trying remove TestDB: %v", err)
	}

}

func TestDisplayHTMLStg(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	var bodyString string
	err := a.MakeMigration()
	if err != nil {
		fmt.Printf("Error with migration in TestDisplayHtml: %v", err)
	}

	err = a.InsertToDB("Test_run_app_1", "1", "UnitTest_1", "stg", "hotfix1", "")
	if err != nil {
		fmt.Printf("Error while insert data in TestDisplayHtml: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/stg", a.DisplayHTMLStg).Methods("GET")

	ts := httptest.NewServer(r)
	defer ts.Close()

	url := ts.URL + "/stg"

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

	assert.Contains(t, bodyString, "Hello There!")
	assert.Contains(t, bodyString, "Test_run_app_1")
	assert.Contains(t, bodyString, "<td>1</td>")
	assert.Contains(t, bodyString, "<td>stg</td>")
	assert.Contains(t, bodyString, "UnitTest_1")

	err = os.Remove("TestDB.db")
	if err != nil {
		fmt.Printf("Error while trying remove TestDB: %v", err)
	}

}
