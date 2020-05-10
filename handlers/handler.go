package handlers

import "net/http"

//Handlers Handlers
type Handlers interface {
	FindFlicListByZip(w http.ResponseWriter, r *http.Request)
	FindFlicByKey(w http.ResponseWriter, r *http.Request)
	SetFlicTable(w http.ResponseWriter, r *http.Request)
	SetLogLevel(w http.ResponseWriter, r *http.Request)
}
