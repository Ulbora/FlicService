package handlers

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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	//ph "github.com/Ulbora/AnalyticPusher"
	mg "github.com/Ulbora/FlicService/managers"
	lg "github.com/Ulbora/Level_Logger"
)

//FlicHandler FlicHandler
type FlicHandler struct {
	Manager mg.Manager
	Log     *lg.Logger
	//AnalyticPusher ph.AnalyticPusher
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

//Flic Flic
type Flic struct {
	Key            string `json:"id"`
	Lic            string `json:"license"`
	ExpDate        string `json:"expDate"`
	LicName        string `json:"licenseName"`
	BusName        string `json:"businessName"`
	PremiseAddress string `json:"premiseAddress"`
	Address        string `json:"address"`
	City           string `json:"city"`
	State          string `json:"state"`
	PremiseZip     string `json:"premiseZip"`
	MailingAddress string `json:"mailingAddress"`
	Phone          string `json:"phone"`
}

//GetNew GetNew
func (h *FlicHandler) GetNew() Handlers {
	return h
}

//FindFlicListByZip FindFlicListByZip
func (h *FlicHandler) FindFlicListByZip(w http.ResponseWriter, r *http.Request) {
	//SetupResponse(&w, r)
	// if (r).Method == "OPTIONS" {
	// 	//w.WriteHeader(http.StatusOK)
	// 	//h.Log.Debug("In preflight: ", w.WriteHeader)

	// } else {

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
			colonInd := strings.Index(host, ":")
			if colonInd > 0 {
				h.Log.Debug("stripped host: ", host[0:colonInd])
				host = host[0:colonInd]
			}

			colonInd = strings.LastIndex(origin, ":")
			//h.Log.Debug("stripped origin port: ", origin[colonInd:])
			if colonInd > 0 && len(origin[colonInd:]) == 5 {
				h.Log.Debug("stripped origin 1: ", origin[0:colonInd])
				origin = origin[0:colonInd]
			}
			slashInd := strings.LastIndex(origin, "/")
			origin = origin[slashInd+1:]
			h.Log.Debug("stripped origin 2: ", origin)

			wwwInd := strings.Index(origin, "www.")
			if wwwInd > -1 {
				origin = origin[wwwInd+4:]
				h.Log.Debug("stripped origin 3: ", origin)
			}

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
			suc, flicList := h.Manager.FindFlicListByZip(&flicReq)
			if suc {
				h.Log.Debug("flicList: ", *flicList)
				w.WriteHeader(http.StatusOK)
				resJSON, _ := json.Marshal(flicList)
				fmt.Fprint(w, string(resJSON))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}

		}
	}
	//	}

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
			colonInd := strings.Index(host, ":")
			if colonInd > 0 {
				h.Log.Debug("stripped host: ", host[0:colonInd])
				host = host[0:colonInd]
			}

			colonInd = strings.LastIndex(origin, ":")
			//h.Log.Debug("stripped origin port: ", origin[colonInd:])
			if colonInd > 0 && len(origin[colonInd:]) == 5 {
				h.Log.Debug("stripped origin 1: ", origin[0:colonInd])
				origin = origin[0:colonInd]
			}
			slashInd := strings.LastIndex(origin, "/")
			origin = origin[slashInd+1:]
			h.Log.Debug("stripped origin 2: ", origin)

			wwwInd := strings.Index(origin, "www.")
			if wwwInd > -1 {
				origin = origin[wwwInd+4:]
				h.Log.Debug("stripped origin 3: ", origin)
			}

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
			suc, flic := h.Manager.FindFlicByKey(&flicReq)
			if suc {
				layoutUS := "January 2, 2006"

				var flc Flic
				flc.Key = flic.Key
				flc.Lic = flic.Lic
				flc.BusName = flic.BusName
				flc.LicName = flic.LicName
				flc.PremiseAddress = flic.PremiseAddress
				flc.Address = flic.Address
				flc.City = flic.City
				flc.State = flic.State
				flc.PremiseZip = flic.PremiseZip
				flc.MailingAddress = flic.MailingAddress
				flc.Phone = flic.Phone
				flc.ExpDate = flic.ExpDate.Format(layoutUS)
				dstr := flic.ExpDate.Format(layoutUS)
				h.Log.Debug("date : ", flic.ExpDate)
				h.Log.Debug("date str: ", dstr)
				w.WriteHeader(http.StatusOK)
				resJSON, _ := json.Marshal(flc)
				fmt.Fprint(w, string(resJSON))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}

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
