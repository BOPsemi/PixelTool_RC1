package controllers

import (
	"PixelTool_RC1/util"
	"strconv"
)

/*
DeltaEvalController :delta-E evaluation controller
*/
type DeltaEvaluationController interface {
	EvaluateDeltaE(refDataPath, compDataPath string, kvalues []float64) ([]float64, bool)
	SaveDeltaEResultData(savepath, filename string) bool
}

// structure
type deltaEvaluationController struct {
	deltaEResults []float64
}

/*
NewDeltaEvaluationController : initializer
*/
func NewDeltaEvaluationController() DeltaEvaluationController {
	obj := new(deltaEvaluationController)

	obj.deltaEResults = make([]float64, 0)

	return obj
}

/*
EvaluateDeltaE
	in	:refDataPath, compDataPath string, kvalues []float64
	out	:bool
*/
func (dc *deltaEvaluationController) EvaluateDeltaE(refDataPath, compDataPath string, kvalues []float64) ([]float64, bool) {
	status := false

	// io handler
	ioHandler := util.NewIOUtil()

	// data buffer
	refData := make([][]float64, 0)
	devData := make([][]float64, 0)
	deltaEResults := make([]float64, 0)

	// data array extractor
	makeDataArray := func(rawdata [][]string) [][]float64 {
		// buffer
		var convertedData [][]float64

		for _, patch := range rawdata {
			// buffer data
			eachPatchData := make([]float64, 0)

			// extract data and convert data
			for index, data := range patch {
				if index > 1 && index < 5 {
					value, err := strconv.ParseFloat(data, 64)
					if err == nil {
						eachPatchData = append(eachPatchData, value)
					}
				}
			}
			// update each patch data to ref data
			convertedData = append(convertedData, eachPatchData)
		}
		// return the converted data
		return convertedData
	}

	// open reference csv file
	if refPatchRawData, ok := ioHandler.ReadCSVFile(refDataPath); ok {
		refData = makeDataArray(refPatchRawData)
	} else {
		return deltaEResults, status
	}

	// open dev csv file
	if devPatchRawData, ok := ioHandler.ReadCSVFile(compDataPath); ok {
		devData = makeDataArray(devPatchRawData)
	} else {
		return deltaEResults, status
	}

	/*
		calculate deltaE
		data[0]	:Red
		data[1]	:Green
		data[2]	:Blue
	*/
	deltaEcalculator := util.NewDeltaLabCalculator()
	for index := 0; index < 24; index++ {
		deltaE := deltaEcalculator.DeltaLab(refData[index], devData[index], kvalues)
		deltaEResults = append(deltaEResults, deltaE)
	}

	if len(deltaEResults) == 0 {
		return deltaEResults, status
	}

	// update status
	status = true
	dc.deltaEResults = deltaEResults

	return deltaEResults, status
}

/*
SaveDeltaEResultData
	in	:savepath, filename string
	out	:bool
*/
func (dc *deltaEvaluationController) SaveDeltaEResultData(savepath, filename string) bool {
	status := false

	if len(dc.deltaEResults) == 24 {
		if savepath != "" && filename != "" {

			// make string
			data := make([][]string, 0)
			for index, obj := range dc.deltaEResults {
				str := []string{strconv.Itoa(index + 1), strconv.FormatFloat(obj, 'f', 4, 64)}
				data = append(data, str)
			}

			// save data
			iohandler := util.NewIOUtil()
			if iohandler.WriteCSVFile(savepath, filename, data) {
				// status update
				status = true
			}
		}
	}

	return status
}
