package models

/*
ColorSpace :enum for light source
*/
type ColorSpace int

/*ColorSpace*/
const (
	CIE ColorSpace = iota
	NTSC
	SRGB
)
