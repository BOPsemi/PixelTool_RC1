package models

import (
	"PixelTool_RC1/util"
	"strconv"
)

/*
IlluminationInterface : definition of illumination str interface
*/
type IlluminationInterface interface {
	GetWavelangth() int
	GetIntensity() float64
}

/*
Illumination :illumination structure
*/
type Illumination struct {
	wavelength int
	intensity  float64
}

/*
GetWavelangth :getter
*/
func (il *Illumination) GetWavelangth() int {
	return il.wavelength
}

/*
GetIntensity :getter
*/
func (il *Illumination) GetIntensity() float64 {
	return il.intensity
}

/*
IlluminationMapper : data mapper
*/
func illuminationMapper(data []string) (*Illumination, bool) {
	ill := new(Illumination)
	status := false

	// mapping
	if len(data) > 0 {
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
			mapping
		*/
		ill.wavelength = strToInt(data[0])
		ill.intensity = strToFloat64(data[1])

		// status update
		status = true

	}

	return ill, status
}

/*
ReadIllumination :read Illumination CSV file and map the data to object
*/
func ReadIllumination(path string) []Illumination {
	ills := make([]Illumination, 0)

	// setup csv reader
	reader := util.NewIOUtil()

	// read csv file
	rawdata, status := reader.ReadCSVFile(path)

	// read csv was successful
	if status {
		if len(rawdata) > 0 {
			for _, data := range rawdata {
				ill, mappingstatus := illuminationMapper(data)
				if mappingstatus {
					ills = append(ills, *ill)
				}
			}
		}
	}
	return ills
}
