package compare

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInclude(t *testing.T) {
	for _, item := range includeTable {
		comparator := NewComparator(DefaultParser{}, item.option...)
		result := comparator.Equal(item.source, item.target)
		assert.Equal(t, item.except, result)
	}
}

var simpleMap = map[string]interface{}{
	"key": "value",
}

var includeTable = []struct {
	source interface{}
	target interface{}
	option []Option
	except bool
}{
	{simpleMap, simpleMap, []Option{IncludeField("key")}, true},
}
