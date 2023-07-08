package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"supply/svc"
)

func Initialize() {
	http.HandleFunc("/space", HandleService)
	http.HandleFunc("/space/", HandleService)
}

func HandleService(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	urlSplit := strings.Split(urlPath, "/")
	urlSplitCount := len(urlSplit)
	id := ""

	if urlSplitCount > 4 || (urlSplitCount == 4 && urlSplit[3] != "") || urlSplitCount < 2 {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	} else if urlSplitCount == 2 {
		id = ""
	} else {
		id = urlSplit[2]
	}

	var resp any
	var err error

	if id == "" && r.Method == http.MethodPost {
		resp, err = svc.Create(w, r) // 			POST 		/space
	} else if id == "" && r.Method == http.MethodGet {
		resp, err = svc.GetAll(w, r) // 			GET 		/space
	} else if id != "" && r.Method == http.MethodGet {
		resp, err = svc.GetOne(w, r, id) // 	GET 		/space/:id
	} else if id != "" && r.Method == http.MethodPatch {
		resp, err = svc.Update(w, r, id) // 	PATCH 	/space/:id
	} else if id != "" && r.Method == http.MethodDelete {
		resp, err = svc.Delete(w, r, id) // 	DELETE 	/space/:id
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defaultReturn(w, r, resp, err)
}

func defaultReturn(w http.ResponseWriter, r *http.Request, resp any, err error) {
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
