package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONB used to store map[string]any to postgres jsonb column
type JSONB map[string]any

func (j *JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	encodedJSONCategories, err := json.Marshal(*j)
	if err != nil {
		return nil, err
	}
	return string(encodedJSONCategories), nil
}

func (j *JSONB) Scan(src any) error {
	mapData := map[string]any{}
	switch src := src.(type) {
	case string:
		err := json.Unmarshal([]byte(src), &mapData)
		if err != nil {
			return err
		}
	case []byte:
		if err := json.Unmarshal(src, &mapData); err != nil {
			return err
		}
	case nil:
		*j = nil
		return nil
	default:
		*j = nil
		return errors.New("unknown jsonb map")
	}
	*j = mapData
	return nil
}

func (j *JSONB) SetMap(jsonbMapData map[string]any) error {
	dataVal, err := j.Value()
	if err != nil {
		return err
	}
	castData, ok := dataVal.(string)
	if !ok {
		return errors.New("error casting json data to string")
	}
	return json.Unmarshal([]byte(castData), &jsonbMapData)
}
