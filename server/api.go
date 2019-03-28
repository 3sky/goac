package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//HelloStruct for valid JSON
type HelloStruct struct {
	ID  int    `json:"ID,omitempty"`
	Say string `json:"Info,omitempty"`
}

// AppStatusStruct - Similar struct to db.StatusStruct, but without gorm.Model
type AppStatusStruct struct {
	ID          int    `json:"id"`
	AppName     string `json:"app_name"`
	AppVersion  string `json:"app_version"`
	Environment string `json:"environment"`
	Branch      string `json:"branch"`
	IP          string `json:"ip"`
	UpdateDate  string `json:"update_date"`
	UpdateBy    string `json:"update_by"`
}

// AllApp - list of some app struct
type AllApp struct {
	Name string
	App  []AppStatusStruct
}

//SayHello seyhello
func SayHello(w http.ResponseWriter, r *http.Request) {
	h := &HelloStruct{ID: 1, Say: "Hello Kuba"}
	json.NewEncoder(w).Encode(h)
}

// DisplayAppByID - Display App status By ID (GET)
func (a *App) DisplayAppByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.Atoi(id)
	if err != nil {
		errorWithNonDigit(w, err)
	} else {

		allIds, err := a.GetAllID()
		if err != nil {
			log.Printf("Cannot get all IDs from DB: %v", err)
		}

		if Ifcontains(allIds, int64(i)) {
			//var app interface{}
			tmp, err := a.SelectFromDBWhereID(int64(i))
			if err != nil {
				log.Printf("Cannot get row from DB: %v", err)
			}
			//app = GetAppStatusStructFromStatusStruct(&tmp)
			json.NewEncoder(w).Encode(tmp)
		} else {
			h := &HelloStruct{ID: i, Say: "No app with given ID"}
			json.NewEncoder(w).Encode(h)
		}
	}

}

// AddNewApp - Add New App (POST)
func (a *App) AddNewApp(w http.ResponseWriter, r *http.Request) {

	var app AppStatusStruct
	var updater string

	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		errorWhileDecode(w, err)
	} else {

		if a.validateInsert(app.AppName, app.Environment) {

			if len(app.UpdateBy) == 0 {
				updater = "random guy"
			} else {
				updater = app.UpdateBy
			}

			err = a.InsertToDB(app.AppName, app.AppVersion, updater, app.Environment, app.Branch, app.IP)
			if err != nil {
				log.Printf("Cannot insert row into DB: %v", err)
			}

		} else {
			h := &HelloStruct{Say: "This app already exits on this environment! "}
			json.NewEncoder(w).Encode(h)
		}
	}
}

//SearchApp - search app with name
func (a *App) SearchApp(w http.ResponseWriter, r *http.Request) {

	var app *AppStatusStruct

	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		errorWhileDecode(w, err)
	} else {

		if (len(app.AppName) != 0) && (len(app.Environment) != 0) {

			tmp, err := a.SearchInDB(app.AppName, app.Environment)

			switch {
			case err == nil:
				log.Printf("Error while searching app: %v", err)
			case err.Error() == "record not found":
			default:
				log.Printf("Error while searching app: %v", err)
			}
			if len(tmp.AppName) == 0 {
				h := &HelloStruct{Say: "No such app ! "}
				json.NewEncoder(w).Encode(h)
			} else {
				json.NewEncoder(w).Encode(tmp)
			}

		} else {
			h := &HelloStruct{Say: "Application name and environment are mandatory ! "}
			json.NewEncoder(w).Encode(h)
		}
	}
}

// UpdateData Update Data (PUT)
func (a *App) UpdateData(w http.ResponseWriter, r *http.Request) {

	var app AppStatusStruct

	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		errorWithNonDigit(w, err)
	} else {

		allIds, err := a.GetAllID()
		if err != nil {
			log.Printf("Cannot get all IDs from DB: %v", err)
		}

		if Ifcontains(allIds, i) {
			err = json.NewDecoder(r.Body).Decode(&app)
			if err != nil {
				errorWhileDecode(w, err)
			} else {

				if len(app.AppName) > 0 {
					if err = a.UpdateSelectedColumn(int64(i), "app_name", app.AppName); err != nil {
						log.Printf("Cannot update DB row: %v", err)
					}
				}

				if len(app.UpdateBy) > 0 {
					if err = a.UpdateSelectedColumn(int64(i), "update_by", app.UpdateBy); err != nil {
						log.Printf("Cannot update DB row: %v", err)
					}
				}

				if len(app.AppVersion) > 0 {
					if err = a.UpdateSelectedColumn(int64(i), "app_version", app.AppVersion); err != nil {
						log.Printf("Cannot update DB row: %v", err)
					}
				}

				if len(app.Environment) > 0 {
					if err = a.UpdateSelectedColumn(int64(i), "environment", app.Environment); err != nil {
						log.Printf("Cannot update DB row: %v", err)
					}
				}

				if len(app.IP) > 0 {
					if err = a.UpdateSelectedColumn(int64(i), "ip", app.IP); err != nil {
						log.Printf("Cannot update DB row: %v", err)
					}
				}
				if len(app.Branch) > 0 {
					if err = a.UpdateSelectedColumn(int64(i), "branch", app.Branch); err != nil {
						log.Printf("Cannot update DB row: %v", err)
					}
				}
				err := a.UpdateCurrentDate(int64(i))
				if err != nil {
					log.Printf("Cannot update insert time : %v", err)
				}
				//var appAfterUpdate interface{}
				tmp, err := a.SelectFromDBWhereID(int64(i))
				if err != nil {
					log.Printf("Cannot get row from DB: %v", err)
				}
				//appAfterUpdate = GetAppStatusStructFromStatusStruct(&tmp)
				json.NewEncoder(w).Encode(tmp)

			}

		} else {
			h := &HelloStruct{ID: 1, Say: "No app with given ID"}
			json.NewEncoder(w).Encode(h)
		}

	}

}

// DeleteData - Delete app (DELETE)
func (a *App) DeleteData(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.Atoi(id)

	if err != nil {
		errorWithNonDigit(w, err)
	} else {

		allIds, err := a.GetAllID()
		if err != nil {
			log.Printf("Cannot get all IDs from DB: %v", err)
		}

		if Ifcontains(allIds, int64(i)) {
			err = a.DeleteRowByID(int64(i))
			if err != nil {
				log.Printf("Cannot delete row from DB: %v", err)
			}
			h := &HelloStruct{ID: i, Say: "Record was deleted successfully !"}
			json.NewEncoder(w).Encode(h)
		} else {
			h := &HelloStruct{ID: i, Say: "No app with given ID"}
			json.NewEncoder(w).Encode(h)
		}
	}

}

// DisplayAllApp - Display all app (GET)
func (a *App) DisplayAllApp(w http.ResponseWriter, r *http.Request) {

	Apps := []AppStatusStruct{}

	data := AllApp{
		Name: "Apps",
		App:  Apps,
	}

	allIds, err := a.GetAllID()
	if err != nil {
		log.Printf("Cannot get all IDs from DB: %v", err)
	}

	for _, k := range allIds {
		tmp, err := a.SelectFromDBWhereID(int64(k))
		if err != nil {
			log.Printf("Cannot get row from DB: %v", err)
		}
		data.AddApp(&tmp)
	}
	json.NewEncoder(w).Encode(data)
}

//AddApp - add apps to one big struct
func (aa *AllApp) AddApp(app *AppStatusStruct) []AppStatusStruct {
	valueOffApp := *app
	aa.App = append(aa.App, valueOffApp)
	return aa.App
}

// GetAppStatusStructFromStatusStruct - Convert StatusStruct to AppStruct
/**
func GetAppStatusStructFromStatusStruct(s *AppStatusStruct) *AppStatusStruct {
	s2 := &AppStatusStruct{
		ID:          s.ID,
		AppName:     s.AppName,
		AppVersion:  s.AppVersion,
		Environment: s.Environment,
		Branch:      s.Branch,
		IP:          s.IP,
		//UpdateDate:  s.UpdateDate.Format("2006-01-02 15:04:05"),
		UpdateDate: s.UpdateDate,
		UpdateBy:   s.UpdateBy,
	}
	return s2
}
**/

// Ifcontains - check if list contain ID
func Ifcontains(s []int, e int64) bool {
	for _, a := range s {
		if int64(a) == e {
			return true
		}
	}
	return false
}

func errorWhileDecode(w http.ResponseWriter, err error) {
	http.Error(w, "", 400)
	log.Printf("error with decode payload: %v", err)
	h := &HelloStruct{Say: "Wrong payload syntax"}
	json.NewEncoder(w).Encode(h)
}

func errorWithNonDigit(w http.ResponseWriter, err error) {
	http.Error(w, "", 400)
	log.Printf("User insert non-digit in request: %v", err)
	h := &HelloStruct{Say: "In /api/app/ endpoint ID should be digit"}
	json.NewEncoder(w).Encode(h)
}

func (a *App) validateInsert(name, env string) bool {

	if (len(name) != 0) && (len(env) != 0) {

		_, err := a.SearchInDB(name, env)
		if err == nil {
			return false
		} else if err.Error() == "record not found" {
			return true
		} else {
			return false
		}
	}

	return false
}
