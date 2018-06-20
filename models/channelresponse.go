package models

/*
ChannelResponse :define channel response
*/
type ChannelResponse struct {
	CheckerNumber int     // color checker number
	Gr            float64 // Gr channel
	Gb            float64 // Gb channel
	R             float64 // blue channel
	B             float64 // red channel
}
