package models

import "testing"
import "github.com/stretchr/testify/assert"
import "fmt"

func Test_NewIllumination(t *testing.T) {
	obj := new(Illumination)

	assert.NotNil(t, obj)
}

func Test_ReadIllumination(t *testing.T) {
	path := "/Users/kazufumiwatanabe/go/src/PixelTool/data/illumination_D65.csv"
	ills := ReadIllumination(path)

	fmt.Println(ills)
}
