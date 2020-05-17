package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	ph "github.com/Ulbora/AnalyticPusher"
	pu "github.com/Ulbora/BigQueryPuller"
	mg "github.com/Ulbora/FlicService/managers"
	fdb "github.com/Ulbora/FlicService/mysqldb"
	lg "github.com/Ulbora/Level_Logger"
	db "github.com/Ulbora/dbinterface"
	mdb "github.com/Ulbora/dbinterface_mysql"
	"google.golang.org/api/option"
)

func TestFlicHandler_FindFlicListByZip(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	// var ayn ph.MockPusher
	//fh.AnalyticPusher = &ayn
	//ayn.MockPushSuccess = true
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
	getRow.Row = []string{"1", "some user", "456211111", "test.com", "customer", "1"}
	mydb.MockRow1 = &getRow

	dbi = &mydb

	udb.DB = dbi
	dbi.Connect()

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udb.Log = &l

	fm.FlicDB = &udb
	fm.Log = &l
	fh.Manager = &fm
	fh.Log = &l
	var bp pu.MockPuller
	var mres [][]bigquery.Value
	var mr1 []bigquery.Value
	var v1 bigquery.Value
	v1 = "158097011D35255"
	mr1 = append(mr1, v1)
	//v1 = "15-80-9701-1D-35255"
	//mr1 = append(mr1, v1)
	//v1 = time.Now()
	//mr1 = append(mr1, v1)
	v1 = "Bobs OUtdoors"
	mr1 = append(mr1, v1)
	v1 = "Bobs OUtdoors"
	mr1 = append(mr1, v1)
	v1 = "123 Bobs street, Bobtown PR"
	mr1 = append(mr1, v1)
	//v1 = "PO Box 123, Bobtown PR"
	//mr1 = append(mr1, v1)
	//v1 = "129-358-1234"
	//mr1 = append(mr1, v1)

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

	h := fh.GetNew()
	aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"zip":"12345"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", aJSON)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	r.Host = "test.com:8080"
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Origin", "http://test.com:8080")
	r.Header.Set("customer-key", "customer1234")
	w := httptest.NewRecorder()
	h.FindFlicListByZip(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 200 || w.Header().Get("Content-Type") != "application/json" {
		t.Fail()
	}
}

func TestFlicHandler_FindFlicListByZipAuth(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	// var ayn ph.MockPusher
	//fh.AnalyticPusher = &ayn
	//ayn.MockPushSuccess = true
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
	fh.Manager = &fm
	fh.Log = &l
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

	h := fh.GetNew()
	aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"zip":"12345"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", aJSON)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	r.Host = "test.com"
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Origin", "test.com")
	r.Header.Set("customer-key", "customer1234")
	w := httptest.NewRecorder()
	h.FindFlicListByZip(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 401 || w.Header().Get("Content-Type") != "application/json" {
		t.Fail()
	}
}

func TestFlicHandler_FindFlicListByZipHost(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	// var ayn ph.MockPusher
	// fh.AnalyticPusher = &ayn
	// ayn.MockPushSuccess = true
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
	fh.Manager = &fm
	fh.Log = &l
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

	h := fh.GetNew()
	aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"zip":"12345"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", aJSON)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	r.Host = "test.com"
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("customer-key", "customer1234")
	//r.Header.Set("Origin", "test.com")
	w := httptest.NewRecorder()
	h.FindFlicListByZip(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 401 || w.Header().Get("Content-Type") != "application/json" {
		t.Fail()
	}
}

func TestFlicHandler_FindFlicListByZipBadBody(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	var ayn ph.MockPusher
	//fh.AnalyticPusher = &ayn
	ayn.MockPushSuccess = true
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
	fh.Manager = &fm
	fh.Log = &l
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

	h := fh.GetNew()
	//aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"zip":"12345"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", nil)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	r.Host = "test.com"
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.FindFlicListByZip(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 400 {
		t.Fail()
	}
}

func TestFlicHandler_FindFlicListByZipMedia(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	var ayn ph.MockPusher
	//fh.AnalyticPusher = &ayn
	ayn.MockPushSuccess = true
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
	fh.Manager = &fm
	fh.Log = &l
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

	h := fh.GetNew()
	aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"zip":"12345"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", aJSON)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	r.Host = "test.com"
	//r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.FindFlicListByZip(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 415 {
		t.Fail()
	}
}

type testObj struct {
	Valid bool   `json:"valid"`
	Code  string `json:"code"`
}

func TestFlicHandler_processBody(t *testing.T) {
	var oh FlicHandler
	var l lg.Logger
	oh.Log = &l
	var robj testObj
	robj.Valid = true
	robj.Code = "3"
	// var res http.Response
	// res.Body = ioutil.NopCloser(bytes.NewBufferString(`{"valid":true, "code":"1"}`))
	var sURL = "http://localhost/test"
	aJSON, _ := json.Marshal(robj)
	r, _ := http.NewRequest("POST", sURL, bytes.NewBuffer(aJSON))
	var obj testObj
	suc, _ := oh.processBody(r, nil)
	if suc || obj.Valid != false || obj.Code != "" {
		t.Fail()
	}
}

func TestFlicHandler_FindFlicByKey(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	// var ayn ph.MockPusher
	// fh.AnalyticPusher = &ayn
	// ayn.MockPushSuccess = true
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
	getRow.Row = []string{"1", "some user", "456211111", "test.com", "customer", "1"}
	mydb.MockRow1 = &getRow

	dbi = &mydb

	udb.DB = dbi
	dbi.Connect()

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udb.Log = &l

	fm.FlicDB = &udb
	fm.Log = &l
	fh.Manager = &fm
	fh.Log = &l
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

	h := fh.GetNew()
	aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"id":"12345"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", aJSON)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	r.Host = "test.com:8080"
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Origin", "http://test.com:8089")
	r.Header.Set("customer-key", "customer1234")
	w := httptest.NewRecorder()
	h.FindFlicByKey(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 200 || w.Header().Get("Content-Type") != "application/json" {
		t.Fail()
	}
}

func TestFlicHandler_FindFlicByKeyAuth(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	// var ayn ph.MockPusher
	// fh.AnalyticPusher = &ayn
	// ayn.MockPushSuccess = true
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
	fh.Manager = &fm
	fh.Log = &l
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

	h := fh.GetNew()
	aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"id":"12345"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", aJSON)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	r.Host = "test.com"
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Origin", "test.com")
	r.Header.Set("customer-key", "customer1234")
	w := httptest.NewRecorder()
	h.FindFlicByKey(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 401 || w.Header().Get("Content-Type") != "application/json" {
		t.Fail()
	}
}

func TestFlicHandler_FindFlicByKeyHost(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	// var ayn ph.MockPusher
	// fh.AnalyticPusher = &ayn
	// ayn.MockPushSuccess = true
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
	fh.Manager = &fm
	fh.Log = &l
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

	h := fh.GetNew()
	aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"id":"12345"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", aJSON)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	r.Host = "test.com"
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("customer-key", "customer1234")
	//r.Header.Set("Origin", "test.com")
	w := httptest.NewRecorder()
	h.FindFlicByKey(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 401 || w.Header().Get("Content-Type") != "application/json" {
		t.Fail()
	}
}

func TestFlicHandler_FindFlicByKeyBody(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	var ayn ph.MockPusher
	//fh.AnalyticPusher = &ayn
	ayn.MockPushSuccess = true
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
	fh.Manager = &fm
	fh.Log = &l
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

	h := fh.GetNew()
	//aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"id":"12345"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", nil)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	r.Host = "test.com"
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Origin", "test.com")
	w := httptest.NewRecorder()
	h.FindFlicByKey(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 400 {
		t.Fail()
	}
}

func TestFlicHandler_FindFlicByKeyMedia(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	var ayn ph.MockPusher
	//fh.AnalyticPusher = &ayn
	ayn.MockPushSuccess = true
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
	fh.Manager = &fm
	fh.Log = &l
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

	h := fh.GetNew()
	aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"id":"12345"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", aJSON)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	r.Host = "test.com"
	//r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Origin", "test.com")
	w := httptest.NewRecorder()
	h.FindFlicByKey(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 415 {
		t.Fail()
	}
}

func TestFlicHandler_SetFlicTable(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	var ayn ph.MockPusher
	//fh.AnalyticPusher = &ayn
	ayn.MockPushSuccess = true
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

	dbi = &mydb

	udb.DB = dbi
	dbi.Connect()

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udb.Log = &l

	fm.FlicDB = &udb
	fm.Log = &l
	fh.Manager = &fm
	fh.Log = &l

	fm.AnalyticPusher = &ayn

	h := fh.GetNew()
	aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"name":"new_table_name"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", aJSON)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	//r.Host = "test.com"
	r.Header.Set("Content-Type", "application/json")
	//r.Header.Set("Origin", "test.com")
	r.Header.Set("api-key", "456211111")
	w := httptest.NewRecorder()
	h.SetFlicTable(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 200 {
		t.Fail()
	}
}

func TestFlicHandler_SetFlicTableApiKey(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	var ayn ph.MockPusher
	//fh.AnalyticPusher = &ayn
	ayn.MockPushSuccess = true
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
	getRow.Row = []string{"1", "some user", "456211111aa", "", "api", "0"}
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
	fh.Manager = &fm
	fh.Log = &l

	fm.AnalyticPusher = &ayn

	h := fh.GetNew()
	aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"name":"new_table_name"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", aJSON)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	//r.Host = "test.com"
	r.Header.Set("Content-Type", "application/json")
	//r.Header.Set("Origin", "test.com")
	r.Header.Set("api-key", "456211111")
	w := httptest.NewRecorder()
	h.SetFlicTable(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 401 {
		t.Fail()
	}
}

func TestFlicHandler_SetFlicTableBody(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	var ayn ph.MockPusher
	//fh.AnalyticPusher = &ayn
	ayn.MockPushSuccess = true
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

	dbi = &mydb

	udb.DB = dbi
	dbi.Connect()

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udb.Log = &l

	fm.FlicDB = &udb
	fm.Log = &l
	fh.Manager = &fm
	fh.Log = &l

	fm.AnalyticPusher = &ayn

	h := fh.GetNew()
	//aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"name":"new_table_name"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", nil)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	//r.Host = "test.com"
	r.Header.Set("Content-Type", "application/json")
	//r.Header.Set("Origin", "test.com")
	r.Header.Set("api-key", "456211111")
	w := httptest.NewRecorder()
	h.SetFlicTable(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 400 {
		t.Fail()
	}
}

func TestFlicHandler_SetFlicTableMedia(t *testing.T) {
	var fh FlicHandler
	var fm mg.FlicManager
	var ayn ph.MockPusher
	//fh.AnalyticPusher = &ayn
	ayn.MockPushSuccess = true
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

	dbi = &mydb

	udb.DB = dbi
	dbi.Connect()

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	udb.Log = &l

	fm.FlicDB = &udb
	fm.Log = &l
	fh.Manager = &fm
	fh.Log = &l

	fm.AnalyticPusher = &ayn

	h := fh.GetNew()
	aJSON := ioutil.NopCloser(bytes.NewBufferString(`{"name":"new_table_name"}`))
	//aJSON, _ := json.Marshal(robj)
	//fmt.Println("aJSON: ", aJSON)
	r, _ := http.NewRequest("POST", "/ffllist", aJSON)
	//r, _ := http.NewRequest("POST", "/ffllist", nil)
	//r.Host = "test.com"
	//r.Header.Set("Content-Type", "application/json")
	//r.Header.Set("Origin", "test.com")
	r.Header.Set("api-key", "456211111")
	w := httptest.NewRecorder()
	h.SetFlicTable(w, r)
	hd := w.Header()
	fmt.Println("w content type", hd.Get("Content-Type"))
	if w.Code != 415 {
		t.Fail()
	}
}
