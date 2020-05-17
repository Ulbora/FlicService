package manager

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	ph "github.com/Ulbora/AnalyticPusher"
	pu "github.com/Ulbora/BigQueryPuller"

	fdb "github.com/Ulbora/FlicService/mysqldb"
	lg "github.com/Ulbora/Level_Logger"
	db "github.com/Ulbora/dbinterface"
	mdb "github.com/Ulbora/dbinterface_mysql"
	"google.golang.org/api/option"
)

func TestFlicManager_FindFlicListByZip(t *testing.T) {
	var fm FlicManager
	var dbi db.Database
	var udb fdb.UserDB
	var mydb mdb.MyDBMock
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

	var mTestRow db.DbRow
	mTestRow.Row = []string{"1"}
	mydb.MockTestRow = &mTestRow

	var getRow db.DbRow
	getRow.Row = []string{"1", "some user", "456211111", "somedomain.com", "customer", "1"}
	mydb.MockRow1 = &getRow

	var rows [][]string
	row1 := []string{"1", "some_bq_table"}
	rows = append(rows, row1)
	var dbrows db.DbRows
	dbrows.Rows = rows
	mydb.MockRows1 = &dbrows

	dbi = &mydb

	udb.DB = dbi
	dbi.Connect()

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udb.Log = &l

	fm.FlicDB = &udb
	fm.Log = &l
	var bp pu.MockPuller
	var mres [][]bigquery.Value
	var mr1 []bigquery.Value
	var v1 bigquery.Value
	v1 = "158097011D35255"
	mr1 = append(mr1, v1)
	v1 = "Bobs OUtdoors"
	mr1 = append(mr1, v1)
	v1 = "Bobs OUtdoors"
	mr1 = append(mr1, v1)
	v1 = "123 Bobs street, Bobtown PR"
	mr1 = append(mr1, v1)

	mres = append(mres, mr1)
	bp.MockResp = mres

	fm.GcpProject = "august-gantry-192521"
	fm.DatasetName = "ulboralabs"
	//fm.Table = "flic_May_5_2020_18_28_26"

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
	req.CustomerKey = "61616dfggdf5g64gf4"
	req.Zip = "30134"
	req.Domain = "somedomain.com"

	f := fm.GetNew()
	f.SetTableName("flic_May_5_2020_18_28_26")
	_, res := f.FindFlicListByZip(&req)
	if len(*res) == 0 {
		t.Fail()
	} else {
		fmt.Println((*res)[0].LicName)
		fmt.Println((*res)[0].BusName)
	}

}

func TestFlicManager_FindFlicByKey(t *testing.T) {
	var fm FlicManager
	var dbi db.Database
	var udb fdb.UserDB
	var mydb mdb.MyDBMock
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

	var mTestRow db.DbRow
	mTestRow.Row = []string{"1"}
	mydb.MockTestRow = &mTestRow

	var getRow db.DbRow
	getRow.Row = []string{"1", "some user", "456211111", "", "api", "1"}
	mydb.MockRow1 = &getRow

	dbi = &mydb

	udb.DB = dbi
	dbi.Connect()

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udb.Log = &l

	fm.FlicDB = &udb
	fm.Log = &l
	var bp pu.MockPuller
	var mres [][]bigquery.Value
	var mr1 []bigquery.Value
	var v1 bigquery.Value
	v1 = "158097011D35255"
	mr1 = append(mr1, v1)
	v1 = "15-80-9701-1D-35255"
	mr1 = append(mr1, v1)
	v1 = time.Now()
	mr1 = append(mr1, v1)
	v1 = "Bobs OUtdoors"
	mr1 = append(mr1, v1)
	v1 = "Bobs OUtdoors"
	mr1 = append(mr1, v1)
	v1 = "123 Bobs street, Bobtown PR"
	mr1 = append(mr1, v1)
	v1 = "123 Bobs street"
	mr1 = append(mr1, v1)
	v1 = "Bobtown"
	mr1 = append(mr1, v1)
	v1 = "PR"
	mr1 = append(mr1, v1)
	v1 = "12345"
	mr1 = append(mr1, v1)
	v1 = "PO Box 123, Bobtown PR"
	mr1 = append(mr1, v1)
	v1 = "129-358-1234"
	mr1 = append(mr1, v1)

	mres = append(mres, mr1)
	bp.MockResp = mres

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
	_, res := f.FindFlicByKey(&req)
	if (*res).Key == "" {
		t.Fail()
	} else {
		fmt.Println(res.LicName)
		fmt.Println(res.BusName)
	}
}

func TestFlicManager_SetBqTable(t *testing.T) {
	var fm FlicManager
	var dbi db.Database
	var udb fdb.UserDB
	var mydb mdb.MyDBMock
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

	var mTestRow db.DbRow
	mTestRow.Row = []string{"1"}
	mydb.MockTestRow = &mTestRow

	var getRow db.DbRow
	getRow.Row = []string{"1", "some user", "456211111", "", "api", "1"}
	mydb.MockRow1 = &getRow

	var rows [][]string
	row1 := []string{"1", "some_bq_table"}
	rows = append(rows, row1)
	var dbrows db.DbRows
	dbrows.Rows = rows
	mydb.MockRows1 = &dbrows

	mydb.MockUpdateSuccess1 = true

	dbi = &mydb

	udb.DB = dbi
	dbi.Connect()

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udb.Log = &l

	fm.FlicDB = &udb
	fm.Log = &l
	var bp pu.MockPuller
	var mres [][]bigquery.Value
	var mr1 []bigquery.Value
	var v1 bigquery.Value
	v1 = "158097011D35255"
	mr1 = append(mr1, v1)
	v1 = "15-80-9701-1D-35255"
	mr1 = append(mr1, v1)
	v1 = time.Now()
	mr1 = append(mr1, v1)
	v1 = "Bobs OUtdoors"
	mr1 = append(mr1, v1)
	v1 = "Bobs OUtdoors"
	mr1 = append(mr1, v1)
	v1 = "123 Bobs street, Bobtown PR"
	mr1 = append(mr1, v1)
	v1 = "PO Box 123, Bobtown PR"
	mr1 = append(mr1, v1)
	v1 = "129-358-1234"
	mr1 = append(mr1, v1)

	mres = append(mres, mr1)
	bp.MockResp = mres

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
	suc := f.SetTableName("test_table")
	if !suc {
		t.Fail()
	}
}

func TestFlicManager_InitialBqTableName(t *testing.T) {
	var fm FlicManager
	var dbi db.Database
	var udb fdb.UserDB
	var mydb mdb.MyDBMock
	mydb.Host = "localhost:3306"
	mydb.User = "admin"
	mydb.Password = "admin"
	mydb.Database = "flic_service"

	var mTestRow db.DbRow
	mTestRow.Row = []string{"1"}
	mydb.MockTestRow = &mTestRow

	var getRow db.DbRow
	getRow.Row = []string{"1", "some user", "456211111", "somedomain.com", "customer", "1"}
	mydb.MockRow1 = &getRow

	var rows [][]string
	row1 := []string{"1", "some_bq_table"}
	rows = append(rows, row1)
	var dbrows db.DbRows
	dbrows.Rows = rows
	mydb.MockRows1 = &dbrows

	dbi = &mydb

	udb.DB = dbi
	dbi.Connect()

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udb.Log = &l

	fm.FlicDB = &udb
	fm.Log = &l

	f := fm.GetNew()
	suc := f.InitialBqTableName()
	if !suc {
		t.Fail()
	}
}
