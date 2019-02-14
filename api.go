package main

import (
	
	"encoding/json"
	"net/http"
	"strconv"
	"log"

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
	ENVIRONMENT string `json:"env"`
	BRANCH string `json:"branch"`
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
func (a *App) DisplayAppByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.Atoi(id)
	if err != nil {
		errorWithNonDigit(w, err)
	} else {

		all_ids, err := a.GetAllID()
		if err != nil {
			log.Printf("Cannot get all IDs from DB: %v", err)
		}

		if Ifcontains(all_ids, i) {
			var app interface{};
			tmp, err := a.SelectFromDBWhereID(int64(i))
			if err != nil {
				log.Printf("Cannot get row from DB: %v", err)
			}
			app = GetAppStatusStructFromStatusStruct(&tmp)
			json.NewEncoder(w).Encode(app)
		} else {
			h := &HelloStruct{ID: i, SAY: "No app with given ID"}
			json.NewEncoder(w).Encode(h)
		}
	}



}

// Add New App (POST)
func (a *App) AddNewApp(w http.ResponseWriter, r *http.Request) {

	var app AppStatusStruct
	var updater string
	
	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		errorWhileDecode(w, err)
	} else {

		if ( len(app.APP_NAME) != 0 ) && ( len(app.APP_VERSION) != 0 ){
		
			if len(app.UPDATE_BY) == 0 {
				updater = "random guy"
			} else {
				updater = app.UPDATE_BY
			}
	
			err = a.InsertToDB(app.APP_NAME, app.APP_VERSION, updater, app.ENVIRONMENT, app.BRANCH)
			if err != nil {
				log.Printf("Cannot insert row into DB: %v", err)
			}
	
		} else {
			h := &HelloStruct{SAY: "Application name and version are mandatory ! "}
			json.NewEncoder(w).Encode(h)
		}
	}
}

// Update Data (PUT)
func (a *App) UpdateData(w http.ResponseWriter, r *http.Request) {

	var app AppStatusStruct

	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.Atoi(id)
	
	if err != nil {
		errorWithNonDigit(w, err)
	} else { 

		all_ids, err := a.GetAllID()
		if err != nil {
			log.Printf("Cannot get all IDs from DB: %v", err)
		}

		if Ifcontains(all_ids, i) {
			err = json.NewDecoder(r.Body).Decode(&app)
			if err != nil {
				errorWhileDecode(w, err)
			} else {

				if len(app.APP_NAME) > 0 {
					if err = a.UpdateSelectedColumn(int64(i), "app_name", app.APP_NAME); err != nil { log.Printf("Cannot update DB row: %v", err)}
				}
				
				if len(app.UPDATE_BY) > 0 {
					if err = a.UpdateSelectedColumn(int64(i), "updated_by", app.UPDATE_BY); err != nil { log.Printf("Cannot update DB row: %v", err)}
				} 
		
				if len(app.APP_VERSION) > 0 {
					if err = a.UpdateSelectedColumn(int64(i), "app_version", app.APP_VERSION); err != nil { log.Printf("Cannot update DB row: %v", err)}
				}
				
				if len(app.ENVIRONMENT) > 0 {
					if err = a.UpdateSelectedColumn(int64(i), "env", app.ENVIRONMENT); err != nil { log.Printf("Cannot update DB row: %v", err)}
				} 
		
				if len(app.BRANCH) > 0 {
					if err = a.UpdateSelectedColumn(int64(i), "branch", app.BRANCH); err != nil { log.Printf("Cannot update DB row: %v", err)}
				} 
				
				var app_after_update interface{};
				tmp, err := a.SelectFromDBWhereID(int64(i))
				if err != nil {
					log.Printf("Cannot get row from DB: %v", err)
				}
				app_after_update = GetAppStatusStructFromStatusStruct(&tmp)
				json.NewEncoder(w).Encode(app_after_update)

			}

		} else {
			h := &HelloStruct{ID: 1, SAY: "No app with given ID"}
			json.NewEncoder(w).Encode(h)
		}		

	}
	
}


// Delete data (DELETE)
func (a *App) DeleteData(w http.ResponseWriter, r *http.Request) {
	
	var app AppStatusStruct

	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.Atoi(id)

	if err != nil {
		errorWithNonDigit(w, err)
	} else {

		all_ids, err := a.GetAllID()
		if err != nil {
			log.Printf("Cannot get all IDs from DB: %v", err)
		}

		if Ifcontains(all_ids, i) {

			err = json.NewDecoder(r.Body).Decode(&app)
			if err != nil {
				errorWhileDecode(w, err)
			} else {
				err = a.DeleteRowByID(int64(i))
				if err != nil {
					log.Printf("Cannot delete row from DB: %v", err)
				}
				h := &HelloStruct{ID: i, SAY: "Record was deleted successfully !"}
				json.NewEncoder(w).Encode(h)
			}
		} else {
			h := &HelloStruct{ID: i, SAY: "No app with given ID"}
			json.NewEncoder(w).Encode(h)
		}
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

	all_ids, err := a.GetAllID()
	if err != nil {
		log.Printf("Cannot get all IDs from DB: %v", err)
	}

	for _, k := range all_ids {
		tmp, err := a.SelectFromDBWhereID(int64(k))
		if err != nil {
			log.Printf("Cannot get row from DB: %v", err)
		}
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
		ENVIRONMENT: s.ENVIRONMENT,
		BRANCH: s.BRANCH,
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

func errorWhileDecode(w http.ResponseWriter, err error) {
	http.Error(w, "", 400)
	log.Printf("error with decode payload: %v", err)
	h := &HelloStruct{SAY: "Wrong payload syntax"}
	json.NewEncoder(w).Encode(h)
}

func errorWithNonDigit(w http.ResponseWriter, err error) {
	http.Error(w, "", 400)
	log.Printf("User insert non-digit in request: %v", err)
	h := &HelloStruct{SAY: "In /api/app/ endpoint ID should be digit"}
	json.NewEncoder(w).Encode(h)
}
