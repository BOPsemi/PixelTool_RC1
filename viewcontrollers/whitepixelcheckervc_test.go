package viewcontrollers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewVC(t *testing.T) {
	obj := NewWhitePixelCheckerViewController()
	assert.NotNil(t, obj)
}

func Test_CreateWhitePixelPatch(t *testing.T) {
	obj := NewWhitePixelCheckerViewController()

	csvfilepath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/white_pixel.csv"
	filesavepath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"
	dirname := "white_pixel"

	status := obj.CreateWhitePixelPatch(csvfilepath, filesavepath, dirname, 100, 100)

	assert.True(t, status)
}
