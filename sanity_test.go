package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonUnmarshalNumberToString(t *testing.T) {
	// Unmarshalling a JSON number into a struct field of type string fails
	var outA struct{ X string }
	err := json.Unmarshal([]byte(`{"X": 123}`), &outA)
	assert.Error(t, err)
	assert.Equal(t, "", outA.X)

	// But unmarshalling into json.Number works
	var outB struct{ X json.Number }
	err = json.Unmarshal([]byte(`{"X": 123}`), &outB)
	assert.NoError(t, err)
	assert.Equal(t, "123", outB.X.String())

	// And unmarshalling a string into json.Number works
	err = json.Unmarshal([]byte(`{"X": "456"}`), &outB)
	assert.NoError(t, err)
	assert.Equal(t, "456", outB.X.String())

}

func TestJsonUnmarshalInt64(t *testing.T) {
	// Unmarshalling a JSON number into a struct field of type string fails
	var outA struct{ X int64 }
	err := json.Unmarshal([]byte(`{"X": 123}`), &outA)
	assert.NoError(t, err)
	assert.Equal(t, int64(123), outA.X)
}
