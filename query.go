package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//StatusStruct struck directly from DB
type StatusStruct struct {
	gorm.Model
	AppName     string    `main:"app_name"`
	AppVersion  string    `main:"app_version"`
	Environment string    `main:"environment"`
	Branch      string    `main:"branch"`
	UpdateDate  time.Time `main:"update_date"`
	UpdateBy    string    `main:"update_by"`
}

//DeleteRowByID delete from DB
func (a *App) DeleteRowByID(id int64) {

	var status StatusStruct

	a.DB.First(&status, id)
	a.DB.Delete(&status)
}

//UpdateSelectedColumn - Update Selected Column
func (a *App) UpdateSelectedColumn(id int64, col, newVal string) {

	var status StatusStruct

	a.DB.First(&status, id)

	if col == "AppName" {
		a.DB.Model(&status).Update("app_name", newVal)
	} else if col == "UpdateBy" {
		a.DB.Model(&status).Update("update_by", newVal)
	} else if col == "AppVersion" {
		a.DB.Model(&status).Update("app_version", newVal)
	} else if col == "Env" {
		a.DB.Model(&status).Update("environment", newVal)
	} else if col == "Branch" {
		a.DB.Model(&status).Update("branch", newVal)
	}

}

//SelectFromDBWhereID - Select From DB Where ID
func (a *App) SelectFromDBWhereID(id int64) StatusStruct {

	var status StatusStruct

	a.DB.First(&status, id)

	return status

}

//GetAllID - Get All ID
func (a *App) GetAllID() []int {

	var statuses []StatusStruct
	var ID []int

	a.DB.Find(&statuses)

	for _, data := range statuses {
		ID = append(ID, int(data.Model.ID))
	}

	return ID

}

//InsertToDB - Insert To DB
func (a *App) InsertToDB(app, version, updater, env, branch string) {

	a.DB.Create(&StatusStruct{
		AppName:     app,
		AppVersion:  version,
		Environment: env,
		Branch:      branch,
		UpdateDate:  time.Now(),
		UpdateBy:    updater})

}

//MakeMigration - Make Migration
func (a *App) MakeMigration() {

	a.DB.AutoMigrate(&StatusStruct{})
}
