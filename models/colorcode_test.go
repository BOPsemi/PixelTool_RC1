package models

import "testing"
import "github.com/stretchr/testify/assert"
import "fmt"

func Test_New(t *testing.T) {
	obj := new(ColorCode)

	assert.NotNil(t, obj)
}

func Test_ReadColorCode(t *testing.T) {
	path := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/macbeth_patch_code.csv"
	colorCodes := ReadColorCode(path)

	assert.EqualValues(t, 24, len(colorCodes))
}

func Test_SetColorCode(t *testing.T) {
	colorcode := SetColorCode(1, "test", 255, 255, 255, 255)
	fmt.Println(colorcode)
}

func Test_SerializeData(t *testing.T) {
	mocdata := SetColorCode(1, "test", 255, 255, 255, 255)

	serializeddata := mocdata.SerializeData()
	fmt.Println(serializeddata)
}
