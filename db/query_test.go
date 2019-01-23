package db

import (
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)


func TestMakeMIgration(t *testing.T) {

	MakeMIgration()

	db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	CheckErr(err)
	defer db.Close()
	
	assert.Equal(t, err, nil)
}

func TestInsertToDB(t *testing.T) {
	
	
	var status StatusStruct

	db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	CheckErr(err)
	defer db.Close()

	InsertToDB("Test_run_app", "1", "UnitTest")
	db.Where("app_name = ?", "Test_run_app").First(&status)
	assert.Equal(t, "Test_run_app", string(status.APP_NAME))
	assert.Equal(t, "1",            string(status.APP_VERSION))
	assert.Equal(t, "UnitTest",     string(status.UPDATE_BY))

}


func TestSelectFromDBWhereID(t *testing.T) {

	var test_status StatusStruct
	var status StatusStruct

	db, err := gorm.Open("sqlite3", "./SimpleDB.db")
	CheckErr(err)

	db.Where("app_name = ?", "Test_run_app").First(&test_status)

	status = SelectFromDBWhereID(int64(test_status.Model.ID))

	assert.Equal(t, test_status.APP_NAME,    status.APP_NAME)
	assert.Equal(t, test_status.APP_VERSION, status.APP_VERSION)
	assert.Equal(t, test_status.UPDATE_BY,   status.UPDATE_BY)
	
}


func TestGetAllID(t *testing.T) {

	InsertToDB("Test 1", "1", "Admin 1")
	InsertToDB("Test 2", "2", "Admin 2")
	InsertToDB("Test 3", "3", "Admin 3")
	IDs := GetAllID()
	assert.Len(t, IDs, 4)
	
}

func TestUpdateSelectedColumn(t *testing.T) {

	var s1, s2, s3 StatusStruct
	UpdateSelectedColumn(2,"app_name", "Test_pass")
	UpdateSelectedColumn(3,"updated_by", "Greate Tester")
	UpdateSelectedColumn(4,"app_version", "10")
	UpdateSelectedColumn(5,"dupa", "dyap")
	
	s1 = SelectFromDBWhereID(int64(2))
	s2 = SelectFromDBWhereID(int64(3))
	s3 = SelectFromDBWhereID(int64(4))
	

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

	IDs := GetAllID()

	for _, k := range IDs {
		DeleteRowByID(int64(k))
	}

	te := GetAllID()

	assert.Len(t, te, 0)

	err := os.Remove("./SimpleDB.db")
	CheckErr(err)
}
