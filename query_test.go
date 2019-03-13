package main

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestMakeMigration(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	err := a.MakeMigration()
	if err != nil {
		fmt.Printf("Cannot make schema migration: %v", err)
	}

	db, err := gorm.Open("sqlite3", "TestDB.db")
	if err != nil {
		fmt.Printf("Cannot open TestDB.db: %v", err)
	}
	defer db.Close()

	assert.Equal(t, err, nil)
}

func TestInsertToDB(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	var status StatusStruct

	err := a.InsertToDB("Test_run_app", "1", "UnitTest", "dev", "testing")
	if err != nil {
		fmt.Printf("Cannot insert data in TestInsertToDB: %v", err)
	}

	err = a.DB.Where("app_name = ?", "Test_run_app").First(&status).Error
	if err != nil {
		fmt.Printf("Cannot find data in TestInsertToDB.db: %v", err)
	}
	assert.Equal(t, "Test_run_app", string(status.AppName))
	assert.Equal(t, "1", string(status.AppVersion))
	assert.Equal(t, "UnitTest", string(status.UpdateBy))
	assert.Equal(t, "dev", string(status.Environment))
	assert.Equal(t, "testing", string(status.Branch))

}

func TestSelectFromDBWhereID(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	var testStatus StatusStruct
	var status StatusStruct

	err := a.DB.Where("app_name = ?", "Test_run_app").First(&testStatus).Error
	if err != nil {
		fmt.Printf("Cannot find data in TestDB.db: %v", err)
	}
	status, err = a.SelectFromDBWhereID(int64(testStatus.Model.ID))
	if err != nil {
		fmt.Printf("Cannot get row from DB in TestSelectFromDBWhereID: %v", err)
	}

	assert.Equal(t, testStatus.AppName, status.AppName)
	assert.Equal(t, testStatus.AppVersion, status.AppVersion)
	assert.Equal(t, testStatus.UpdateBy, status.UpdateBy)

}

func TestGetAllID(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	err := a.InsertToDB("Test 1", "1", "Admin 1", "dev1", "testing1")
	if err != nil {
		fmt.Printf("Cannot insert data in TestGetAllID: %v", err)
	}

	err = a.InsertToDB("Test 2", "2", "Admin 2", "dev2", "testing2")
	if err != nil {
		fmt.Printf("Cannot insert data in TestGetAllID: %v", err)
	}

	err = a.InsertToDB("Test 3", "3", "Admin 3", "dev2", "testing3")
	if err != nil {
		fmt.Printf("Cannot insert data in TestGetAllID: %v", err)
	}

	IDs, err := a.GetAllID()
	if err != nil {
		fmt.Printf("Cannot get All IDs from TestGetAllID: %v", err)
	}

	assert.Len(t, IDs, 8)

}

func TestUpdateSelectedColumn(t *testing.T) {

	a := createTestDBConnection()
	defer a.DB.Close()

	var s1, s2, s3 StatusStruct
	a.UpdateSelectedColumn(6, "app_name", "Test_pass")
	a.UpdateSelectedColumn(7, "update_by", "Greate Tester")
	a.UpdateSelectedColumn(7, "environment", "stage")
	a.UpdateSelectedColumn(8, "app_version", "10")
	a.UpdateSelectedColumn(8, "branch", "hotfix")
	a.UpdateSelectedColumn(9, "dupa", "dyap")

	s1, err := a.SelectFromDBWhereID(int64(6))
	if err != nil {
		fmt.Printf("Cannot get row from DB in TestDeleteRowByID: %v", err)
	}

	s2, err = a.SelectFromDBWhereID(int64(7))
	if err != nil {
		fmt.Printf("Cannot get row from DB in TestDeleteRowByID: %v", err)
	}

	s3, err = a.SelectFromDBWhereID(int64(8))
	if err != nil {
		fmt.Printf("Cannot get row from DB in TestDeleteRowByID: %v", err)
	}

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

	a := createTestDBConnection()
	defer a.DB.Close()

	IDs, err := a.GetAllID()
	if err != nil {
		fmt.Printf("Cannot get All IDs from TestDeleteRowByID: %v", err)
	}

	for _, k := range IDs {
		err := a.DeleteRowByID(int64(k))
		if err != nil {
			fmt.Printf("Error with deleteing row by ID: %v", err)
		}
	}

	te, err := a.GetAllID()
	if err != nil {
		fmt.Printf("Cannot get Alls ID from TestDeleteRowByID: %v", err)
	}

	assert.Len(t, te, 0)

}

func createTestDBConnection() *App {

	a := &App{}
	var err error

	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	if err != nil {
		fmt.Printf("Cannot open TestDB.db: %v", err)
	}

	return a

}
