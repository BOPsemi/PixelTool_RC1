package models

import "testing"

import "github.com/stretchr/testify/assert"
import "fmt"

func Test_ReadWhitePoint(t *testing.T) {
	path := "/Users/kazufumiwatanabe/go/src/PixelTool/json/whitepoint.json"
	if obj, ok := ReadWhitePoint(path); ok {
		assert.True(t, ok)
		fmt.Println(obj)
	}
}
