package models

import (
	"PixelTool_RC1/util"
	"strconv"
)

/*
ColorCheckerInterface :interface of ColorChecker
*/
type ColorCheckerInterface interface {
	GetWavelength() int
	GetIntensity() float64
}

/*
ColorChecker :color chekcer structure
*/
type ColorChecker struct {
	wavelength int     // wavelength
	intensity  float64 // reflection intensity
}

// GetWavelength :getter
func (cc *ColorChecker) GetWavelength() int {
	return cc.wavelength
}

// GetIntensity :getter
func (cc *ColorChecker) GetIntensity() float64 {
	return cc.intensity
}

/*
ColorCheckerMapper :mapper for colorchecker
*/
func colorCheckerMapper(data []string, order int) (*ColorChecker, bool) {
	// initialize status
	status := false

	// initialize ColorChecker object
	checker := new(ColorChecker)

	if len(data) != 0 && order >= 0 {
		if order < 25 {
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
				strToFloat64 :converter from string to Float64
			*/
			strToFloat64 := func(str string) float64 {
				number, err := strconv.ParseFloat(str, 64)
				if err != nil {
					number = 0.0
				}

				return number
			}

			// mapping
			checker.wavelength = strToInt(data[0])
			checker.intensity = strToFloat64(data[order])

			// update status
			status = true
		}

	}

	return checker, status
}

/*
ReadColorChecker :read color code CSV file and map the data to object
*/
func ReadColorChecker(path string) [][]ColorChecker {
	// initialize buffer
	colorchckers := make([][]ColorChecker, 0)

	// setup csv reader
	reader := util.NewIOUtil()

	// read csv file
	rawdata, status := reader.ReadCSVFile(path)

	// read csv was successful
	if status {
		if len(rawdata) > 0 {
			for i := 1; i < 25; i++ {

				// init checker
				checkers := make([]ColorChecker, 0)

				for _, data := range rawdata {
					checker, mappingstatus := colorCheckerMapper(data, i)
					if mappingstatus {
						checkers = append(checkers, *checker)
					}
				}

				// stack
				colorchckers = append(colorchckers, checkers)
			}
		}
	}

	return colorchckers
}
