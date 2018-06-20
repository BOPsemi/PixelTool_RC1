package viewcontrollers

import "testing"
import "github.com/stretchr/testify/assert"

const (
	BASE  = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/std_patch/Blue.png"
	NOISE = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/white_pixel/white_pixel.png"

	filesavepath = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"
	dirname      = "add_wp_noise"
)

func Test_NewNoiseAdditionViewController(t *testing.T) {
	obj := NewNoiseAdditionViewController()

	assert.NotNil(t, obj)
}

func Test_SetImageDataForWhitePixelAddition(t *testing.T) {
	obj := NewNoiseAdditionViewController()

	assert.True(t, obj.SetImageDataForWhitePixelAddition(BASE, NOISE))
	//assert.False(t, obj.SetImageDataForWhitePixelAddition("", NOISE))
	//assert.False(t, obj.SetImageDataForWhitePixelAddition(BASE, ""))
}

func Test_CreateImageWithWhitePixel(t *testing.T) {
	obj := NewNoiseAdditionViewController()

	status := obj.SetImageDataForWhitePixelAddition(BASE, NOISE)
	assert.True(t, status)

	obj.CreateImageWithWhitePixel(25, "test", filesavepath, dirname)

}
