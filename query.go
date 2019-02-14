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

func (a *App) UpdateSelectedColumn(id int64, col, new_val string) error {

	var status StatusStruct

	err := a.DB.First(&status, id).Error
	if err != nil {
		return err
	}

	if col == "app_name" {
		if err := a.DB.Model(&status).Update("APP_NAME", new_val).Error; err != nil {return err}
	} else if col == "updated_by" {
		if err := a.DB.Model(&status).Update("UPDATE_BY", new_val).Error; err != nil {return err}
	} else if col == "app_version" {
		if err := a.DB.Model(&status).Update("APP_VERSION", new_val).Error; err != nil {return err}
	} else if col == "env" {
		if err := a.DB.Model(&status).Update("ENVIRONMENT", new_val).Error; err != nil {return err}
	} else if col == "branch" {
		if err := a.DB.Model(&status).Update("BRANCH", new_val).Error; err != nil {return err}
	}

	return nil
}

func (a *App) SelectFromDBWhereID(id int64) (StatusStruct, error) {
	
	var status StatusStruct
	
	err := a.DB.First(&status, id).Error 
	
	if err != nil {
		return status, err
	}

	return status, nil
	
}

func (a *App) GetAllID() ([]int, error) {

	var statuses []StatusStruct
	var ID []int

	err := a.DB.Find(&statuses).Error
	if err != nil {
		return ID, err
	} else {
		for _, data := range statuses {
			ID = append(ID, int(data.Model.ID))
		}
	}
	

	return ID, nil

}

func (a *App) InsertToDB(app, version, updater, env, branch string) error {

	err := a.DB.Create(&StatusStruct{
		APP_NAME: app, 
		APP_VERSION: version, 
		ENVIRONMENT: env,
		BRANCH: branch,
		UPDATE_DATE: time.Now(),
		UPDATE_BY: updater}).Error
	if err != nil {
		return err
	}

	return nil

}

func (a *App) MakeMigration() error {

	err := a.DB.AutoMigrate(&StatusStruct{}).Error
	if err != nil {
		return err
	}

	return nil
}

