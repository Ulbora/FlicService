package manager

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

import (
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
	ph "github.com/Ulbora/AnalyticPusher"
	pu "github.com/Ulbora/BigQueryPuller"
	db "github.com/Ulbora/FlicService/mysqldb"
	lg "github.com/Ulbora/Level_Logger"
)

const (
	searchTypeZip = "zip"
	searchTypeID  = "id"

	analyticsTable = "flic_analytics"
)

//Flic Flic
type Flic struct {
	Key            string    `json:"id"`
	Lic            string    `json:"license"`
	ExpDate        time.Time `json:"expDate"`
	LicName        string    `json:"licenseName"`
	BusName        string    `json:"businessName"`
	PremiseAddress string    `json:"premiseAddress"`
	//PremiseZip     string    `json:"premiseZip"`
	MailingAddress string `json:"mailingAddress"`
	Phone          string `json:"phone"`
}

//FlicList FlicList
type FlicList struct {
	Key            string `json:"id"`
	LicName        string `json:"licenseName"`
	BusName        string `json:"businessName"`
	PremiseAddress string `json:"premiseAddress"`
}

//FlicRequest FlicRequest
type FlicRequest struct {
	ID          string
	Zip         string
	CustomerKey string
	APIKey      string
	Domain      string
}

//FlicAnalytics FlicAnalytics
type FlicAnalytics struct {
	CustomerKey string    `bigquery:"customer_key"`
	APIKey      string    `bigquery:"api_key"`
	Domain      string    `bigquery:"domain"`
	Entered     time.Time `bigquery:"entered"`
	Success     bool      `bigquery:"success"`
	SearchType  string    `bigquery:"type_search"`
}

//Manager Manager
type Manager interface {
	FindFlicListByZip(req *FlicRequest) (bool, *[]FlicList)
	FindFlicByKey(req *FlicRequest) (bool, *Flic)
	SetTableName(table string) bool
	ValidateUser(req *FlicRequest) bool
	InitialBqTableName() bool
}

//FlicManager FlicManager
type FlicManager struct {
	FlicDB         db.FlicDB
	Log            *lg.Logger
	AnalyticPusher ph.AnalyticPusher
	Puller         pu.Puller
	GcpProject     string //= "august-gantry-192521"
	DatasetName    string // = "ulboralabs"
	Table          string
}

//GetNew GetNew
func (m *FlicManager) GetNew() Manager {
	return m
}

//FindFlicListByZip FindFlicListByZip
func (m *FlicManager) FindFlicListByZip(req *FlicRequest) (bool, *[]FlicList) {
	var rtn = []FlicList{}
	var suc bool
	if m.ValidateUser(req) {
		suc = true
		var query = "SELECT key, lic_name, bus_name, premise_address " +
			" FROM " + m.GcpProject + "." + m.DatasetName + "." + m.Table +
			" WHERE premise_zip like @zip "
		var qp []bigquery.QueryParameter
		var par bigquery.QueryParameter
		par.Name = "zip"
		par.Value = req.Zip + "%"
		qp = append(qp, par)

		m.Log.Debug("query: ", query)

		recs := m.Puller.Pull(query, &qp)
		m.Log.Debug("res: ", *recs)

		for _, f := range *recs {
			var flic FlicList
			flic.Key = f[0].(string)
			flic.LicName = f[1].(string)
			flic.BusName = f[2].(string)
			flic.PremiseAddress = f[3].(string)
			rtn = append(rtn, flic)
		}
	}
	return suc, &rtn
}

//FindFlicByKey FindFlicByKey
func (m *FlicManager) FindFlicByKey(req *FlicRequest) (bool, *Flic) {
	var rtn Flic
	var suc bool
	if m.ValidateUser(req) {
		suc = true
		var query = "SELECT key, lic, exp_date, lic_name, bus_name, premise_address, mailing_address, phone " +
			" FROM " + m.GcpProject + "." + m.DatasetName + "." + m.Table +
			" WHERE key = @key "
		var qp []bigquery.QueryParameter
		var par bigquery.QueryParameter
		par.Name = "key"
		par.Value = req.ID
		qp = append(qp, par)

		m.Log.Debug("query: ", query)

		recs := m.Puller.Pull(query, &qp)
		m.Log.Debug("res: ", *recs)
		for _, f := range *recs {
			var flic Flic
			flic.Key = f[0].(string)
			flic.Lic = f[1].(string)
			flic.ExpDate = f[2].(time.Time)
			flic.LicName = f[3].(string)
			flic.BusName = f[4].(string)
			flic.PremiseAddress = f[5].(string)
			flic.MailingAddress = f[6].(string)
			flic.Phone = f[7].(string)
			rtn = flic
		}
	}
	return suc, &rtn
}

//SetTableName SetTableName
func (m *FlicManager) SetTableName(table string) bool {
	var rtn bool
	if table != "" {
		res := m.FlicDB.GetFlicTable()
		if res.ID != 0 {
			res.Name = table
			suc := m.FlicDB.SetFlicTable(res)
			if suc {
				m.Table = table
				rtn = suc
			}
		}
	}
	return rtn
}

//InitialBqTableName InitialBqTableName
func (m *FlicManager) InitialBqTableName() bool {
	var rtn bool
	res := m.FlicDB.GetFlicTable()
	m.Log.Debug("InitialBqTableName: ", res.Name)
	if res.ID != 0 {
		m.Table = res.Name
		rtn = true
		m.Log.Debug("InitialBqTableName Success: ", rtn)
	}
	return rtn
}

//ValidateUser ValidateUser
func (m *FlicManager) ValidateUser(req *FlicRequest) bool {
	var rtn bool
	m.Log.Debug("api key: ", req.APIKey)
	if req.APIKey != "" {
		res := m.FlicDB.GetUser(req.APIKey)
		m.Log.Debug("api key from db: ", *res)
		if res != nil && res.UserType == "api" && res.Enabled {
			rtn = true
		}
	} else if req.CustomerKey != "" {
		res := m.FlicDB.GetUser(req.CustomerKey)
		m.Log.Debug("user req: ", *req)
		m.Log.Debug("user res: ", *res)
		m.Log.Debug("user res.UserType: ", *&res.UserType)
		m.Log.Debug("user res.Domain: ", *&res.Domain)
		m.Log.Debug("user enabled: ", *&res.Enabled)
		if res != nil && res.UserType == "customer" && res.Domain == req.Domain && res.Enabled {
			rtn = true
		}
	}
	m.Log.Debug("user validated: ", rtn)
	var fl FlicAnalytics
	fl.APIKey = req.APIKey
	fl.CustomerKey = req.CustomerKey
	fl.Domain = req.Domain
	if req.ID != "" {
		fl.SearchType = searchTypeID
	} else if req.Zip != "" {
		fl.SearchType = searchTypeZip
	}
	fl.Entered = time.Now()
	fl.Success = rtn
	// go m.AnalyticPusher.Push(fl, analyticsTable)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(rc interface{}) {
		defer wg.Done()
		suc := m.AnalyticPusher.Push(rc, analyticsTable)
		m.Log.Debug("Analytic Push: ", suc)
	}(fl)
	wg.Wait()
	return rtn
}
