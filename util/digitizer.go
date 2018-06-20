package util

/*
Digitizer :degitize the data
*/
type Digitizer interface {
	D8bitDigitizeData(data float64, refLevel uint8) uint8
	D10bitDigitizeData(data float64, refLevel uint16) uint16
}

// structure definition
type digitizer struct {
}

/*
NewDigitizer :initializer of Digitizer
*/
func NewDigitizer() Digitizer {
	obj := new(digitizer)

	return obj
}

/*
DigitizeDataAt8bitBy
	in	;data float64, refLevel uint8
	out	;uint8
*/
func (di *digitizer) D8bitDigitizeData(data float64, refLevel uint8) uint8 {
	var digital uint8
	if refLevel == 0 {
		// no ref level, it means 1.0 -> 255
		buff := uint8(data * 255.0)
		if buff > 255 {
			digital = 255
		} else {
			digital = buff
		}
	} else {
		// with ref level, it means 1.0 -> refLevel
		buff := uint8(data * float64(refLevel))
		if buff > 255 {
			digital = 255
		} else {
			digital = buff
		}
	}

	return digital
}

/*
D10bitDigitizeData(
	in	;data float64, refLevel uint16
	out	;uint16
*/
func (di *digitizer) D10bitDigitizeData(data float64, refLevel uint16) uint16 {
	var digital uint16
	if refLevel == 0 {
		// no ref level, it means 1.0 -> 1023
		digital = uint16(data) * 1023
	} else {
		// with ref level, it means 1.0 -> refLevel
		digital = uint16(data) * refLevel
	}

	return digital
}
