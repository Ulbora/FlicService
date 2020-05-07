package mysqldb

import (
	"fmt"
	"strconv"

	lg "github.com/Ulbora/Level_Logger"
	dbi "github.com/Ulbora/dbinterface"
)

/*
 Copyright (C) 2020 Ulbora Labs LLC. (www.ulboralabs.com)
 All rights reserved.
 Copyright (C) 2020 Ken Williamson
 All rights reserved.
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.
 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

const (
	userTest   = "select count(*) from user "
	userSearch = "SELECT u.id, u.user, u.domain, u.key, ut.user_type " +
		"FROM user u " +
		"inner join user_type ut " +
		"on u.user_type_id = ut.id " +
		"WHERE u.key = ? "
)

//User User
type User struct {
	ID       int64
	User     string
	Domain   string
	Key      string
	UserType string
}

//FlicDB FlicDB
type FlicDB interface {
	GetUser(key string) *User
}

//UserDB UserDB
type UserDB struct {
	DB  dbi.Database
	Log *lg.Logger
}

//GetNew GetNew
func (d *UserDB) GetNew() FlicDB {
	return d
}

//GetUser GetUser
func (d *UserDB) GetUser(key string) *User {
	d.Log.Debug("in get user")
	var rtn *User
	if !d.testConnection() {
		d.DB.Connect()
	}
	var a []interface{}
	a = append(a, key)
	row := d.DB.Get(userSearch, a...)
	if row != nil && len(row.Row) != 0 {
		foundRow := row.Row
		rtn = parseClientRow(&foundRow)
	}

	return rtn
}

func (d *UserDB) testConnection() bool {
	d.Log.Debug("in testConnection")
	var rtn = false
	var a []interface{}
	d.Log.Debug("d.DB: ", fmt.Sprintln(d.DB))
	rowPtr := d.DB.Test(userTest, a...)
	d.Log.Debug("rowPtr", *rowPtr)
	d.Log.Debug("after testConnection test", *rowPtr)
	if len(rowPtr.Row) != 0 {
		foundRow := rowPtr.Row
		int64Val, err := strconv.ParseInt(foundRow[0], 10, 0)
		d.Log.Debug("int64Val", int64Val)
		if err != nil {
			d.Log.Error(err)
		}
		if int64Val >= 0 {
			rtn = true
		}
	}
	return rtn
}

func parseClientRow(foundRow *[]string) *User {
	var rtn User
	if len(*foundRow) > 0 {
		int64Val, err := strconv.ParseInt((*foundRow)[0], 10, 64)
		if err == nil {
			rtn.ID = int64Val
			rtn.User = (*foundRow)[1]
			rtn.Key = (*foundRow)[2]
			rtn.Domain = (*foundRow)[3]
			rtn.UserType = (*foundRow)[4]
		}
	}
	return &rtn
}
