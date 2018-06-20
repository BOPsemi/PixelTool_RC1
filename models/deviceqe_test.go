package models

import "testing"
import "github.com/stretchr/testify/assert"
import "fmt"

func Test_NewDeviceQE(t *testing.T) {
	obj := new(DeviceQE)

	assert.NotNil(t, obj)
}

func Test_ReadDeviceQE(t *testing.T) {
	path := "/Users/kazufumiwatanabe/go/src/PixelTool/data/device_QE.csv"
	deviceQEs := ReadDeviceQE(path)

	fmt.Println(deviceQEs)
}
