package flaggio

import (
	"encoding/json"
	"strconv"
)

var _ json.Unmarshaler = (*UserContext)(nil)

// UserContext is a map of strings and one of:
// int64, float64, bool, string
type UserContext map[string]interface{}

// UnmarshalJSON unmarshals the bytes into UserContext
func (a UserContext) UnmarshalJSON(b []byte) error {
	data := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	for k, v := range data {
		strV := string(v)
		if n, err := strconv.ParseInt(strV, 10, 64); err == nil {
			a[k] = n
		} else if n, err := strconv.ParseFloat(strV, 64); err == nil {
			a[k] = n
		} else if b, err := strconv.ParseBool(strV); err == nil {
			a[k] = b
		} else {
			// everything else is treated as a string, even null
			// strings will still be quoted on the json.RawMessage so we
			// try to unmarshal them. it will fail for objects and arrays
			// in that case, ignore the error and return the raw string
			_ = json.Unmarshal(v, &strV)
			a[k] = strV
		}
	}
	return nil
}
