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
	UPDATE_DATE time.Time `main:"updated_date"`
	UPDATE_BY string `main:"updated_by"`
}


func (a *App) DeleteRowByID(id int64){

	var status StatusStruct

	db := a.DB

	db.First(&status, id)
	db.Delete(&status)
}

func (a *App) UpdateSelectedColumn(id int64, col, new_val string) {

	var status StatusStruct

	// db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	// CheckErr(err)
	// defer db.Close()

	db := a.DB

	db.First(&status, id)

	if col == "app_name" {
		db.Model(&status).Update("APP_NAME", new_val)
	} else if col == "updated_by" {
		db.Model(&status).Update("UPDATE_BY", new_val)
	} else if col == "app_version" {
		db.Model(&status).Update("APP_VERSION", new_val)
	}

}

func (a *App) SelectFromDBWhereID(id int64) StatusStruct {
	
	var status StatusStruct
	
	// db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	// CheckErr(err)
	// defer db.Close()
	db := a.DB

	db.First(&status, id)


	return status
	
}

func (a *App) GetAllID() []int {

	var statuses []StatusStruct
	var ID []int

	db := a.DB
	// db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	// CheckErr(err)
	// defer db.Close()

	db.Find(&statuses)

	for _, data := range statuses {
		ID = append(ID, int(data.Model.ID))
	}

	return ID

}

func (a *App) InsertToDB(app, version, updater string) {

	// db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	// CheckErr(err)
	// defer db.Close()

	db := a.DB

	db.Create(&StatusStruct{APP_NAME: app, APP_VERSION: version, UPDATE_DATE: time.Now(), UPDATE_BY: updater})

}

func (a *App) MakeMigration() {

	// db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	// CheckErr(err)
	// defer db.Close()

	// db := 

	a.DB.AutoMigrate(&StatusStruct{})
}

