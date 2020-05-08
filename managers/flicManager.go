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
	"time"

	"cloud.google.com/go/bigquery"
	pu "github.com/Ulbora/BigQueryPuller"
	db "github.com/Ulbora/FlicService/mysqldb"
	lg "github.com/Ulbora/Level_Logger"
)

//Flic Flic
type Flic struct {
	Key            string    // `bigquery:"key"`
	Lic            string    //`bigquery:"lic"`
	ExpDate        time.Time //`bigquery:"exp_date"`
	LicName        string    //`bigquery:"lic_name"`
	BusName        string    //`bigquery:"bus_name"`
	PremiseAddress string    //`bigquery:"premise_address"`
	PremiseZip     string    //`bigquery:"premise_zip"`
	MailingAddress string    // `bigquery:"mailing_address"`
	Phone          string    //`bigquery:"phone"`
}

//FlicRequest FlicRequest
type FlicRequest struct {
	ID          string
	Zip         string
	CustomerKey string
	APIKey      string
	Domain      string
}

//Manager Manager
type Manager interface {
	FindFlicListByZip(req *FlicRequest) *[]Flic
	FindFlicByKey(req *FlicRequest) *Flic
}

//FlicManager FlicManager
type FlicManager struct {
	FlicDB db.FlicDB
	Log    *lg.Logger
	//Pusher
	Puller      pu.Puller
	GcpProject  string //= "august-gantry-192521"
	DatasetName string // = "ulboralabs"
	Table       string
}

//GetNew GetNew
func (m *FlicManager) GetNew() Manager {
	return m
}

//FindFlicListByZip FindFlicListByZip
func (m *FlicManager) FindFlicListByZip(req *FlicRequest) *[]Flic {
	var rtn []Flic
	if m.validateUser(req) {
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
			var flic Flic
			flic.Key = f[0].(string)
			flic.LicName = f[1].(string)
			flic.BusName = f[2].(string)
			flic.PremiseAddress = f[3].(string)
			rtn = append(rtn, flic)
		}
	}
	return &rtn
}

//FindFlicByKey FindFlicByKey
func (m *FlicManager) FindFlicByKey(req *FlicRequest) *Flic {
	var rtn Flic
	if m.validateUser(req) {
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
	return &rtn
}

func (m *FlicManager) validateUser(req *FlicRequest) bool {
	var rtn bool
	if req.APIKey != "" {
		res := m.FlicDB.GetUser(req.APIKey)
		if res != nil && res.UserType == "api" {
			rtn = true
		}
	} else if req.CustomerKey != "" {
		res := m.FlicDB.GetUser(req.CustomerKey)
		m.Log.Debug("user req: ", *req)
		m.Log.Debug("user res: ", *res)
		m.Log.Debug("user res.UserType: ", *&res.UserType)
		m.Log.Debug("user res.Domain: ", *&res.Domain)
		if res != nil && res.UserType == "customer" && res.Domain == req.Domain {
			rtn = true
		}
	}
	m.Log.Debug("user validated: ", rtn)
	return rtn
}
