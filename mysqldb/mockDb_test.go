package mysqldb

import (
	"fmt"
	"testing"
)

func TestMockDB_GetUser(t *testing.T) {
	// var dbii db.Database
	var udbii MockUserDB
	var u User
	u.ID = 1
	udbii.MockUser = u

	udbi := udbii.GetNew()
	fu := udbi.GetUser("61616dfggdf5g64gf4")
	fmt.Println("u", fu)
	if u.ID == 0 {
		t.Fail()
	}
}

func TestMockDB_GetBqTable(t *testing.T) {
	// var dbii db.Database
	var udbii MockUserDB
	var bqt FlicTable
	bqt.ID = 5
	bqt.Name = "test_table"
	udbii.MockFlicTable = bqt

	udbi := udbii.GetNew()
	fu := udbi.GetFlicTable()
	fmt.Println("fu", fu)
	if fu.ID == 0 {
		t.Fail()
	}
}

func TestMockDB_SetBqTable(t *testing.T) {
	// var dbii db.Database
	var udbii MockUserDB
	udbii.MockSetSuccess = true
	var bqt FlicTable
	bqt.ID = 5
	bqt.Name = "test_table"

	udbi := udbii.GetNew()
	suc := udbi.SetFlicTable(&bqt)
	fmt.Println("suc", suc)
	if !suc {
		t.Fail()
	}
}
