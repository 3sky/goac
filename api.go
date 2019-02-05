package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//HelloStruct for valid JSON
type HelloStruct struct {
	ID  int    `json:"ID,omitempty"`
	SAY string `json:"INFO,omitempty"`
}

// AppStatusStruct is similar struct to db.StatusStruct, but without model
type AppStatusStruct struct {
	ID          uint
	AppName     string `json:"AppName"`
	AppVersion  string `json:"AppVersion"`
	Environment string `json:"Env"`
	Branch      string `json:"Branch"`
	UpdateDate  string `json:"UpdateDate"`
	UpdateBy    string `json:"UpdateBy"`
}

// AllApp struct
type AllApp struct {
	Name string
	App  []AppStatusStruct
}

//SayHello seyhello
func SayHello(w http.ResponseWriter, r *http.Request) {
	h := &HelloStruct{ID: 1, SAY: "Hello Kuba"}
	json.NewEncoder(w).Encode(h)
}

//DisplaAppByID - Display App status By ID (GET)
func (a *App) DisplaAppByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	i, _ := strconv.Atoi(id)
	allIds := a.GetAllID()

	if Ifcontains(allIds, i) {
		var app interface{}
		tmp := a.SelectFromDBWhereID(int64(i))
		app = GetAppStatusStructFromStatusStruct(&tmp)
		json.NewEncoder(w).Encode(app)
	} else {
		h := &HelloStruct{ID: i, SAY: "No app with given ID"}
		json.NewEncoder(w).Encode(h)
	}

}

// AddNewApp (POST)
func (a *App) AddNewApp(w http.ResponseWriter, r *http.Request) {

	var app AppStatusStruct
	var updater string

	_ = json.NewDecoder(r.Body).Decode(&app)

	if (len(app.AppName) != 0) && (len(app.AppVersion) != 0) {

		if len(app.UpdateBy) == 0 {
			updater = "random guy"
		} else {
			updater = app.UpdateBy
		}

		a.InsertToDB(app.AppName, app.AppVersion, updater, app.Environment, app.Branch)

	} else {
		h := &HelloStruct{SAY: "Application name and version are mandatory ! "}
		json.NewEncoder(w).Encode(h)
	}
}

// UpdateData (PUT)
func (a *App) UpdateData(w http.ResponseWriter, r *http.Request) {

	var app AppStatusStruct

	vars := mux.Vars(r)
	id := vars["id"]
	i, _ := strconv.Atoi(id)
	allIds := a.GetAllID()

	if Ifcontains(allIds, i) {
		_ = json.NewDecoder(r.Body).Decode(&app)

		if len(app.AppName) > 0 {
			a.UpdateSelectedColumn(int64(i), "AppName", app.AppName)
		}

		if len(app.UpdateBy) > 0 {
			a.UpdateSelectedColumn(int64(i), "UpdateBy", app.UpdateBy)
		}

		if len(app.AppVersion) > 0 {
			a.UpdateSelectedColumn(int64(i), "AppVersion", app.AppVersion)
		}

		if len(app.Environment) > 0 {
			a.UpdateSelectedColumn(int64(i), "Env", app.Environment)
		}

		if len(app.Branch) > 0 {
			a.UpdateSelectedColumn(int64(i), "Branch", app.Branch)
		}

		var appAfterUpdate interface{}
		tmp := a.SelectFromDBWhereID(int64(i))
		appAfterUpdate = GetAppStatusStructFromStatusStruct(&tmp)
		json.NewEncoder(w).Encode(appAfterUpdate)

	} else {
		h := &HelloStruct{ID: 1, SAY: "No app with given ID"}
		json.NewEncoder(w).Encode(h)
	}

}

// DeleteData (DELETE)
func (a *App) DeleteData(w http.ResponseWriter, r *http.Request) {

	var app AppStatusStruct

	vars := mux.Vars(r)
	id := vars["id"]
	i, _ := strconv.Atoi(id)
	allIds := a.GetAllID()

	if Ifcontains(allIds, i) {
		_ = json.NewDecoder(r.Body).Decode(&app)
		a.DeleteRowByID(int64(i))
		h := &HelloStruct{ID: i, SAY: "Record was deleted successfully !"}
		json.NewEncoder(w).Encode(h)
	} else {
		h := &HelloStruct{ID: i, SAY: "No app with given ID"}
		json.NewEncoder(w).Encode(h)
	}
}

// DisplayAllApp (GET)
func (a *App) DisplayAllApp(w http.ResponseWriter, r *http.Request) {

	var app *AppStatusStruct

	Apps := []AppStatusStruct{}

	data := AllApp{
		Name: "Apps",
		App:  Apps,
	}

	allIds := a.GetAllID()
	for _, k := range allIds {
		tmp := a.SelectFromDBWhereID(int64(k))
		app = GetAppStatusStructFromStatusStruct(&tmp)
		data.AddApp(app)
	}
	json.NewEncoder(w).Encode(data)
}

// AddApp add App to list of app
func (aa *AllApp) AddApp(app *AppStatusStruct) []AppStatusStruct {
	valueOffApp := *app
	aa.App = append(aa.App, valueOffApp)
	return aa.App
}

// GetAppStatusStructFromStatusStruct convert StatusStruct to AppStruct
func GetAppStatusStructFromStatusStruct(s *StatusStruct) *AppStatusStruct {
	s2 := &AppStatusStruct{
		ID:          s.Model.ID,
		AppName:     s.AppName,
		AppVersion:  s.AppVersion,
		Environment: s.Environment,
		Branch:      s.Branch,
		UpdateDate:  s.UpdateDate.Format("2006-01-02 15:04:05"),
		UpdateBy:    s.UpdateBy,
	}
	return s2
}

// Ifcontains check if list contain ID
func Ifcontains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
