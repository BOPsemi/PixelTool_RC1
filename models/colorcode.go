/*
Definition of color code
Basically, this structure handles Macbeth color code
*/

package models

import (
	"PixelTool_RC1/util"
	"image/color"
	"strconv"
)

/*
ColorCodeInterface :interface definition
*/
type ColorCodeInterface interface {
	GetNumber() int
	GetName() string
	GetGreenSignal() uint8
	GetRedSignal() uint8
	GetBlueSignal() uint8
	GetASignal() uint8
	GenerateColorRGBA() *color.RGBA

	SerializeData() []string
}

/*
ColorCode :Macbeth Color Code structure
*/
type ColorCode struct {
	number int    // code number
	name   string // code name
	r      uint8  // red, should be unit8
	g      uint8  // green, should be unit8
	b      uint8  // blue, should be unit8
	a      uint8  // a, should be unit8
}

// GetNumber :getter
func (cl *ColorCode) GetNumber() int {
	return cl.number
}

// GetName :getter
func (cl *ColorCode) GetName() string {
	return cl.name
}

// GetGreenSignal :getter
func (cl *ColorCode) GetGreenSignal() uint8 {
	return cl.g
}

// GetRedSignal :getter
func (cl *ColorCode) GetRedSignal() uint8 {
	return cl.r
}

// GetBlueSignal :getter
func (cl *ColorCode) GetBlueSignal() uint8 {
	return cl.b
}

// GetASignal :getter
func (cl *ColorCode) GetASignal() uint8 {
	return cl.a
}

/*
ColorCodeMapper :mapper for ColorCode
*/
func colorCodeMapper(data []string) (*ColorCode, bool) {
	// initialize status
	status := false

	// initialize ColorCode Object
	code := new(ColorCode)

	// mapping
	if len(data) != 0 {

		/*
			strToInt :convert string to Int
			If the error was detected, the function return -1
		*/
		strToInt := func(str string) int {
			number, err := strconv.Atoi(str)
			if err != nil {
				number = -1
			}
			return number
		}

		/*
			strToUint8 :convert string to uint8
			If the error was detected, the function return 0
		*/
		strToUint8 := func(str string) uint8 {
			number, err := strconv.Atoi(str)
			if err != nil {
				number = 0
			}

			return uint8(number)
		}

		/*
			Mapping
		*/
		code.number = strToInt(data[0])
		code.name = data[1]
		code.r = strToUint8(data[2])
		code.g = strToUint8(data[3])
		code.b = strToUint8(data[4])
		code.a = 255

		// update status
		status = true
	}

	return code, status
}

/*
ReadColorCode :read color code CSV file and map the data to object
*/
func ReadColorCode(path string) []ColorCode {
	// initialize buffer
	colorcodes := make([]ColorCode, 0)

	// setup csv reader
	reader := util.NewIOUtil()

	// read csv file
	rawdata, status := reader.ReadCSVFile(path)

	// read csv was successful
	if status {
		if len(rawdata) > 0 {
			for _, data := range rawdata {
				colorcode, mappingstatus := colorCodeMapper(data)
				if mappingstatus {
					colorcodes = append(colorcodes, *colorcode)
				}
			}
		}
	}

	return colorcodes
}

/*
SetColorCode :set color code
	in	;patchNumber int, patchName string, rin, gin, bin, ain uint8
	out	;*ColorCode
*/
func SetColorCode(patchNumber int, patchName string, rin, gin, bin, ain uint8) *ColorCode {
	colorcode := &ColorCode{
		number: patchNumber,
		name:   patchName,
		r:      rin,
		g:      gin,
		b:      bin,
		a:      ain,
	}

	return colorcode
}

/*
GenerateColorRGBA :generate color.RGBA str from self data
*/
func (cl *ColorCode) GenerateColorRGBA() *color.RGBA {
	return &color.RGBA{
		R: cl.r,
		G: cl.g,
		B: cl.b,
		A: cl.a,
	}
}

/*
SerializeData : serialize the data
	in	;
	out	;[]string
*/
func (cl *ColorCode) SerializeData() []string {

	return []string{
		strconv.Itoa(cl.number),
		cl.name,
		strconv.Itoa(int(cl.r)),
		strconv.Itoa(int(cl.g)),
		strconv.Itoa(int(cl.b)),
		strconv.Itoa(int(cl.a)),
	}
}
