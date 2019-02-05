package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestDisplayHTML(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	var bodyString string
	a.MakeMigration()
	a.InsertToDB("Test_run_app_1", "1", "UnitTest_1", "dev", "hotfix1")

	r := mux.NewRouter()
	r.HandleFunc("/", a.DisplayHTML).Methods("GET")

	ts := httptest.NewServer(r)
	defer ts.Close()

	url := ts.URL + "/"

	resp, err := http.Get(url)
	CheckErr(err)

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		bodyString = string(bodyBytes)
	}

	assert.Contains(t, bodyString, "Hello There!")
	assert.Contains(t, bodyString, "Test_run_app_1")
	assert.Contains(t, bodyString, "<td>1 </td>")
	assert.Contains(t, bodyString, "UnitTest_1")

	err = os.Remove("TestDB.db")
	CheckErr(err)

}
