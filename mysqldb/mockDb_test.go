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
