package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	b64 "encoding/base64"

	"cloud.google.com/go/bigquery"
	ph "github.com/Ulbora/AnalyticPusher"
	pu "github.com/Ulbora/BigQueryPuller"
	flh "github.com/Ulbora/FlicService/handlers"
	mg "github.com/Ulbora/FlicService/managers"
	fdb "github.com/Ulbora/FlicService/mysqldb"
	lg "github.com/Ulbora/Level_Logger"
	mdb "github.com/Ulbora/dbinterface_mysql"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

func main() {
	var fh flh.FlicHandler
	var fm mg.FlicManager
	var udb fdb.UserDB
	var mydb mdb.MyDB
	var flicDbHost string
	var flicDbUser string
	var flicDbPassword string
	var flicDbName string
	var l lg.Logger
	l.LogLevel = lg.AllLevel

	var gcpCreds string

	if os.Getenv("FLIC_DB_HOST") != "" {
		flicDbHost = os.Getenv("FLIC_DB_HOST")
	} else {
		flicDbHost = "localhost:3306"
	}

	if os.Getenv("FLIC_DB_USER") != "" {
		flicDbUser = os.Getenv("FLIC_DB_USER")
	} else {
		flicDbUser = "admin"
	}

	if os.Getenv("FLIC_DB_PASSWORD") != "" {
		flicDbPassword = os.Getenv("FLIC_DB_PASSWORD")
	} else {
		flicDbPassword = "admin"
	}

	if os.Getenv("FLIC_DB_DATABASE") != "" {
		flicDbName = os.Getenv("FLIC_DB_DATABASE")
	} else {
		flicDbName = "flic_service"
	}
	mydb.Host = flicDbHost         // "localhost:3306"
	mydb.User = flicDbUser         // "admin"
	mydb.Password = flicDbPassword // "admin"
	mydb.Database = flicDbName     //"flic_service"

	if os.Getenv("GCP_CREDS") != "" {
		creds, _ := b64.StdEncoding.DecodeString(os.Getenv("GCP_CREDS"))
		err := ioutil.WriteFile("creds.json", creds, 0644)
		l.Debug("creds in main err: ", err)
		gcpCreds = "./creds.json"
	} else {
		gcpCreds = "../gcpCreds.json"
	}

	udb.Log = &l
	udb.DB = &mydb
	udb.DB.Connect()
	fm.Log = &l
	fm.FlicDB = &udb
	fh.Log = &l
	var bp pu.BigQueryPuller
	fm.GcpProject = "august-gantry-192521"
	fm.DatasetName = "ulboralabs"
	fm.Table = "flic_May_5_2020_18_28_26"

	fm.InitialBqTableName()

	ctx := context.Background()
	bp.SetContext(ctx)
	client, err := bigquery.NewClient(ctx, fm.GcpProject, option.WithCredentialsFile(gcpCreds))
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
	fh.Manager = &fm

	router := mux.NewRouter()
	port := "3000"
	envPort := os.Getenv("PORT")
	if envPort != "" {
		portInt, _ := strconv.Atoi(envPort)
		if portInt != 0 {
			port = envPort
		}
	}
	h := fh.GetNew()
	l.LogLevel = lg.OffLevel

	router.HandleFunc("/rs/findByZip", h.FindFlicListByZip).Methods("POST")
	router.HandleFunc("/rs/findById", h.FindFlicByKey).Methods("POST")
	router.HandleFunc("/rs/table", h.SetFlicTable).Methods("POST")
	router.HandleFunc("/rs/loglevel", h.SetLogLevel).Methods("POST")

	fmt.Println("Starting Server on " + port)
	http.ListenAndServe(":"+port, router)

}

// go mod init github.com/Ulbora/FlicService
