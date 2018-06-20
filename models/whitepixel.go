/*
Definition of white pixel
*/

package models

import (
	"PixelTool_RC1/util"
	"strconv"
)

/*
WhitePixelInterface :interface of WhitePixel model
*/
type WhitePixelInterface interface {
	GetLevel() int
	GetCount() int
}

/*
WhitePixel :white pixel structure
*/
type WhitePixel struct {
	level int // level, unit is DN
	count int // count, not ppm
}

/*
GetLevel :getter
*/
func (wh *WhitePixel) GetLevel() int {
	return wh.level
}

/*
GetCount :getter
*/
func (wh *WhitePixel) GetCount() int {
	return wh.count
}

/*
WhitePixelMapper : white pixel mapper
*/
func whitePixelMapper(data []string) (*WhitePixel, bool) {
	wp := new(WhitePixel)
	status := false

	// mapping
	if len(data) > 0 {
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
			mapping to wp structure
		*/
		wp.level = strToInt(data[0])
		wp.count = strToInt(data[1])

		// status update
		status = true

	}

	return wp, status
}

/*
ReadWhitePixel :read white pixel CSV file and map the data to object
*/
func ReadWhitePixel(path string) []WhitePixel {
	// initialize buffer
	wps := make([]WhitePixel, 0)

	// setup csv reader
	reader := util.NewIOUtil()

	// read csv file
	rawdata, status := reader.ReadCSVFile(path)

	// read csv was successful
	if status {
		if len(rawdata) > 0 {
			for _, data := range rawdata {
				wp, mappingstatus := whitePixelMapper(data)
				if mappingstatus {
					wps = append(wps, *wp)
				}
			}
		}
	}
	return wps
}

/*
SetWhitePixel :
	in 	;level, count
	out	;white pixel
*/
func SetWhitePixel(level, count int) *WhitePixel {
	return &WhitePixel{level: level, count: count}
}
