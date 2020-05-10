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
	getRow.Row = []string{"1", "some user", "somedomain.com", "456211111", "customer", "1"}
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
	getRow.Row = []string{"1", "some user", "somedomain.com", "456211111", "customer", "1"}
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
	getRow.Row = []string{"1", "some user", "somedomain.com", "456211111", "customer", "1"}
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

func TestUserDB_GetBqTable(t *testing.T) {

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

	var rows [][]string
	row1 := []string{"1", "some_bq_table"}
	rows = append(rows, row1)
	var dbrows db.DbRows
	dbrows.Rows = rows
	mydb.MockRows1 = &dbrows

	dbii = &mydb

	udbii.DB = dbii

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udbii.Log = &l

	dbii.Connect()

	udbi := udbii.GetNew()
	bqt := udbi.GetFlicTable()
	fmt.Println("bqt", bqt)
	if bqt.ID == 0 {
		t.Fail()
	}
}

func TestUserDB_GetBqTableConnect(t *testing.T) {

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

	var rows [][]string
	row1 := []string{"1", "some_bq_table"}
	rows = append(rows, row1)
	var dbrows db.DbRows
	dbrows.Rows = rows
	mydb.MockRows1 = &dbrows

	dbii = &mydb

	udbii.DB = dbii

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udbii.Log = &l

	dbii.Connect()

	udbi := udbii.GetNew()
	bqt := udbi.GetFlicTable()
	fmt.Println("bqt", bqt)
	if bqt.ID == 0 {
		t.Fail()
	}
}

func TestUserDB_SetBqTable(t *testing.T) {

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

	mydb.MockUpdateSuccess1 = true

	dbii = &mydb

	udbii.DB = dbii

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udbii.Log = &l

	dbii.Connect()

	udbi := udbii.GetNew()
	var ft FlicTable
	ft.ID = 1
	ft.Name = "test_table_1"
	suc := udbi.SetFlicTable(&ft)
	fmt.Println("bqt suc", suc)
	if !suc {
		t.Fail()
	}
}

func TestUserDB_SetBqTableConnect(t *testing.T) {

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

	mydb.MockUpdateSuccess1 = true

	dbii = &mydb

	udbii.DB = dbii

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udbii.Log = &l

	dbii.Connect()

	udbi := udbii.GetNew()
	var ft FlicTable
	ft.ID = 1
	ft.Name = "test_table_1"
	suc := udbi.SetFlicTable(&ft)
	fmt.Println("bqt suc", suc)
	if !suc {
		t.Fail()
	}
}
