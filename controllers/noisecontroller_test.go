package controllers

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewNoiseController(t *testing.T) {
	obj := NewNoiseController()
	assert.NotNil(t, obj)
}

func Test_AddWhitePixelNoise(t *testing.T) {
	// base moc
	baseMoc := &color.RGBA{
		R: 10,
		G: 20,
		B: 30,
		A: 255,
	}

	base := []color.RGBA{*baseMoc, *baseMoc}

	// noise moc
	noiseMoc := &color.RGBA{
		R: 55,
		G: 65,
		B: 75,
		A: 255,
	}
	noise := []color.RGBA{*noiseMoc, *noiseMoc}

	// test
	obj := NewNoiseController()
	result := obj.AddWhitePixelNoise(base, noise, 50)

	// check
	assert.Equal(t, 2, len(result))
	fmt.Println(result[0].G, result[0].R)

}
