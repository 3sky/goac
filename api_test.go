package main

import (
	"testing"
	"time"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/gorilla/mux"
)


func TestIfcontains(t *testing.T) {

	t1 := []int{1, 2, 3, 4}
	t2 := []int{2, 3, 90, 12}

	assert.True(t, Ifcontains(t1, 1))
	assert.False(t, Ifcontains(t1, 10))

	assert.True(t, Ifcontains(t2, 2))
	assert.False(t, Ifcontains(t1, 13))
}

func TestGetAppStatusStructFromStatusStruct(t *testing.T) {
	
	h := &StatusStruct{
		Model: gorm.Model{ID: 1},
		APP_NAME: "Test",
		APP_VERSION: "1",
		UPDATE_DATE: time.Now(),
		UPDATE_BY: "Admin 1",
		ENVIRONMENT: "dev",
	    BRANCH: "testing",
		
	}

	app := GetAppStatusStructFromStatusStruct(h)

	assert.Equal(t, h.APP_NAME, app.APP_NAME)
	assert.Equal(t, h.APP_VERSION, app.APP_VERSION)
	assert.Equal(t, h.UPDATE_BY, app.UPDATE_BY)
	assert.Equal(t, h.ENVIRONMENT, app.ENVIRONMENT)
	assert.Equal(t, h.BRANCH, app.BRANCH)

}

func TestSayHello(t *testing.T) {

	var h HelloStruct

	req := httptest.NewRequest("GET", "http://127.0.0.1:5000/hello", nil)
	w := httptest.NewRecorder()

	SayHello(w, req)

	resp := w.Result()
	
	err := json.NewDecoder(resp.Body).Decode(&h)
	if err != nil {
		fmt.Printf("Error while decode in TestSayHello: %v", err)
	}
	assert.Equal(t, 200,          resp.StatusCode)
	assert.Equal(t, 1,            h.ID)
	assert.Equal(t, "Hello Kuba", h.SAY)
		
}

func TestDisplayAppByID(t *testing.T) {
	
	a := createTestDBConnection()
	defer a.DB.Close()

	var h HelloStruct
	var a1, a2 AppStatusStruct

	err := a.MakeMigration()
	if err != nil {
		fmt.Printf("Error with insert into DB in TestDisplayAppByID: %v", err)
	}
	
	err = a.InsertToDB("Test_run_app_1", "1", "UnitTest_1", "dev", "fix_12")
	if err != nil {
		fmt.Printf("Error with insert into DB in TestDisplayAppByID: %v", err)
	}

	err = a.InsertToDB("Test_run_app_2", "2", "UnitTest_2", "stage", "new_feature")
	if err != nil {
		fmt.Printf("Error with insert into DB in TestDisplayAppByID: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/app/{id}", a.DisplayAppByID).Methods("GET")
	
    ts := httptest.NewServer(r)
    defer ts.Close()


	url1 := ts.URL + "/api/app/" + "1"
	url2 := ts.URL + "/api/app/" + "2"
	url3 := ts.URL + "/api/app/" + "3"

	resp1, err := http.Get(url1)
	if err != nil {
		fmt.Printf("Error while making GET Request: %v", err)
	}

	resp2, err := http.Get(url2)
	if err != nil {
		fmt.Printf("Error while making GET Request: %v", err)
	}

	resp3, err := http.Get(url3)
	if err != nil {
		fmt.Printf("Error while making GET Request: %v", err)
	}

	err = json.NewDecoder(resp1.Body).Decode(&a1)
	if err != nil {
		fmt.Printf("Error while decode in TestDisplayAppByID: %v", err)
	}

	err = json.NewDecoder(resp2.Body).Decode(&a2)
	if err != nil {
		fmt.Printf("Error while decode in TestDisplayAppByID: %v", err)
	}

	err = json.NewDecoder(resp3.Body).Decode(&h)
	if err != nil {
		fmt.Printf("Error while decode in TestDisplayAppByID: %v", err)
	}


	assert.Equal(t, 200,              resp1.StatusCode)
	assert.Equal(t, 1,                int(a1.ID))
	assert.Equal(t, "Test_run_app_1", string(a1.APP_NAME) )
	assert.Equal(t, "1",              a1.APP_VERSION)
	assert.Equal(t, "UnitTest_1",     a1.UPDATE_BY )
	assert.Equal(t, "dev",            a1.ENVIRONMENT)
	assert.Equal(t, "fix_12",         a1.BRANCH )

	assert.Equal(t, 200,              resp2.StatusCode)
	assert.Equal(t, 2,                int(a2.ID))
	assert.Equal(t, "Test_run_app_2", string(a2.APP_NAME))
	assert.Equal(t, "2",              a2.APP_VERSION)
	assert.Equal(t, "UnitTest_2",     a2.UPDATE_BY)
	assert.Equal(t, "stage",          a2.ENVIRONMENT)
	assert.Equal(t, "new_feature",    a2.BRANCH )

	assert.Equal(t, 200,                    resp3.StatusCode)
	assert.Equal(t, "No app with given ID", h.SAY)

}

func TestAddNewApp(t *testing.T) {
	
	a := createTestDBConnection()
	defer a.DB.Close()

	var st1, st2 StatusStruct
	var h HelloStruct

    P1 := &AppStatusStruct{
		APP_NAME: "New_app",
		APP_VERSION: "1.01",
		UPDATE_BY:  "test1",
		ENVIRONMENT: "dev",
		BRANCH: "testing",
	}

	P2 := &AppStatusStruct{
		APP_NAME: "New_app_2",
		APP_VERSION: "11.1",
	}

	P3 := &AppStatusStruct{
		APP_NAME: "New_app_2",
		UPDATE_BY: "test3",
	}

	jsonP1, err := json.Marshal(P1)
	if err != nil {
		fmt.Printf("Error while marshall in TestAddNewApp: %v", err)
	}

	jsonP2, err := json.Marshal(P2)
	if err != nil {
		fmt.Printf("Error while marshall in TestAddNewApp: %v", err)
	}

	jsonP3, err := json.Marshal(P3)
	if err != nil {
		fmt.Printf("Error while marshall in TestAddNewApp: %v", err)
	}


	req1, err := http.NewRequest("POST", "http://127.0.0.1:5000/api/app/new", bytes.NewBuffer(jsonP1))
	if err != nil {
		fmt.Printf("Error while make new POST Request: %v", err)
	}

	req2, err := http.NewRequest("POST", "http://127.0.0.1:5000/api/app/new", bytes.NewBuffer(jsonP2))
	if err != nil {
		fmt.Printf("Error while make new POST Request: %v", err)
	}

	req3, err := http.NewRequest("POST", "http://127.0.0.1:5000/api/app/new", bytes.NewBuffer(jsonP3))
	if err != nil {
		fmt.Printf("Error while make new POST Request: %v", err)
	}
	
	res1 := httptest.NewRecorder()
	res2 := httptest.NewRecorder()
	res3 := httptest.NewRecorder()

	a.AddNewApp(res1, req1)
	a.AddNewApp(res2, req2)
	a.AddNewApp(res3, req3)

	st1, err = a.SelectFromDBWhereID(int64(3))
	if err != nil {
		fmt.Printf("Error with get row in TestAddNewApp: %v", err)
	}

	st2, err = a.SelectFromDBWhereID(int64(4))
	if err != nil {
		fmt.Printf("Error with get row in TestAddNewApp: %v", err)
	}

	err = json.NewDecoder(res3.Body).Decode(&h)
	if err != nil {
		fmt.Printf("Error while decode in TestAddNewApp: %v", err)
	}

	assert.Equal(t, 200,       res1.Code )
	assert.Equal(t, 3,         int(st1.Model.ID))
	assert.Equal(t, "New_app", st1.APP_NAME)
	assert.Equal(t, "1.01",    st1.APP_VERSION)
	assert.Equal(t, "test1",   st1.UPDATE_BY)
	assert.Equal(t, "dev",     st1.ENVIRONMENT)
	assert.Equal(t, "testing", st1.BRANCH )

	assert.Equal(t, 200,          res2.Code)
	assert.Equal(t, 4,            int(st2.Model.ID))
	assert.Equal(t, "New_app_2",  st2.APP_NAME)
	assert.Equal(t, "11.1",       st2.APP_VERSION)
	assert.Equal(t, "random guy", st2.UPDATE_BY)
	assert.Equal(t, "",           st2.ENVIRONMENT)
	assert.Equal(t, "",           st2.BRANCH )

	assert.Equal(t, 200,                                             res3.Code  )
	assert.Equal(t, "Application name and version are mandatory ! ", h.SAY)

	
}

func TestDisplayAllApp(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	var allApp AllApp
	var app1, app2, app3, app4 AppStatusStruct

	req := httptest.NewRequest("GET", "http://127.0.0.1:5000/api/apps", nil)
	w := httptest.NewRecorder()

	a.DisplayAllApp(w, req)

	resp := w.Result()
	
	
	err := json.NewDecoder(resp.Body).Decode(&allApp)
	if err != nil {
		fmt.Printf("Error while decode in TestDisplayAllApp: %v", err)
	}

	app1 = allApp.App[0]
	app2 = allApp.App[1]
	app3 = allApp.App[2]
	app4 = allApp.App[3]

	
	assert.Equal(t, 1,                int(app1.ID))
	assert.Equal(t, "Test_run_app_1", app1.APP_NAME)
	assert.Equal(t, "1",              app1.APP_VERSION)
	assert.Equal(t, "UnitTest_1",     app1.UPDATE_BY)

	assert.Equal(t, 2,                int(app2.ID))
	assert.Equal(t, "Test_run_app_2", app2.APP_NAME)
	assert.Equal(t, "2",              app2.APP_VERSION)
	assert.Equal(t, "UnitTest_2",     app2.UPDATE_BY)

	assert.Equal(t, 3,         int(app3.ID))
	assert.Equal(t, "New_app", app3.APP_NAME)
	assert.Equal(t, "1.01",    app3.APP_VERSION)
	assert.Equal(t, "test1",   app3.UPDATE_BY)
 
	assert.Equal(t, 4,            int(app4.ID))
	assert.Equal(t, "New_app_2",  app4.APP_NAME)
	assert.Equal(t, "11.1",       app4.APP_VERSION)
	assert.Equal(t, "random guy", app4.UPDATE_BY)


}




/** at this moment test doesn't work PUT Method test
func TestDeleteData(t *testing.T) {

	
	var h1 HelloStruct

	r := mux.NewRouter()
	r.HandleFunc("/api/app/{id}", DisplaAppByID).Methods("DELETE")
	
    ts := httptest.NewServer(r)
    defer ts.Close()

	req, err := http.NewRequest("DELETE", "1", nil)
	w := httptest.NewRecorder()

	DeleteData(w, req)

	resp := w.Result()
	
	
	_ = json.NewDecoder(resp.Body).Decode(&h1)


	fmt.Println(h1)

}
**/


/** at this moment test doesn't work PUT Method test
func TestUpdateData(t * testing.T) {

	//var st1web, st1db db.StatusStruct
	var st1web db.StatusStruct
	// var st2web, st2db db.StatusStruct
	// var st3db db.StatusStruct
	// var st3web HelloStruct

	//p2 := []byte(`{"app_name": "2.0" }`)
	p1 := []byte(`{"app_name": "Not_test_now", "app_version": "10.12"`)
	// p3 := []byte(`{"app_name": "API Proxy Newier", "app_version": "1.2", "updated_by": "Kuba 2" }`)
	
	r := mux.NewRouter()
	ts := httptest.NewServer(r)
	client := &http.Client{}
	url1 := ts.URL + "/api/app/1"
	req1, _ := http.NewRequest("PUT", url1, bytes.NewBuffer(p1))
	
	

	defer ts.Close()
	
	res1, err := client.Do(req1)


    //eq1, _ := http.NewRequest("PUT", "http://127.0.0.1:5000/api/app/1", bytes.NewBuffer(p1))
	//res1 := httptest.NewRecorder()	
	//handler := http.HandlerFunc(UpdateData)

	//handler.ServeHTTP(res1, req1)

	fmt.Printf("from DB - %+v\n", res1)
	fmt.Printf("from DB - %+v\n", req1)
	//UpdateData(res1, req1)


	fmt.Printf("from DB - %+v\n", res1.Body)

    //json.Unmarshal(res1.Body.Bytes(), &st1web)
	// _ = json.NewDecoder(res1.Body).Decode(&st1web)
	fmt.Printf("from WEB - %+v\n", st1web)


	// jsonP1, _ := json.Marshal(P1)
	// jsonP2, _ := json.Marshal(P2)
	// jsonP3, _ := json.Marshal(P3)

	// req1, _ := http.NewRequest("PUT", "http://127.0.0.1:5000/api/app/1", bytes.NewBuffer(jsonP1))

	// req2, _ := http.NewRequest("PUT", "http://127.0.0.1:5000/api/app/2", bytes.NewBuffer(jsonP2))
	// req3, _ := http.NewRequest("PUT", "http://127.0.0.1:5000/api/app/3", bytes.NewBuffer(jsonP3))

	// w1 := httptest.NewRecorder()
	// w2 := httptest.NewRecorder()
	// w3 := httptest.NewRecorder()

	// UpdateData(w1, req1)
	// UpdateData(w2, req2)
	// UpdateData(w3, req3)

	// res1 := w1.Result()
	// res2 := w2.Result()
	// res3 := w3.Result()
	
	// _ = json.NewDecoder(res1.Body).Decode(&st1web)
	// _ = json.NewDecoder(res2.Body).Decode(&st2web)
	// _ = json.NewDecoder(res3.Body).Decode(&st3web)

	// st1db = db.SelectFromDBWhereID(int64(1))
	// st2db = db.SelectFromDBWhereID(int64(2))

	// fmt.Printf("from DB - %+v\n", st1db)
	// //fmt.Printf("from DB - %+v\n", st2db)
	// //fmt.Printf("from DB - %+v\n", st3db)
	// fmt.Printf("from WEB - %+v\n", st1web)
	// // fmt.Printf("%+v\n", st2web)
	// // fmt.Printf("%+v\n", st3web)
	// //Case 1
	// assert.Equal(t, res1.StatusCode, 200)
	// assert.Equal(t, int(st1web.Model.ID), 1)
	// assert.Equal(t, st1web.Model.ID, st1db.Model.ID)
	// assert.Equal(t, st1web.APP_VERSION, "2.0")
	// assert.Equal(t, st1web.APP_VERSION, st1db.APP_VERSION)
	// assert.Equal(t, st1web.APP_NAME, st1db.APP_NAME)
	
	// //Case 2
	// assert.Equal(t, res2.StatusCode, 200)
	// assert.Equal(t, int(st2web.Model.ID), 2)
	// assert.Equal(t, st2web.Model.ID, st2db.Model.ID)
	// assert.Equal(t, st2web.APP_VERSION, "10.12")
	// assert.Equal(t, st2web.APP_VERSION, st2db.APP_VERSION)
	// assert.Equal(t, st2web.APP_NAME, st2db.APP_NAME)
	// assert.Equal(t, st2web.APP_NAME, "Not_test_now")

	// //Case 3
	// assert.Equal(t, res3.StatusCode, 200)
	// assert.Equal(t, int(st3db.Model.ID), 3)
	// assert.Equal(t, st3web.SAY, "You cannot update UPDATE_DATE")
}
**/
