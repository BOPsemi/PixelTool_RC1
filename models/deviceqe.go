/*
The model of device QE (quantum efficiency)
*/

package models

import (
	"PixelTool_RC1/util"
	"strconv"
)

/*
DeviceQEInterface :interface of Device QE str
*/
type DeviceQEInterface interface {
	GetWavelength() int
	GetGrSignal() float64
	GetGbSignal() float64
	GetRedSignal() float64
	GetBlueSignal() float64
}

/*
DeviceQE :device QE structure
*/
type DeviceQE struct {
	wavelength int
	gr         float64 // Gr channel QE
	gb         float64 // Gb channel QE
	r          float64 // Red channel QE
	b          float64 // Blue channel QE
}

// GetWavelength :getter
func (de *DeviceQE) GetWavelength() int {
	return de.wavelength
}

// GetGrSignal :getter
func (de *DeviceQE) GetGrSignal() float64 {
	return de.gr
}

// GetGbSignal :getter
func (de *DeviceQE) GetGbSignal() float64 {
	return de.gb
}

// GetRedSignal :getter
func (de *DeviceQE) GetRedSignal() float64 {
	return de.r
}

// GetBlueSignal :getter
func (de *DeviceQE) GetBlueSignal() float64 {
	return de.b
}

/*
DeviceQEMapper :data mapper
*/
func deviceQEMapper(data []string) (*DeviceQE, bool) {
	qe := new(DeviceQE)
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
		qe.wavelength = strToInt(data[0])
		qe.gr = strToFloat64(data[1])
		qe.gb = strToFloat64(data[2])
		qe.r = strToFloat64(data[3])
		qe.b = strToFloat64(data[4])

		// status update
		status = true

	}

	return qe, status
}

/*
ReadDeviceQE :read device QE CSV file and map the data to object
*/
func ReadDeviceQE(path string) []DeviceQE {
	qes := make([]DeviceQE, 0)

	// setup csv reader
	reader := util.NewIOUtil()

	// read csv file
	rawdata, status := reader.ReadCSVFile(path)

	// read csv was successful
	if status {
		if len(rawdata) > 0 {
			for _, data := range rawdata {
				qe, mappingstatus := deviceQEMapper(data)
				if mappingstatus {
					qes = append(qes, *qe)
				}
			}
		}
	}
	return qes
}
