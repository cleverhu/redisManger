package textUtils

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestPullMiddleText(t *testing.T) {
	str := "1234556"
	assert.Equal(t, "45", PullMiddleText(str, "123", "56"))
}

