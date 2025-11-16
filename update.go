package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// UpdateHandler triggers backend update logic
func updateHandler(w http.ResponseWriter, r *http.Request) {
	task := r.URL.Query().Get("task")
	if task == "1" {
		if err := updateBackend(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("Update complete"))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Invalid task"))
}

// updateBackend runs the update logic for all servers
func updateBackend() error {
	servers, err := getAllServers()
	if err != nil {
		return err
	}
	for _, server := range servers {
		if err := fetchData("time", server.Code); err != nil {
			log.Printf("Error updating time for %s: %v", server.Code, err)
		}
		if err := fetchData("trains", server.Code); err != nil {
			log.Printf("Error updating trains for %s: %v", server.Code, err)
		}
		if err := fetchData("stations", server.Code); err != nil {
			log.Printf("Error updating stations for %s: %v", server.Code, err)
		}
		// Add delay calculation and other logic as needed
	}
	return nil
}

// fetchData fetches and saves data for a given type and server
func fetchData(dataType, server string) error {
	var url string
	switch dataType {
	case "servers":
		url = "https://panel.simrail.eu:8084/servers-open"
	case "time":
		url = "https://panel.simrail.eu:8084/time-open?serverCode=" + server
	case "trains":
		url = "https://panel.simrail.eu:8084/trains-open?serverCode=" + server
	case "stations":
		url = "https://panel.simrail.eu:8084/stations-open?serverCode=" + server
	default:
		return nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	filePath := "../files/" + server + "." + dataType + ".json"
	return ioutil.WriteFile(filePath, body, 0644)
}

// Server struct for server list
type Server struct {
	Code string
	Name string
}

// getAllServers loads server list from file or API
func getAllServers() ([]Server, error) {
	// For now, load from servers-open API
	resp, err := http.Get("https://panel.simrail.eu:8084/servers-open")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data struct {
		Data []struct {
			ServerCode string `json:"ServerCode"`
			ServerName string `json:"ServerName"`
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	servers := make([]Server, 0, len(data.Data))
	for _, s := range data.Data {
		servers = append(servers, Server{Code: s.ServerCode, Name: s.ServerName})
	}
	return servers, nil
}
