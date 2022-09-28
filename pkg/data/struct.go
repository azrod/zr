package data

import (
	"encoding/json"

	"github.com/azrod/zr/pkg/format"
)

type DataHotReload struct {
	LogLevel  string           `json:"log_level"`
	LogFormat format.LogFormat `json:"log_format"`
}

type DataHotReloadJSON []byte

// String returns the JSON representation of the object.
func (d DataHotReloadJSON) String() string {
	return string(d)
}

// JSONMarshal marshals the object into JSON.
func (d DataHotReload) JSONMarshal() DataHotReloadJSON {
	b, _ := json.Marshal(d)
	return DataHotReloadJSON(b)
}

// // UnmarshalJSON unmarshals the JSON into the object.
// func (d *DataHotReload) UnmarshalJSON(b []byte) error {
// 	// var x *DataHotReload
// 	err := json.Unmarshal(b, &d)
// 	// d = x
// 	return err
// }
