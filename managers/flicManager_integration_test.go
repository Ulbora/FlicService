// +build integration move to top

package manager

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/bigquery"
	ph "github.com/Ulbora/AnalyticPusher"
	pu "github.com/Ulbora/BigQueryPuller"

	fdb "github.com/Ulbora/FlicService/mysqldb"
	lg "github.com/Ulbora/Level_Logger"
	db "github.com/Ulbora/dbinterface"
	mdb "github.com/Ulbora/dbinterface_mysql"
	"google.golang.org/api/option"
)

func TestFlicManageri_FindFlicListByZip(t *testing.T) {
	var fm FlicManager
	var db db.Database
	var udb fdb.UserDB
	var mydb mdb.MyDB
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

	db = &mydb

	udb.DB = db
	db.Connect()

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udb.Log = &l

	fm.FlicDB = &udb
	fm.Log = &l
	var bp pu.BigQueryPuller
	fm.GcpProject = "august-gantry-192521"
	fm.DatasetName = "ulboralabs"
	fm.Table = "flic_May_5_2020_18_28_26"
	ctx := context.Background()
	bp.SetContext(ctx)
	client, err := bigquery.NewClient(ctx, fm.GcpProject, option.WithCredentialsFile("../../gcpCreds.json"))
	if err != nil {
		fmt.Println("bq err: ", err)
	}
	bp.SetClient(client)
	fm.Puller = &bp

	var ayn ph.Pusher
	ayn.GcpProject = "august-gantry-192521"
	ayn.Client = client
	ayn.Ctx = ctx
	ayn.DatasetName = "ulboralabs"
	fm.AnalyticPusher = &ayn

	var req FlicRequest
	req.APIKey = "61616dfggdf5g64gf4"
	req.Zip = "30134"

	f := fm.GetNew()
	res := f.FindFlicListByZip(&req)
	if len(*res) == 0 {
		t.Fail()
	} else {
		fmt.Println((*res)[0].LicName)
		fmt.Println((*res)[0].BusName)
	}

}

func TestFlicManageri_FindFlicByKey(t *testing.T) {
	var fm FlicManager
	var db db.Database
	var udb fdb.UserDB
	var mydb mdb.MyDB
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

	db = &mydb

	udb.DB = db
	db.Connect()

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udb.Log = &l

	fm.FlicDB = &udb
	fm.Log = &l
	var bp pu.BigQueryPuller
	fm.GcpProject = "august-gantry-192521"
	fm.DatasetName = "ulboralabs"
	fm.Table = "flic_May_5_2020_18_28_26"
	ctx := context.Background()
	bp.SetContext(ctx)
	client, err := bigquery.NewClient(ctx, fm.GcpProject, option.WithCredentialsFile("../../gcpCreds.json"))
	if err != nil {
		fmt.Println("bq err: ", err)
	}
	bp.SetClient(client)
	fm.Puller = &bp

	var ayn ph.Pusher
	ayn.GcpProject = "august-gantry-192521"
	ayn.Client = client
	ayn.Ctx = ctx
	ayn.DatasetName = "ulboralabs"
	fm.AnalyticPusher = &ayn

	var req FlicRequest
	req.APIKey = "61616dfggdf5g64gf4"
	req.ID = "158097011D35284"

	f := fm.GetNew()
	res := f.FindFlicByKey(&req)
	if (*res).Key == "" {
		t.Fail()
	} else {
		fmt.Println(res.LicName)
		fmt.Println(res.BusName)
	}

}
