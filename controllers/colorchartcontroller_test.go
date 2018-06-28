package controllers

import (
	"PixelTool_RC1/models"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewColorChartController(t *testing.T) {
	obj := NewColorChartController()

	assert.NotNil(t, obj)
}

func Test_setEnv(t *testing.T) {
	obj := new(colorChartController)

	obj.setEnv(
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/linearmatrix_elm.csv",
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/",
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/device_QE.csv",
		models.D65,
	)

	assert.EqualValues(t, 0, obj.lithSource)
	assert.Equal(t, 6, len(obj.linearMatrixElm))
	assert.NotEqual(t, 0, len(obj.filepath))

	log.Println(obj.linearMatrixElm)
	log.Println(obj.filepath)
}

func mocinit() colorChartController {
	obj := new(colorChartController)

	// set enviroment
	obj.setEnv(
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/linearmatrix_elm.csv",
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/",
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/device_QE.csv",
		models.D65,
	)

	// return object
	return *obj
}

func Test_calculateDeviceResponse(t *testing.T) {
	obj := mocinit()

	assert.True(t, obj.calculateDeviceResponse(0.42))
}

func Test_calculateLinearMatrix(t *testing.T) {
	obj := mocinit()
	elm := []float64{0.136, 0.061, 0.104, 0.063, 0.041, 0.057}

	assert.True(t, obj.calculateDeviceResponse(0.42))
	assert.True(t, obj.calculateLinearMatrix(elm))

	for index, data := range obj.rawLinearizedResponses {
		fmt.Println(index+1, data)
	}
}

func Test_calculateWhiteBalance(t *testing.T) {
	obj := mocinit()
	elm := []float64{0.136, 0.061, 0.104, 0.063, 0.041, 0.057}

	assert.True(t, obj.calculateDeviceResponse(0.42))
	assert.True(t, obj.calculateLinearMatrix(elm))
	assert.True(t, obj.calculateWhiteBalance(22))

	for index, data := range obj.rawWBData {
		fmt.Println(index+1, data)
	}
}

func Test_calculateGammaCorrection(t *testing.T) {
	obj := mocinit()
	elm := []float64{0.636, 0.061, 0.054, 0.063, 0.241, 0.857}

	assert.True(t, obj.calculateDeviceResponse(0.42))
	assert.True(t, obj.calculateLinearMatrix(elm))
	assert.True(t, obj.calculateWhiteBalance(22))
	assert.True(t, obj.calculateGammaCorrection(0.42))

	for index, data := range obj.rawWBData {
		fmt.Println(index+1, data)
	}
}

func Test_convert8BitData(t *testing.T) {
	obj := mocinit()
	elm := []float64{0.136, 0.061, 0.104, 0.063, 0.041, 0.057}

	assert.True(t, obj.calculateDeviceResponse(0.42))
	assert.True(t, obj.calculateLinearMatrix(elm))
	assert.True(t, obj.calculateWhiteBalance(22))
	assert.True(t, obj.calculateGammaCorrection(0.45))
	assert.True(t, obj.convert8BitData())

	for index, data := range obj.devColorCodes {
		fmt.Println(index+1, data)
	}
}

func Test_RunDevice(t *testing.T) {
	obj := NewColorChartController()
	elm := []float64{0.136, 0.061, 0.104, 0.063, 0.041, 0.057}

	results := obj.RunDevice(
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/linearmatrix_elm.csv",
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/",
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/device_QE.csv",
		models.D65,
		0.45,
		elm,
	)

	for _, data := range results {
		fmt.Println(data)
	}

	results = obj.RunDevice(
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/linearmatrix_elm.csv",
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/",
		"/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/device_QE.csv",
		models.D65,
		0.45,
		[]float64{},
	)

	// --- print --
	for _, data := range results {
		fmt.Println(data)
	}
}

func Test_RunStandard(t *testing.T) {
	obj := NewColorChartController()

	path := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/std_color_patch.csv"
	results := obj.RunStandad(path)

	assert.NotEqual(t, 0, len(results))

	// --- print --
	for _, data := range results {
		fmt.Println(data)
	}
}

func Test_GenerateColorPatchImage(t *testing.T) {
	obj := NewColorChartController()

	path := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/std_color_patch.csv"
	savePath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"

	results := obj.RunStandad(path)
	assert.True(t, obj.GenerateColorPatchImage(results, savePath, "std_patch", 100, 100))
	assert.True(t, obj.Generate24MacbethColorChart(false, savePath+"/std_patch/"))
}

func Test_SaveColorPatchCSVData(t *testing.T) {
	obj := NewColorChartController()

	path := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/std_color_patch.csv"
	savePath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/std_patch/"

	results := obj.RunStandad(path)

	assert.True(t, obj.SaveColorPatchCSVData(results, savePath, std24ColorChartName))
}
