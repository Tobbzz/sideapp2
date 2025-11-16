package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// apiHandler serves layout and train data
func apiHandler(w http.ResponseWriter, r *http.Request) {
	typeParam := r.URL.Query().Get("type")
	server := r.URL.Query().Get("server")
	layoutNr := r.URL.Query().Get("layout")

	if typeParam == "layout" {
		serveLayout(w, layoutNr)
		return
	}
	if typeParam == "trains" {
		serveTrains(w, server, layoutNr)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Invalid type"))
}

func serveLayout(w http.ResponseWriter, layoutNr string) {
	data, err := ioutil.ReadFile("../layouts/layouts.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Layout file error"))
		return
	}
	var layouts struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &layouts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Layout JSON error"))
		return
	}
	for _, layout := range layouts.Data {
		if layout["number"] == layoutNr {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(layout)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Layout not found"))
}

func serveTrains(w http.ResponseWriter, server, layoutNr string) {
	filePath := "../files/" + server + ".z_readydata.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Train data not found"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
