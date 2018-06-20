package viewcontrollers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewColorCheckerViewController(t *testing.T) {
	obj := NewColorCheckerViewController()

	assert.NotNil(t, obj)
}

func Test_ColorCheckerViewController(t *testing.T) {
	obj := NewColorCheckerViewController()

	csvfilepath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/Macbeth_Patch_Code.csv"
	filesavepath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"
	dirname := "std_patch"

	obj.CreateColorCodePatch(csvfilepath, filesavepath, dirname, 100, 100)
}

func Test_SaveColorCodePatchData(t *testing.T) {
	obj := NewColorCheckerViewController()

	csvfilepath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/Macbeth_Patch_Code.csv"
	filesavepath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"
	dirname := "std_patch"

	status := obj.CreateColorCodePatch(csvfilepath, filesavepath, dirname, 100, 100)
	assert.True(t, status)

	filepath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"
	filename := "std_color_patch"

	status = obj.SaveColorCodePatchData(filepath, filename)
	assert.True(t, status)

}
