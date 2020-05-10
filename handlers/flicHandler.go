package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	ph "github.com/Ulbora/AnalyticPusher"
	mg "github.com/Ulbora/FlicService/managers"
	lg "github.com/Ulbora/Level_Logger"
)

//FlicHandler FlicHandler
type FlicHandler struct {
	Manager        mg.Manager
	Log            *lg.Logger
	AnalyticPusher ph.AnalyticPusher
}

//Request Request
type Request struct {
	Zip string `json:"zip"`
	ID  string `json:"id"`
}

//TableRequest TableRequest
type TableRequest struct {
	Name string `json:"name"`
}

//GetNew GetNew
func (h *FlicHandler) GetNew() Handlers {
	return h
}

//FindFlicListByZip FindFlicListByZip
func (h *FlicHandler) FindFlicListByZip(w http.ResponseWriter, r *http.Request) {
	h.setContentType(w)
	h.Log.Debug("Host: ", r.Host)
	h.Log.Debug("URL: ", r.URL)
	h.Log.Debug("origin: ", r.Header.Get("Origin"))
	contOk := h.checkContent(r)
	h.Log.Debug("conOk: ", contOk)
	if !contOk {
		http.Error(w, "json required", http.StatusUnsupportedMediaType)
	} else {
		var req Request
		suc, err := h.processBody(r, &req)
		h.Log.Debug("suc: ", suc)
		h.Log.Debug("req: ", req)
		h.Log.Debug("err: ", err)
		if !suc && err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			origin := r.Header.Get("Origin")
			host := r.Host
			apiKey := r.Header.Get("api-key")
			customerKey := r.Header.Get("customer-key")
			var flicReq mg.FlicRequest
			flicReq.APIKey = apiKey
			flicReq.CustomerKey = customerKey
			flicReq.Zip = req.Zip
			if origin != "" {
				flicReq.Domain = origin
			} else if host != "" {
				flicReq.Domain = host
			}
			flicList := h.Manager.FindFlicListByZip(&flicReq)
			w.WriteHeader(http.StatusOK)
			resJSON, _ := json.Marshal(flicList)
			fmt.Fprint(w, string(resJSON))
		}
	}
}

//FindFlicByKey FindFlicByKey
func (h *FlicHandler) FindFlicByKey(w http.ResponseWriter, r *http.Request) {
	h.setContentType(w)
	h.Log.Debug("Host: ", r.Host)
	h.Log.Debug("URL: ", r.URL)
	h.Log.Debug("origin: ", r.Header.Get("Origin"))
	contOk := h.checkContent(r)
	h.Log.Debug("conOk: ", contOk)
	if !contOk {
		http.Error(w, "json required", http.StatusUnsupportedMediaType)
	} else {
		var req Request
		suc, err := h.processBody(r, &req)
		h.Log.Debug("suc: ", suc)
		h.Log.Debug("req: ", req)
		h.Log.Debug("err: ", err)
		if !suc && err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			origin := r.Header.Get("Origin")
			host := r.Host
			apiKey := r.Header.Get("api-key")
			customerKey := r.Header.Get("customer-key")
			var flicReq mg.FlicRequest
			flicReq.APIKey = apiKey
			flicReq.CustomerKey = customerKey
			flicReq.ID = req.ID
			if origin != "" {
				flicReq.Domain = origin
			} else if host != "" {
				flicReq.Domain = host
			}
			flic := h.Manager.FindFlicByKey(&flicReq)
			w.WriteHeader(http.StatusOK)
			resJSON, _ := json.Marshal(flic)
			fmt.Fprint(w, string(resJSON))
		}
	}
}

//SetFlicTable SetFlicTable
func (h *FlicHandler) SetFlicTable(w http.ResponseWriter, r *http.Request) {
	contOk := h.checkContent(r)
	h.Log.Debug("conOk: ", contOk)
	if !contOk {
		http.Error(w, "json required", http.StatusUnsupportedMediaType)
	} else {
		var req TableRequest
		suc, err := h.processBody(r, &req)
		h.Log.Debug("suc: ", suc)
		h.Log.Debug("req: ", req)
		h.Log.Debug("err: ", err)
		if !suc && err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			apiKey := r.Header.Get("api-key")
			var flicReq mg.FlicRequest
			flicReq.APIKey = apiKey
			suc := h.Manager.ValidateUser(&flicReq)
			if suc {
				h.Manager.SetTableName(req.Name)
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	}
}

//CheckContent CheckContent
func (h *FlicHandler) checkContent(r *http.Request) bool {
	var rtn bool
	cType := r.Header.Get("Content-Type")
	if cType == "application/json" {
		rtn = true
	}
	return rtn
}

//SetContentType SetContentType
func (h *FlicHandler) setContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

//ProcessBody ProcessBody
func (h *FlicHandler) processBody(r *http.Request, obj interface{}) (bool, error) {
	var suc bool
	var err error
	if r.Body != nil {
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(obj)
		if err != nil {
			h.Log.Error("Decode Error: ", err.Error())
		} else {
			suc = true
		}
	} else {
		err = errors.New("Bad Body")
	}
	return suc, err
}
