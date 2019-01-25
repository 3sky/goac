package static

import (
	"p2go/db"
	"testing"
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
	"net/http/httptest"
	

	"github.com/stretchr/testify/assert"
	"github.com/gorilla/mux"
)

/** 
At this moment still not work
**/

func TestDisplayHtml(t *testing.T) {

	var bodyString string
	db.MakeMIgration()
	db.InsertToDB("Test_run_app_1", "1", "UnitTest_1")

	r := mux.NewRouter()
	r.HandleFunc("/", DisplayHtml).Methods("GET")
	
    ts := httptest.NewServer(r)
	defer ts.Close()
	
	url := ts.URL + "/"

	resp, err := http.Get(url)
	db.CheckErr(err)


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

	err = os.Remove("./SimpleDB.db")
	db.CheckErr(err)
		
}