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

// Similar struct to db.StatusStruct, but without model
type AppStatusStruct struct {
	ID uint 
	APP_NAME string `json:"app_name"` 
	APP_VERSION string `json:"app_version"`
	UPDATE_DATE string`json:"updated_date"`
	UPDATE_BY string `json:"updated_by"`
}

// All app struct
type AllApp struct {
	Name string
	App []AppStatusStruct
}



//SayHello seyhello
func SayHello(w http.ResponseWriter, r *http.Request) {
	h := &HelloStruct{ID: 1, SAY: "Hello Kuba"}
	json.NewEncoder(w).Encode(h)
}

// Display App status By ID (GET)
func (a *App) DisplaAppByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	i, _ := strconv.Atoi(id)
	all_ids := a.GetAllID()

	if Ifcontains(all_ids, i) {
		var app interface{};
		tmp := a.SelectFromDBWhereID(int64(i))
		app = GetAppStatusStructFromStatusStruct(&tmp)
		json.NewEncoder(w).Encode(app)
	} else {
		h := &HelloStruct{ID: i, SAY: "No app with given ID"}
		json.NewEncoder(w).Encode(h)
	}

}

// Add New App (POST)
func (a *App) AddNewApp(w http.ResponseWriter, r *http.Request) {

	var app AppStatusStruct
	var updater string
	
	_ = json.NewDecoder(r.Body).Decode(&app)
	//fmt.Printf("%+v\n", app)

	if ( len(app.APP_NAME) != 0 ) && ( len(app.APP_VERSION) != 0 ){
		
		if len(app.UPDATE_BY) == 0 {
			updater = "random guy"
		} else {
			updater = app.UPDATE_BY
		}

		a.InsertToDB(app.APP_NAME, app.APP_VERSION, updater)

	} else {
		h := &HelloStruct{SAY: "Application name and version are mandatory ! "}
		json.NewEncoder(w).Encode(h)
	}
}

// Update Data (PUT)
func (a *App) UpdateData(w http.ResponseWriter, r *http.Request) {

	var app AppStatusStruct

	vars := mux.Vars(r)
	id := vars["id"]
	i, _ := strconv.Atoi(id)
	all_ids := a.GetAllID()

	if Ifcontains(all_ids, i) {
		_ = json.NewDecoder(r.Body).Decode(&app)
		
		if len(app.APP_NAME) > 0 {
			a.UpdateSelectedColumn(int64(i), "app_name", app.APP_NAME)
		}
		
		if len(app.UPDATE_BY) > 0 {
			a.UpdateSelectedColumn(int64(i), "updated_by", app.UPDATE_BY)
		} 

		if len(app.APP_VERSION) > 0 {
			a.UpdateSelectedColumn(int64(i), "app_version", app.APP_VERSION)
		} 
		
		var app_after_update interface{};
		tmp := a.SelectFromDBWhereID(int64(i))
		app_after_update = GetAppStatusStructFromStatusStruct(&tmp)
		json.NewEncoder(w).Encode(app_after_update)

		
	} else {
		h := &HelloStruct{ID: 1, SAY: "No app with given ID"}
		json.NewEncoder(w).Encode(h)
	}

	
}


// Delete data (DELETE)
func (a *App) DeleteData(w http.ResponseWriter, r *http.Request) {
	
	var app AppStatusStruct

	vars := mux.Vars(r)
	id := vars["id"]
	i, _ := strconv.Atoi(id)
	all_ids := a.GetAllID()

	if Ifcontains(all_ids, i) {
		_ = json.NewDecoder(r.Body).Decode(&app)
		a.DeleteRowByID(int64(i))
		h := &HelloStruct{ID: i, SAY: "Record was deleted successfully !"}
		json.NewEncoder(w).Encode(h)
	} else {
		h := &HelloStruct{ID: i, SAY: "No app with given ID"}
		json.NewEncoder(w).Encode(h)
	}
}

// Display all app (GET)
func (a *App) DisplayAllApp(w http.ResponseWriter, r *http.Request) {
	
	var app *AppStatusStruct

	Apps := []AppStatusStruct{}

	data := AllApp{
		Name: "Apps",
		App: Apps,
	}

	all_ids := a.GetAllID()
	for _, k := range all_ids {
		tmp := a.SelectFromDBWhereID(int64(k))
		app = GetAppStatusStructFromStatusStruct(&tmp)
		data.AddApp(app)
	}
	json.NewEncoder(w).Encode(data)
}

func (aa *AllApp) AddApp(app *AppStatusStruct) []AppStatusStruct{
	valueOffApp := *app
	aa.App = append(aa.App, valueOffApp)
	return aa.App
}

// Convert StatusStruct to AppStruct
func GetAppStatusStructFromStatusStruct(s *StatusStruct) *AppStatusStruct {
	s2 := &AppStatusStruct{
		ID: s.Model.ID,
		APP_NAME: s.APP_NAME, 
		APP_VERSION: s.APP_VERSION, 
		UPDATE_DATE: s.UPDATE_DATE.Format("2006-01-02 15:04:05"), 
		UPDATE_BY: s.UPDATE_BY,
	}
	return s2
}

// check if list contain ID
func Ifcontains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}