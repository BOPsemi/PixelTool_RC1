package models

import (
	"encoding/json"
	"io/ioutil"
)

/*
WhitePointInterface :interface of white point
*/
type WhitePointInterface interface {
	ReadWhitePoint(path string) (*WhitePoint, bool)
}

/*
WhitePointRGBPrime : definition of prime in RGB white point
*/
type WhitePointRGBPrime struct {
	X float64 `json: "x"`
	Y float64 `json: "y"`
	Z float64 `json: "z"`
}

/*
WhitePointWitePrime : definition of prime in white point
*/
type WhitePointWitePrime struct {
	X2 float64 `json: "x2"`
	Y2 float64 `json: "y2"`
	Z2 float64
}

/*
WhitePointRGBW : definition of RGBW structure
*/
type WhitePointRGBW struct {
	R  WhitePointRGBPrime  `json: "r"`
	G  WhitePointRGBPrime  `json: "g"`
	B  WhitePointRGBPrime  `json: "b"`
	W  WhitePointWitePrime `json: "w"`
	Wn struct {
		X float64
		Y float64
		Z float64
	}
	Wns []float64
}

/*
WhitePoint : Defitnion of SRGB white point
*/
type WhitePoint struct {
	CIE  WhitePointRGBW `json: "cie"`
	NTSC WhitePointRGBW `json: "ntsc"`
	SRGB WhitePointRGBW `json: "srgb"`
}

/*
ReadWhitePoint :initialize of white point str
*/
func ReadWhitePoint(path string) (*WhitePoint, bool) {
	obj := new(WhitePoint)
	status := false

	if path != "" {
		rawdata, err := ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(rawdata, &obj)
			if err == nil {
				status = true
			}
		}
	}

	return obj, status
}
