package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
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
		Model:       gorm.Model{ID: 1},
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

	_ = json.NewDecoder(resp.Body).Decode(&h)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, 1, h.ID)
	assert.Equal(t, "Hello Kuba", h.SAY)

}

func TestDisplaAppByID(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	var h HelloStruct
	var a1, a2 AppStatusStruct

	a.MakeMigration()
	a.InsertToDB("Test_run_app_1", "1", "UnitTest_1", "dev", "fix_12")
	a.InsertToDB("Test_run_app_2", "2", "UnitTest_2", "stage", "new_feature")

	r := mux.NewRouter()
	r.HandleFunc("/api/app/{id}", a.DisplaAppByID).Methods("GET")

	ts := httptest.NewServer(r)
	defer ts.Close()

	url1 := ts.URL + "/api/app/" + "1"
	url2 := ts.URL + "/api/app/" + "2"
	url3 := ts.URL + "/api/app/" + "3"

	resp1, err := http.Get(url1)
	CheckErr(err)
	resp2, err := http.Get(url2)
	CheckErr(err)
	resp3, err := http.Get(url3)
	CheckErr(err)

	_ = json.NewDecoder(resp1.Body).Decode(&a1)
	_ = json.NewDecoder(resp2.Body).Decode(&a2)
	_ = json.NewDecoder(resp3.Body).Decode(&h)

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
	assert.Equal(t, "No app with given ID", h.SAY)

}

func TestAddNewApp(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
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

	jsonP1, _ := json.Marshal(P1)
	jsonP2, _ := json.Marshal(P2)
	jsonP3, _ := json.Marshal(P3)

	req1, _ := http.NewRequest("POST", "http://127.0.0.1:5000/api/app/new", bytes.NewBuffer(jsonP1))
	req2, _ := http.NewRequest("POST", "http://127.0.0.1:5000/api/app/new", bytes.NewBuffer(jsonP2))
	req3, _ := http.NewRequest("POST", "http://127.0.0.1:5000/api/app/new", bytes.NewBuffer(jsonP3))

	res1 := httptest.NewRecorder()
	res2 := httptest.NewRecorder()
	res3 := httptest.NewRecorder()

	a.AddNewApp(res1, req1)
	a.AddNewApp(res2, req2)
	a.AddNewApp(res3, req3)

	st1 = a.SelectFromDBWhereID(int64(3))
	st2 = a.SelectFromDBWhereID(int64(4))

	_ = json.NewDecoder(res3.Body).Decode(&h)

	assert.Equal(t, 200, res1.Code)
	assert.Equal(t, 3, int(st1.Model.ID))
	assert.Equal(t, "New_app", st1.AppName)
	assert.Equal(t, "1.01", st1.AppVersion)
	assert.Equal(t, "test1", st1.UpdateBy)
	assert.Equal(t, "dev", st1.Environment)
	assert.Equal(t, "testing", st1.Branch)

	assert.Equal(t, 200, res2.Code)
	assert.Equal(t, 4, int(st2.Model.ID))
	assert.Equal(t, "New_app_2", st2.AppName)
	assert.Equal(t, "11.1", st2.AppVersion)
	assert.Equal(t, "random guy", st2.UpdateBy)
	assert.Equal(t, "", st2.Environment)
	assert.Equal(t, "", st2.Branch)

	assert.Equal(t, 200, res3.Code)
	assert.Equal(t, "Application name and version are mandatory ! ", h.SAY)

}

func TestDisplayAllApp(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	var allApp AllApp
	var app1, app2, app3, app4 AppStatusStruct

	req := httptest.NewRequest("GET", "http://127.0.0.1:5000/api/apps", nil)
	w := httptest.NewRecorder()

	a.DisplayAllApp(w, req)

	resp := w.Result()

	_ = json.NewDecoder(resp.Body).Decode(&allApp)

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

	//p2 := []byte(`{"AppName": "2.0" }`)
	p1 := []byte(`{"AppName": "Not_test_now", "AppVersion": "10.12"`)
	// p3 := []byte(`{"AppName": "API Proxy Newier", "AppVersion": "1.2", "updated_by": "Kuba 2" }`)

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
	// assert.Equal(t, st1web.AppVersion, "2.0")
	// assert.Equal(t, st1web.AppVersion, st1db.AppVersion)
	// assert.Equal(t, st1web.AppName, st1db.AppName)

	// //Case 2
	// assert.Equal(t, res2.StatusCode, 200)
	// assert.Equal(t, int(st2web.Model.ID), 2)
	// assert.Equal(t, st2web.Model.ID, st2db.Model.ID)
	// assert.Equal(t, st2web.AppVersion, "10.12")
	// assert.Equal(t, st2web.AppVersion, st2db.AppVersion)
	// assert.Equal(t, st2web.AppName, st2db.AppName)
	// assert.Equal(t, st2web.AppName, "Not_test_now")

	// //Case 3
	// assert.Equal(t, res3.StatusCode, 200)
	// assert.Equal(t, int(st3db.Model.ID), 3)
	// assert.Equal(t, st3web.SAY, "You cannot update UpdateDate")
}
**/
