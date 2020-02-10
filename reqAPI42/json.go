package reqAPI42

import (
	"time"
	"strings"
	"encoding/json"
)

type JSONProjectParent struct {
	API42ProjectParent
}

type JSONTime struct {
	time.Time
}

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
