package flaggio

import (
	"encoding/json"
	"strconv"
)

var _ json.Unmarshaler = (*UserContext)(nil)

// UserContext is a map of strings and any of:
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
			a[k] = strV
		}
	}
	return nil
}
