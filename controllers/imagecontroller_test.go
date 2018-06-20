package controllers

import (
	"PixelTool_RC1/models"
	"PixelTool_RC1/util"
	"fmt"
	"image/color"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewImageController(t *testing.T) {
	obj := NewImageController()

	assert.NotNil(t, obj)
}

func Test_CreateImage(t *testing.T) {
	obj := NewImageController()

	// check initializer
	assert.NotNil(t, obj)

	// load CSV file and create color code obj slices
	path := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/Macbeth_Patch_Code.csv"
	colorCodes := models.ReadColorCode(path)

	// generate rgba data
	rgbadata := make([]color.RGBA, 0)

	for _, data := range colorCodes {
		rgba := data.GenerateColorRGBA()
		rgbadata = append(rgbadata, *rgba)
	}

	// check data size
	assert.EqualValues(t, 24, len(rgbadata))

	// create sample image
	sampleImageData := rgbadata[0]
	//fmt.Println(sampleImageData)

	// create solid image
	rawImage := obj.CreateSolidImage(sampleImageData, 100, 100)
	assert.NotNil(t, rawImage)

	fmt.Println(reflect.TypeOf(rawImage).String())

	// save PNG file
	imagepath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"

	iohandler := util.NewIOUtil()
	iohandler.StreamOutPNGFile(imagepath, colorCodes[0].GetName(), rawImage)

}

func Test_Create24MacbethChart(t *testing.T) {
	path := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/std_patch/"
	fileName := "std_macbeth_chart"

	obj := NewImageController()
	assert.True(t, obj.Create24MacbethChart(path, fileName))

}
