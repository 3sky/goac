package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)


func TestMakeMIgration(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	a.MakeMigration()

	db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	CheckErr(err)
	defer db.Close()
	
	assert.Equal(t, err, nil)
}

func TestInsertToDB(t *testing.T) {
	
	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	var status StatusStruct

	// db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	// CheckErr(err)
	// defer db.Close()

	a.InsertToDB("Test_run_app", "1", "UnitTest")
	a.DB.Where("app_name = ?", "Test_run_app").First(&status)
	assert.Equal(t, "Test_run_app", string(status.APP_NAME))
	assert.Equal(t, "1",            string(status.APP_VERSION))
	assert.Equal(t, "UnitTest",     string(status.UPDATE_BY))

}


func TestSelectFromDBWhereID(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	var test_status StatusStruct
	var status StatusStruct

	// db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	// CheckErr(err)

	a.DB.Where("app_name = ?", "Test_run_app").First(&test_status)

	status = a.SelectFromDBWhereID(int64(test_status.Model.ID))

	assert.Equal(t, test_status.APP_NAME,    status.APP_NAME)
	assert.Equal(t, test_status.APP_VERSION, status.APP_VERSION)
	assert.Equal(t, test_status.UPDATE_BY,   status.UPDATE_BY)
	
}


func TestGetAllID(t *testing.T) {

	a := &App{}
	var err error
	a.DB, err = gorm.Open("sqlite3", "TestDB.db")
	CheckErr(err)
	defer a.DB.Close()

	a.InsertToDB("Test 1", "1", "Admin 1")
	a.InsertToDB("Test 2", "2", "Admin 2")
	a.InsertToDB("Test 3", "3", "Admin 3")
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
	a.UpdateSelectedColumn(6,"app_name", "Test_pass")
	a.UpdateSelectedColumn(7,"updated_by", "Greate Tester")
	a.UpdateSelectedColumn(8,"app_version", "10")
	a.UpdateSelectedColumn(9,"dupa", "dyap")
	
	s1 = a.SelectFromDBWhereID(int64(6))
	s2 = a.SelectFromDBWhereID(int64(7))
	s3 = a.SelectFromDBWhereID(int64(8))
	

	assert.Equal(t, "Test_pass", s1.APP_NAME)
	assert.Equal(t, "1",         s1.APP_VERSION)
	assert.Equal(t, "Admin 1",  s1.UPDATE_BY)

	assert.Equal(t, "Test 2",        s2.APP_NAME)
	assert.Equal(t, "2",             s2.APP_VERSION)
	assert.Equal(t, "Greate Tester", s2.UPDATE_BY)

	assert.Equal(t, "Test 3",  s3.APP_NAME)
	assert.Equal(t, "10",      s3.APP_VERSION)
	assert.Equal(t, "Admin 3", s3.UPDATE_BY, )

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
