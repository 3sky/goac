package main

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestMakeMIgration(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	a.MakeMigration()

	assert.Equal(t, err, nil)
}

func TestInsertToDB(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	var status StatusStruct

	a.InsertToDB("Test_run_app", "1", "UnitTest", "dev", "testing")
	a.DB.Where("app_name = ?", "Test_run_app").First(&status)
	assert.Equal(t, "Test_run_app", string(status.AppName))
	assert.Equal(t, "1", string(status.AppVersion))
	assert.Equal(t, "UnitTest", string(status.UpdateBy))
	assert.Equal(t, "dev", string(status.Environment))
	assert.Equal(t, "testing", string(status.Branch))

}

func TestSelectFromDBWhereID(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	var testStatus StatusStruct
	var status StatusStruct

	a.DB.Where("app_name = ?", "Test_run_app").First(&testStatus)

	status = a.SelectFromDBWhereID(int64(testStatus.Model.ID))

	assert.Equal(t, testStatus.AppName, status.AppName)
	assert.Equal(t, testStatus.AppVersion, status.AppVersion)
	assert.Equal(t, testStatus.UpdateBy, status.UpdateBy)

}

func TestGetAllID(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	a.InsertToDB("Test 1", "1", "Admin 1", "dev1", "testing1")
	a.InsertToDB("Test 2", "2", "Admin 2", "dev2", "testing2")
	a.InsertToDB("Test 3", "3", "Admin 3", "dev2", "testing3")
	IDs := a.GetAllID()
	assert.Len(t, IDs, 8)

}

func TestUpdateSelectedColumn(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	var s1, s2, s3 StatusStruct
	a.UpdateSelectedColumn(6, "AppName", "Test_pass")
	a.UpdateSelectedColumn(7, "UpdateBy", "Greate Tester")
	a.UpdateSelectedColumn(7, "Env", "stage")
	a.UpdateSelectedColumn(8, "AppVersion", "10")
	a.UpdateSelectedColumn(8, "Branch", "hotfix")
	a.UpdateSelectedColumn(9, "dupa", "dyap")

	s1 = a.SelectFromDBWhereID(int64(6))
	s2 = a.SelectFromDBWhereID(int64(7))
	s3 = a.SelectFromDBWhereID(int64(8))

	assert.Equal(t, "Test_pass", s1.AppName)
	assert.Equal(t, "1", s1.AppVersion)
	assert.Equal(t, "Admin 1", s1.UpdateBy)

	assert.Equal(t, "Test 2", s2.AppName)
	assert.Equal(t, "2", s2.AppVersion)
	assert.Equal(t, "Greate Tester", s2.UpdateBy)
	assert.Equal(t, "stage", s2.Environment)

	assert.Equal(t, "Test 3", s3.AppName)
	assert.Equal(t, "10", s3.AppVersion)
	assert.Equal(t, "Admin 3", s3.UpdateBy)
	assert.Equal(t, "hotfix", s3.Branch)

}

func TestDeleteRowByID(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	IDs := a.GetAllID()

	for _, k := range IDs {
		a.DeleteRowByID(int64(k))
	}

	te := a.GetAllID()

	assert.Len(t, te, 0)

}
