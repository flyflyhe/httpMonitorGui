package stringHelper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitStringByLength(t *testing.T) {
	str := "中国abcdefg"
	newStr := "中国;ab;cd;ef;g"

	result := SplitStringByLength(str, ";", 2)

	assert.Equal(t, newStr, result)
}
