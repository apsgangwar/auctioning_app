package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"auction/svc"
)

func Initialize() {
	http.HandleFunc("/auction", HandleService)
	http.HandleFunc("/auction/", HandleService)
}

func HandleService(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	urlSplit := strings.Split(urlPath, "/")
	urlSplitCount := len(urlSplit)
	id := ""

	if urlSplitCount > 5 || (urlSplitCount == 5 && urlSplit[4] != "" && urlSplit[3] != "bid") || (urlSplitCount == 4 && urlSplit[3] != "" && urlSplit[3] != "bid") || urlSplitCount < 2 {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	} else if urlSplitCount == 2 {
		id = ""
	} else {
		id = urlSplit[2]
	}

	var resp any
	var err error

	if urlSplitCount > 3 && urlSplit[3] != "" {
		if urlSplit[3] == "bid" && r.Method == http.MethodPost {
			resp, err = svc.CreateAuctionBid(w, r, id) //	POST 		/auction/:id/bid
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	} else if id == "" && r.Method == http.MethodPost {
		resp, err = svc.Create(w, r) // 								POST 		/auction
	} else if id == "" && r.Method == http.MethodGet {
		resp, err = svc.GetAll(w, r) // 								GET 		/auction
	} else if id != "" && r.Method == http.MethodGet {
		resp, err = svc.GetOne(w, r, id) // 						GET 		/auction/:id
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
