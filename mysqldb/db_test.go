package mysqldb

import (
	"fmt"
	"testing"

	lg "github.com/Ulbora/Level_Logger"
	db "github.com/Ulbora/dbinterface"
	mdb "github.com/Ulbora/dbinterface_mysql"
)

func TestUserDB_GetUser(t *testing.T) {

	var dbii db.Database
	var udbii UserDB
	var mydb mdb.MyDBMock
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

	var mTestRow db.DbRow
	mTestRow.Row = []string{"1"}
	mydb.MockTestRow = &mTestRow

	var getRow db.DbRow
	getRow.Row = []string{"1", "some user", "somedomain.com", "456211111", "customer"}
	mydb.MockRow1 = &getRow

	dbii = &mydb

	udbii.DB = dbii

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udbii.Log = &l

	dbii.Connect()

	udbi := udbii.GetNew()
	u := udbi.GetUser("61616dfggdf5g64gf4")
	fmt.Println("u", u)
	if u.ID == 0 {
		t.Fail()
	}
}

func TestUserDB_GetUserErr(t *testing.T) {

	var dbii db.Database
	var udbii UserDB
	var mydb mdb.MyDBMock
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

	var mTestRow db.DbRow
	mTestRow.Row = []string{}
	mydb.MockTestRow = &mTestRow

	var getRow db.DbRow
	getRow.Row = []string{"1", "some user", "somedomain.com", "456211111", "customer"}
	mydb.MockRow1 = &getRow

	dbii = &mydb

	udbii.DB = dbii

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udbii.Log = &l

	dbii.Connect()

	udbi := udbii.GetNew()
	u := udbi.GetUser("61616dfggdf5g64gf4")
	fmt.Println("u", u)
	if u.ID == 0 {
		t.Fail()
	}
}

func TestUserDB_GetUserBidId(t *testing.T) {

	var dbii db.Database
	var udbii UserDB
	var mydb mdb.MyDBMock
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

	var mTestRow db.DbRow
	mTestRow.Row = []string{"s"}
	mydb.MockTestRow = &mTestRow

	var getRow db.DbRow
	getRow.Row = []string{"1", "some user", "somedomain.com", "456211111", "customer"}
	mydb.MockRow1 = &getRow

	dbii = &mydb

	udbii.DB = dbii

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udbii.Log = &l

	dbii.Connect()

	udbi := udbii.GetNew()
	u := udbi.GetUser("61616dfggdf5g64gf4")
	fmt.Println("u", u)
	if u.ID == 0 {
		t.Fail()
	}
}
