package models

import "testing"
import "github.com/stretchr/testify/assert"
import "fmt"

func Test_NewWhitePixel(t *testing.T) {
	obj := new(WhitePixel)
	assert.NotNil(t, obj)
}

func Test_ReadWhitePixel(t *testing.T) {
	path := "/Users/kazufumiwatanabe/go/src/PixelTool/data/white_pixel.csv"
	wps := ReadWhitePixel(path)

	fmt.Println(wps)
}

func Test_Getter(t *testing.T) {
	path := "/Users/kazufumiwatanabe/go/src/PixelTool/data/white_pixel.csv"
	wps := ReadWhitePixel(path)

	for _, obj := range wps {
		fmt.Println(obj.GetLevel(), obj.GetCount())
	}
}
