package db

import (

	"time"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)


func DeleteRowByID(id int64){

	var status StatusStruct

	db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	CheckErr(err)
	defer db.Close()

	db.First(&status, id)
	db.Delete(&status)
}

func UpdateSelectedColumn(id int64, col, new_val string) {

	var status StatusStruct

	db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	CheckErr(err)
	defer db.Close()

	db.First(&status, id)

	if col == "app_name" {
		db.Model(&status).Update("APP_NAME", new_val)
	} else if col == "updated_by" {
		db.Model(&status).Update("UPDATE_BY", new_val)
	} else if col == "app_version" {
		db.Model(&status).Update("APP_VERSION", new_val)
	}

}

func SelectFromDBWhereID(id int64) StatusStruct {
	
	var status StatusStruct
	
	db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	CheckErr(err)
	defer db.Close()
	
	db.First(&status, id)


	return status
	
}

func GetAllID() []int {

	var statuses []StatusStruct
	var ID []int

	db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	CheckErr(err)
	defer db.Close()

	db.Find(&statuses)

	for _, data := range statuses {
		ID = append(ID, int(data.Model.ID))
	}

	return ID

}

func InsertToDB(app, version, updater string) {

	db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	CheckErr(err)
	defer db.Close()

	db.Create(&StatusStruct{APP_NAME: app, APP_VERSION: version, UPDATE_DATE: time.Now(), UPDATE_BY: updater})

}

func MakeMIgration() {

	db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	CheckErr(err)
	defer db.Close()

	db.AutoMigrate(&StatusStruct{})
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}