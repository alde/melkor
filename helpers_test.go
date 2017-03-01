package melkor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ModifyTags(t *testing.T) {
	var input []interface{}
	input = append(input, map[string]string{"Key": "egg", "Value": "bacon"})
	input = append(input, map[string]string{"Key": "bob", "Value": "hope"})

	var actual []interface{}
	copy(actual, input)
	ModifyTags(actual)

	for _, t0 := range actual {
		tag := t0.(map[string]string)

		key := tag["Key"]
		value := tag["Value"]

		assert.Equal(t, tag[key], value)
	}
}
