package controllers

import (
	"image/color"
)

/*
NoiseController :interface of noise controller
*/
type NoiseController interface {
	AddWhitePixelNoise(base, noise []color.RGBA, darklevel int) []color.RGBA
}

// definition of structure
type noiseController struct {
}

/*
NewNoiseController :initializer of noise controller
*/
func NewNoiseController() NoiseController {
	obj := new(noiseController)

	return obj
}

/*
AddWhitePixelNoise :
	in	;base, noise []color.RGBA
	out	;[]color.RGBA
*/
func (nc *noiseController) AddWhitePixelNoise(base, noise []color.RGBA, darklevel int) []color.RGBA {
	// buffer for saving the results
	result := make([]color.RGBA, 0)

	if len(base) == len(noise) {
		// the size matched case
		/*
			Step-1	:dark level shift
			Step-2	:add white pixel to base data
		*/

		// Step-1	;dark level shift
		levelShifttedNoise := make([]color.RGBA, 0)
		for _, data := range noise {
			shifttedNoise := nc.darkLevelShift(data, darklevel)
			if shifttedNoise != nil {
				levelShifttedNoise = append(levelShifttedNoise, *shifttedNoise)
			}
		}

		// Step-2	;white noise addition
		noiseAddedpixels := make([]color.RGBA, 0)
		for index := 0; index < len(base); index++ {
			pixel := nc.addWhitePixel(base[index], levelShifttedNoise[index])
			if pixel != nil {
				noiseAddedpixels = append(noiseAddedpixels, *pixel)
			}
		}

		// Step-3	;check size
		if len(noiseAddedpixels) == len(base) {
			result = noiseAddedpixels
		}
	}

	return result
}

/*
Add white pixel
*/
func (nc *noiseController) addWhitePixel(base, wp color.RGBA) *color.RGBA {
	return &color.RGBA{
		R: base.R + wp.R,
		G: base.G + wp.G,
		B: base.B + wp.B,
		A: 255,
	}
}

/*
Darklevel shifter
*/
func (nc *noiseController) darkLevelShift(noise color.RGBA, darklevel int) *color.RGBA {
	rLevel := int(noise.R) - darklevel
	gLevel := int(noise.G) - darklevel
	bLevel := int(noise.B) - darklevel

	if rLevel*gLevel*bLevel < 0 {
		return &color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}
	return &color.RGBA{R: uint8(rLevel), G: uint8(gLevel), B: uint8(bLevel), A: 255}
}
