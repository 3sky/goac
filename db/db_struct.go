package db

import (
	"time"
	"github.com/jinzhu/gorm"
)


type StatusStruct struct {
	gorm.Model
	APP_NAME string `db:"app_name"` 
	APP_VERSION string `db:"app_version"`
	UPDATE_DATE time.Time `db:"updated_date"`
	UPDATE_BY string `db:"updated_by"`
}
