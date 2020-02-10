package reqAPI42

import (
	"encoding/json"
	"strings"
	"time"
)

// JSONProjectParent is used to decode JSON values from API 42 for project parent
type JSONProjectParent struct {
	API42ProjectParent
}

// JSONTime is used to decode JSON values from API 42 for time
type JSONTime struct {
	time.Time
}

// UnmarshalJSON is used to deserialize JSONTime
func (jsonVal *JSONTime) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str == "null" {
		*jsonVal = JSONTime{time.Time{}}
		return nil
	}
	timeFormated := strings.Trim(str, "\"")
	timeVal, err := time.Parse(time.RFC3339, timeFormated)
	if err != nil {
		return err
	}
	*jsonVal = JSONTime{timeVal}
	return nil
}

// UnmarshalJSON is used to deserialize JSONProjectParent
func (jsonVal *JSONProjectParent) UnmarshalJSON(b []byte) error {
	str := string(b)

	if str == "null" {
		jsonVal = nil
		return nil
	}

	var projectParent API42ProjectParent
	err := json.Unmarshal(b, &projectParent)
	if err != nil {
		return err
	}
	jsonVal = &JSONProjectParent{projectParent}
	return nil
}
