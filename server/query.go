package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//StatusStruct - data from DB
type StatusStruct struct {
	gorm.Model
	AppName     string    `main:"app_name"`
	AppVersion  string    `main:"app_version"`
	Environment string    `main:"environment"`
	Branch      string    `main:"branch"`
	UpdateDate  time.Time `main:"update_date"`
	UpdateBy    string    `main:"update_by"`
}

//DeleteRowByID - Delete row
func (a *App) DeleteRowByID(id int64) error {

	var status StatusStruct

	err := a.DB.First(&status, id).Error
	if err != nil {
		return err
	}

	err = a.DB.Delete(&status).Error
	if err != nil {
		return err
	}

	return nil
}

//UpdateSelectedColumn - Update column
func (a *App) UpdateSelectedColumn(id int64, col, newVal string) error {

	var status StatusStruct

	err := a.DB.First(&status, id).Error
	if err != nil {
		return err
	}

	if col == "app_name" {
		if err := a.DB.Model(&status).Update("app_name", newVal).Error; err != nil {
			return err
		}
	} else if col == "update_by" {
		if err := a.DB.Model(&status).Update("update_by", newVal).Error; err != nil {
			return err
		}
	} else if col == "app_version" {
		if err := a.DB.Model(&status).Update("app_version", newVal).Error; err != nil {
			return err
		}
	} else if col == "environment" {
		if err := a.DB.Model(&status).Update("environment", newVal).Error; err != nil {
			return err
		}
	} else if col == "branch" {
		if err := a.DB.Model(&status).Update("branch", newVal).Error; err != nil {
			return err
		}
	}

	return nil
}

//UpdateCurrentDate - updating date for current while update row
func (a *App) UpdateCurrentDate(id int64) error {

	var status StatusStruct

	err := a.DB.First(&status, id).Error
	if err != nil {
		return err
	}
	if err := a.DB.Model(&status).Update("update_date", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

//SelectFromDBWhereID - Select row from DB
func (a *App) SelectFromDBWhereID(id int64) (StatusStruct, error) {

	var status StatusStruct

	err := a.DB.First(&status, id).Error

	if err != nil {
		return status, err
	}

	return status, nil

}

//GetAllID - Get all IDs frm DB
func (a *App) GetAllID() ([]int, error) {

	var statuses []StatusStruct
	var ID []int

	err := a.DB.Find(&statuses).Error
	if err != nil {
		return ID, err
	}

	for _, data := range statuses {
		ID = append(ID, int(data.Model.ID))
	}

	return ID, nil

}

//InsertToDB - Insert Data to DB
func (a *App) InsertToDB(app, version, updater, Environment, Branch string) error {

	err := a.DB.Create(&StatusStruct{
		AppName:     app,
		AppVersion:  version,
		Environment: Environment,
		Branch:      Branch,
		UpdateDate:  time.Now(),
		UpdateBy:    updater}).Error
	if err != nil {
		return err
	}

	return nil

}

//MakeMigration - Make schema migration
func (a *App) MakeMigration() error {

	err := a.DB.AutoMigrate(&StatusStruct{}).Error
	if err != nil {
		return err
	}

	return nil
}
