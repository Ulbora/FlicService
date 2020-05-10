// +build integration move to top

package mysqldb

import (
	"fmt"
	"testing"

	lg "github.com/Ulbora/Level_Logger"
	db "github.com/Ulbora/dbinterface"
	mdb "github.com/Ulbora/dbinterface_mysql"
)

func TestUserDBi_GetUser(t *testing.T) {

	var dbii db.Database
	var udbii UserDB
	var mydb mdb.MyDB
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

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

func TestUserDBi_GetBqTable(t *testing.T) {

	var dbii db.Database
	var udbii UserDB
	var mydb mdb.MyDB
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

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

func TestUserDBi_SetBqTable(t *testing.T) {

	var dbii db.Database
	var udbii UserDB
	var mydb mdb.MyDB
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

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
