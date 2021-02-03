package flaggio_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uw-labs/flaggio/internal/flaggio"
)

func TestUserContext_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	ucJSON := []byte(`
	{
		"string": "value",
		"int": 1,
		"float": 2.5,
		"bool": true,
		"null": null,
		"object": {},
		"array": []
	}`)
	uc := make(flaggio.UserContext)
	err := json.Unmarshal(ucJSON, &uc)
	assert.NoError(t, err)
	assert.Equal(t, "value", uc["string"])
	assert.Equal(t, int64(1), uc["int"])
	assert.Equal(t, float64(2.5), uc["float"])
	assert.Equal(t, true, uc["bool"])
	assert.Equal(t, "null", uc["null"])
	assert.Equal(t, "{}", uc["object"])
	assert.Equal(t, "[]", uc["array"])
}
