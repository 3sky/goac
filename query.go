package main

import (

	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

)

type StatusStruct struct {
	gorm.Model
	APP_NAME string `main"app_name"` 
	APP_VERSION string `main:"app_version"`
	ENVIRONMENT string `main:"env"`
	BRANCH string `main:"branch"`
	UPDATE_DATE time.Time `main:"updated_date"`
	UPDATE_BY string `main:"updated_by"`
}


func (a *App) DeleteRowByID(id int64){

	var status StatusStruct

	a.DB.First(&status, id)
	a.DB.Delete(&status)
}

func (a *App) UpdateSelectedColumn(id int64, col, new_val string) {

	var status StatusStruct

	a.DB.First(&status, id)

	if col == "app_name" {
		a.DB.Model(&status).Update("APP_NAME", new_val)
	} else if col == "updated_by" {
		a.DB.Model(&status).Update("UPDATE_BY", new_val)
	} else if col == "app_version" {
		a.DB.Model(&status).Update("APP_VERSION", new_val)
	} else if col == "env" {
		a.DB.Model(&status).Update("ENVIRONMENT", new_val)
	} else if col == "branch" {
		a.DB.Model(&status).Update("BRANCH", new_val)
	}

}

func (a *App) SelectFromDBWhereID(id int64) StatusStruct {
	
	var status StatusStruct
	
	a.DB.First(&status, id)

	return status
	
}

func (a *App) GetAllID() []int {

	var statuses []StatusStruct
	var ID []int

	a.DB.Find(&statuses)

	for _, data := range statuses {
		ID = append(ID, int(data.Model.ID))
	}

	return ID

}

func (a *App) InsertToDB(app, version, updater, env, branch string) {

	a.DB.Create(&StatusStruct{
		APP_NAME: app, 
		APP_VERSION: version, 
		ENVIRONMENT: env,
		BRANCH: branch,
		UPDATE_DATE: time.Now(),
		UPDATE_BY: updater})

}

func (a *App) MakeMigration() {

	a.DB.AutoMigrate(&StatusStruct{})
}

