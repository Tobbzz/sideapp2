package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// Utility to open and parse JSON files
func OpenJSON(path string, v interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// Get all servers from file, with options for codesOnly and activeOnly
func GetAllServers(codesOnly, activeOnly bool) ([]string, error) {
	var servers []struct {
		ServerCode string `json:"ServerCode"`
		IsActive   bool   `json:"IsActive"`
	}
	data, err := ioutil.ReadFile("../files/.servers.json")
	if err != nil {
		return nil, err
	}
	var fileData struct {
		Data []struct {
			ServerCode string `json:"ServerCode"`
			IsActive   bool   `json:"IsActive"`
		}
	}
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, err
	}
	result := []string{}
	serverExcludeList := map[string]bool{"cz1":true,"de3":true,"int3":true,"int4":true,"int5":true,"int6":true,"pl2":true,"pl3":true,"pl4":true}
	for _, dd := range fileData.Data {
		if serverExcludeList[dd.ServerCode] {
			continue
		}
		if activeOnly && !dd.IsActive {
			continue
		}
		result = append(result, dd.ServerCode)
	}
	return result, nil
}

// Get all layouts
func GetAllLayouts() ([]map[string]interface{}, error) {
	data, err := ioutil.ReadFile("../layouts/layouts.json")
	if err != nil {
		return nil, err
	}
	var layouts struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &layouts); err != nil {
		return nil, err
	}
	return layouts.Data, nil
}

// GetIncludedStations returns station IDs for a layout number
func GetIncludedStations(layoutNr string) ([]int, error) {
	data, err := ioutil.ReadFile("../layouts/layouts.json")
	if err != nil {
		return nil, err
	}
	var layouts struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &layouts); err != nil {
		return nil, err
	}
	for _, layout := range layouts.Data {
		if strconv.Itoa(int(layout["number"].(float64))) == layoutNr {
			stations := layout["content"].(map[string]interface{})["stations"].([]interface{})
			ids := []int{}
			for _, stn := range stations {
				stnMap := stn.(map[string]interface{})
				ids = append(ids, int(stnMap["pointID"].(float64)))
			}
			return ids, nil
		}
	}
	return nil, nil
}

// Get line info for a list of station IDs
func GetLineInfo(stationsList []int) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile("../lines.json")
	if err != nil {
		return nil, err
	}
	var linesData struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &linesData); err != nil {
		return nil, err
	}
	// For brevity, only return the first line's data
	if len(linesData.Data) == 0 {
		return nil, nil
	}
	return linesData.Data[0], nil
}

// Get signal name logic
func GetSignalName(prefix string, block interface{}, trackDirection, signalDirection string) string {
	blockStr := ""
	switch v := block.(type) {
	case string:
		blockStr = v
	case float64:
		blockStr = strconv.Itoa(int(v))
	}
	if _, err := strconv.Atoi(blockStr); err != nil {
		return blockStr
	}
	if trackDirection == signalDirection {
		return prefix + blockStr
	}
	return prefix + blockStr + "N"
}

// Convert train/timetable data (stub)
func ConvertData(server string, layoutNr string, lineData map[string]interface{}) (map[string]interface{}, error) {
	// Read timetable, stations, trains
	ttData, err := ioutil.ReadFile("../files/" + server + ".timetable.json")
	if err != nil {
		return map[string]interface{}{"error": "Timetable for server " + server + " could not be downloaded."}, nil
	}
	var tt []map[string]interface{}
	if err := json.Unmarshal(ttData, &tt); err != nil {
		return map[string]interface{}{"error": "Could not open timetable file."}, nil
	}
	stationsData, err := ioutil.ReadFile("../files/" + server + ".stations.json")
	if err != nil {
		return map[string]interface{}{"error": "Could not open stations file."}, nil
	}
	var stationsWrap struct{ Data []map[string]interface{} `json:"data"` }
	if err := json.Unmarshal(stationsData, &stationsWrap); err != nil {
		return map[string]interface{}{"error": "Could not open stations file."}, nil
	}
	trainsData, err := ioutil.ReadFile("../files/" + server + ".trains.json")
	if err != nil {
		return map[string]interface{}{"error": "Could not open trains file."}, nil
	}
	var trainsWrap struct{ Data []map[string]interface{} `json:"data"` }
	if err := json.Unmarshal(trainsData, &trainsWrap); err != nil {
		return map[string]interface{}{"error": "Could not open trains file."}, nil
	}
	trains := trainsWrap.Data
	stations := stationsWrap.Data

	includedStations, _ := GetIncludedStations(layoutNr)
	stationsFrLineData := lineData["stations"].([]interface{})
	stationsNames := lineData["stations_names"].([]interface{})
	stationsIds := lineData["stations_ids"].([]interface{})
	entranceSignals := lineData["entrance_signals"].([]interface{})
	signalRules := lineData["signal_rules"].([]interface{})
	remoteStations := lineData["remote_stations"].([]interface{})

	response := map[string]interface{}{
		"error": nil,
		"data": map[string]interface{}{
			"t": []interface{}{},
			"s": []interface{}{},
		},
	}

	// Process stations
	for _, stn := range stationsFrLineData {
		stnMap := stn.(map[string]interface{})
		thisPointID := int(stnMap["id"].(float64))
		dispatchedBy := GetStationUser(stations, stnMap["name"].(string))
		status := "bot"
		if dispatchedBy != "" {
			status = "user"
		}
		for _, rmStn := range remoteStations {
			rmStnMap := rmStn.(map[string]interface{})
			if int(rmStnMap["id"].(float64)) == thisPointID {
				masterStn := int(rmStnMap["controlled_by"].(float64))
				remoteUser := GetStationUser(stations, GetStationById(stationsFrLineData, masterStn))
				if remoteUser != "" {
					status = "remote"
				}
			}
		}
		response["data"].(map[string]interface{})["s"] = append(response["data"].(map[string]interface{})["s"].([]interface{}), map[string]interface{}{
			"Name": stnMap["name"],
			"id": thisPointID,
			"status": status,
			"dispatched_by": dispatchedBy,
		})
	}

	// Process trains and timetables (simplified, full logic can be expanded)
	for _, thisTT := range tt {
		trainNo := thisTT["trainNoLocal"].(string)
		var thisTrain map[string]interface{}
		for _, train := range trains {
			if train["TrainNoLocal"].(string) == trainNo {
				thisTrain = train
				break
			}
		}
		if thisTrain != nil {
			thisTT["trainObject"] = thisTrain
		} else {
			thisTT["trainObject"] = nil
		}
		response["data"].(map[string]interface{})["t"] = append(response["data"].(map[string]interface{})["t"].([]interface{}), thisTT)
	}

	return response, nil
}

// GetStationUser returns the SteamId of the dispatcher for a station
func GetStationUser(stations []map[string]interface{}, targetStationName string) string {
	for _, stn := range stations {
		if stn["Name"].(string) == targetStationName {
			dispatchedBy := stn["DispatchedBy"].([]interface{})
			if len(dispatchedBy) > 0 {
				return dispatchedBy[0].(map[string]interface{})["SteamId"].(string)
			}
			return ""
		}
	}
	return ""
}

// GetStationById returns the station name for a given ID
func GetStationById(stations []interface{}, idTarget int) string {
	for _, stn := range stations {
		stnMap := stn.(map[string]interface{})
		if int(stnMap["id"].(float64)) == idTarget {
			return stnMap["name"].(string)
		}
	}
	return ""
}

// FilterLayoutObjects filters train and station data by included stations
func FilterLayoutObjects(data map[string]interface{}, includedStations []int) map[string]interface{} {
	response := map[string]interface{}{"error": "", "data": map[string]interface{}{"t": []interface{}{}, "s": []interface{}{}}}
	if data["data"] == nil {
		response["error"] = "Server unavailable"
		return response
	}
	tData := data["data"].(map[string]interface{})["t"].([]interface{})
	sData := data["data"].(map[string]interface{})["s"].([]interface{})
	for _, train := range tData {
		trainMap := train.(map[string]interface{})
		if trainMap["trainObject"] != nil {
			obj := trainMap["trainObject"].(map[string]interface{})
			if containsInt(includedStations, int(obj["tpl_next"].(float64))) || containsInt(includedStations, int(obj["tpl_last"].(float64))) {
				response["data"].(map[string]interface{})["t"] = append(response["data"].(map[string]interface{})["t"].([]interface{}), train)
			}
		}
	}
	for _, stn := range sData {
		stnMap := stn.(map[string]interface{})
		if containsInt(includedStations, int(stnMap["id"].(float64))) {
			response["data"].(map[string]interface{})["s"] = append(response["data"].(map[string]interface{})["s"].([]interface{}), stn)
		}
	}
	return response
}

// Helper: containsInt checks if slice contains int
func containsInt(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// Utility: get delay by date diff (minutes)
func GetDelayByDateDiff(dateA, dateB string) int {
	tA, errA := time.Parse("2006-01-02 15:04:05", dateA)
	tB, errB := time.Parse("2006-01-02 15:04:05", dateB)
	if errA != nil || errB != nil {
		return 0
	}
	delta := tB.Sub(tA)
	return int(delta.Minutes())
}

// Utility: signal name matches station
func SignalNameMatchesStn(signalName, stationSign string) bool {
	if stationSign == "" {
		return false
	}
	return strings.HasPrefix(signalName, stationSign)
}

// GetTimeZone returns the timezone name for an offset in hours
func GetTimeZone(offset int) string {
	// Go does not have timezone_name_from_abbr, so fallback to UTC or use offset
	if offset == 0 {
		return "UTC"
	}
	return "Etc/GMT" + strconv.Itoa(-offset) // e.g. Etc/GMT-1 for +1
}

// GetTimeZoneForServer returns the timezone for a server
func GetTimeZoneForServer(server string) string {
	data, err := ioutil.ReadFile("../files/" + server + ".timezones.json")
	if err != nil {
		return "UTC"
	}
	var offset int
	if err := json.Unmarshal(data, &offset); err != nil {
		return "UTC"
	}
	return GetTimeZone(offset)
}

// IfPassedEntrance checks if train passed entrance signal (stub)
func IfPassedEntrance(lineData map[string]interface{}, currentStation, prevStation int) bool {
	lines := lineData["lines"].([]interface{})
	for _, line := range lines {
		lineMap := line.(map[string]interface{})
		if int(lineMap["point_A"].(float64)) == currentStation ||
			int(lineMap["point_B"].(float64)) == currentStation ||
			int(lineMap["point_A"].(float64)) == prevStation ||
			int(lineMap["point_B"].(float64)) == prevStation {
			return true
		}
	}
	return false
}

// SignalName strips @ from signal name
func SignalName(signalName string) string {
	parts := strings.Split(signalName, "@")
	return parts[0]
}

// GetStationPrefixes returns a map of station name to prefix
func GetStationPrefixes(server string) (map[string]string, error) {
	data, err := ioutil.ReadFile("../files/" + server + ".stations.json")
	if err != nil {
		return nil, err
	}
	var stationsWrap struct{ Data []map[string]interface{} `json:"data"` }
	if err := json.Unmarshal(data, &stationsWrap); err != nil {
		return nil, err
	}
	result := map[string]string{}
	for _, stn := range stationsWrap.Data {
		result[stn["Name"].(string)] = stn["Prefix"].(string)
	}
	return result, nil
}

// StringContains checks if haystack contains needle
func StringContains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}

// More functions (train conversion, etc.) can be ported as needed.
