package main

import (
	"encoding/json"
	"io/ioutil"
)

// Utility to open and parse JSON files
func OpenJSON(path string, v interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
