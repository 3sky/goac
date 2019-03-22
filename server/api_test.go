package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
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
		ID:          1,
		AppName:     "Test",
		AppVersion:  "1",
		UpdateDate:  time.Now(),
		UpdateBy:    "Admin 1",
		Environment: "dev",
		Branch:      "testing",
	}

	app := GetAppStatusStructFromStatusStruct(h)

	assert.Equal(t, h.AppName, app.AppName)
	assert.Equal(t, h.AppVersion, app.AppVersion)
	assert.Equal(t, h.UpdateBy, app.UpdateBy)
	assert.Equal(t, h.Environment, app.Environment)
	assert.Equal(t, h.Branch, app.Branch)

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
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, 1, h.ID)
	assert.Equal(t, "Hello Kuba", h.Say)

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

	assert.Equal(t, 200, resp1.StatusCode)
	assert.Equal(t, 1, int(a1.ID))
	assert.Equal(t, "Test_run_app_1", string(a1.AppName))
	assert.Equal(t, "1", a1.AppVersion)
	assert.Equal(t, "UnitTest_1", a1.UpdateBy)
	assert.Equal(t, "dev", a1.Environment)
	assert.Equal(t, "fix_12", a1.Branch)

	assert.Equal(t, 200, resp2.StatusCode)
	assert.Equal(t, 2, int(a2.ID))
	assert.Equal(t, "Test_run_app_2", string(a2.AppName))
	assert.Equal(t, "2", a2.AppVersion)
	assert.Equal(t, "UnitTest_2", a2.UpdateBy)
	assert.Equal(t, "stage", a2.Environment)
	assert.Equal(t, "new_feature", a2.Branch)

	assert.Equal(t, 200, resp3.StatusCode)
	assert.Equal(t, "No app with given ID", h.Say)

}

func TestAddNewApp(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	var st1, st2 StatusStruct
	var h HelloStruct

	P1 := &AppStatusStruct{
		AppName:     "New_app",
		AppVersion:  "1.01",
		UpdateBy:    "test1",
		Environment: "dev",
		Branch:      "testing",
	}

	P2 := &AppStatusStruct{
		AppName:    "New_app_2",
		AppVersion: "11.1",
	}

	P3 := &AppStatusStruct{
		AppName:  "New_app_2",
		UpdateBy: "test3",
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

	assert.Equal(t, 200, res1.Code)
	assert.Equal(t, 3, st1.ID)
	assert.Equal(t, "New_app", st1.AppName)
	assert.Equal(t, "1.01", st1.AppVersion)
	assert.Equal(t, "test1", st1.UpdateBy)
	assert.Equal(t, "dev", st1.Environment)
	assert.Equal(t, "testing", st1.Branch)

	assert.Equal(t, 200, res2.Code)
	assert.Equal(t, 4, st2.ID)
	assert.Equal(t, "New_app_2", st2.AppName)
	assert.Equal(t, "11.1", st2.AppVersion)
	assert.Equal(t, "random guy", st2.UpdateBy)
	assert.Equal(t, "", st2.Environment)
	assert.Equal(t, "", st2.Branch)

	assert.Equal(t, 200, res3.Code)
	assert.Equal(t, "Application name and version are mandatory ! ", h.Say)

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

	assert.Equal(t, 1, int(app1.ID))
	assert.Equal(t, "Test_run_app_1", app1.AppName)
	assert.Equal(t, "1", app1.AppVersion)
	assert.Equal(t, "UnitTest_1", app1.UpdateBy)

	assert.Equal(t, 2, int(app2.ID))
	assert.Equal(t, "Test_run_app_2", app2.AppName)
	assert.Equal(t, "2", app2.AppVersion)
	assert.Equal(t, "UnitTest_2", app2.UpdateBy)

	assert.Equal(t, 3, int(app3.ID))
	assert.Equal(t, "New_app", app3.AppName)
	assert.Equal(t, "1.01", app3.AppVersion)
	assert.Equal(t, "test1", app3.UpdateBy)

	assert.Equal(t, 4, int(app4.ID))
	assert.Equal(t, "New_app_2", app4.AppName)
	assert.Equal(t, "11.1", app4.AppVersion)
	assert.Equal(t, "random guy", app4.UpdateBy)

}

func TestUpdateData(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	var st1, st2 StatusStruct
	U1 := &AppStatusStruct{
		AppName:  "SomeApp",
		UpdateBy: "test3",
	}

	U2 := &AppStatusStruct{
		AppVersion:  "1.11",
		Environment: "stg",
	}

	jsonU1, err := json.Marshal(U1)
	if err != nil {
		fmt.Printf("Error while marshall in TestUpdateData: %v", err)
	}

	jsonU2, err := json.Marshal(U2)
	if err != nil {
		fmt.Printf("Error while marshall in TestUpdateData: %v", err)
	}

	st1, err = a.SelectFromDBWhereID(int64(1))
	if err != nil {
		fmt.Printf("Error with get row in TestUpdateData: %v \n", err)
	}

	st2, err = a.SelectFromDBWhereID(int64(2))
	if err != nil {
		fmt.Printf("Error with get row in TestUpdateData: %v \n", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/app/{id}", a.UpdateData).Methods("PUT")
	client := &http.Client{}
	ts := httptest.NewServer(r)

	defer ts.Close()

	url1 := ts.URL + "/api/app/" + "1"
	url2 := ts.URL + "/api/app/" + "2"

	req1, err := http.NewRequest("PUT", url1, bytes.NewBuffer(jsonU1))
	req1.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("Error while making PUT Req in TestUpdateData: %v \n", err)
	}

	req2, err := http.NewRequest("PUT", url2, bytes.NewBuffer(jsonU2))
	req2.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("Error while making PUT Req in TestUpdateData: %v \n", err)
	}

	res1, err := client.Do(req1)
	if err != nil {
		fmt.Printf("Error while making PUT Req in TestUpdateData: %v \n", err)
	}
	defer res1.Body.Close()

	res2, err := client.Do(req2)
	if err != nil {
		fmt.Printf("Error while making PUT Req in TestUpdateData: %v \n", err)
	}
	defer res2.Body.Close()

	st1, err = a.SelectFromDBWhereID(int64(1))
	if err != nil {
		fmt.Printf("Error with get row in TestUpdateData: %v \n", err)
	}

	st2, err = a.SelectFromDBWhereID(int64(2))
	if err != nil {
		fmt.Printf("Error with get row in TestUpdateData: %v \n", err)
	}

	assert.Equal(t, 1, int(st1.ID))
	assert.Equal(t, "SomeApp", st1.AppName)
	assert.Equal(t, "1", st1.AppVersion)
	assert.Equal(t, "dev", st1.Environment)
	assert.Equal(t, "test3", st1.UpdateBy)

	assert.Equal(t, 2, int(st2.ID))
	assert.Equal(t, "Test_run_app_2", st2.AppName)
	assert.Equal(t, "1.11", st2.AppVersion)
	assert.Equal(t, "stg", st2.Environment)
	assert.Equal(t, "UnitTest_2", st2.UpdateBy)

}

func TestDeleteData(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/app/{id}", a.DeleteData).Methods("DELETE")
	client := &http.Client{}
	ts := httptest.NewServer(r)

	defer ts.Close()

	url1 := ts.URL + "/api/app/" + "1"

	req1, err := http.NewRequest("DELETE", url1, nil)
	if err != nil {
		fmt.Printf("Error while making PUT Req in TestUpdateData: %v \n", err)
	}

	res1, err := client.Do(req1)
	if err != nil {
		fmt.Printf("Error while making PUT Req in TestUpdateData: %v \n", err)
	}
	defer res1.Body.Close()

	_, err = a.SelectFromDBWhereID(int64(1))

	assert.Contains(t, err.Error(), "record not found")
}
