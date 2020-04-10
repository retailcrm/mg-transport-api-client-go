// +build go1.14

package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func OriginatorMarshalJSONInvalid(t *testing.T) {
	data, err := json.Marshal(Originator(3))
	assert.Nil(t, data)
	assert.EqualError(t, err, "json: error calling MarshalText for type v1.Originator: invalid originator")
}
