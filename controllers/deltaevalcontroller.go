package controllers

import (
	"PixelTool_RC1/models"
	"PixelTool_RC1/tools"
	"PixelTool_RC1/util"
	"strconv"
)

/*
DeltaEvaluationController :delta-E evaluation controller
*/
type DeltaEvaluationController interface {
	SetData(refDataPath, compDataPath string) bool
	RGBtoLabConversion(colorSpace models.ColorSpace, rgb []float64) []float64
	CalculateDeltaE(ref, comp [][]float64, kvalues []float64) ([]float64, float64)
	RunDeltaEEvaluation(colorSpace models.ColorSpace, refData, compData [][]float64, kvalues []float64) ([]float64, float64)

	SaveDeltaEResultData(savepath, filename string) bool
}

// structure
type deltaEvaluationController struct {
	deltaEResults []float64

	refData  [][]float64 // reference data
	compData [][]float64 // comp data
}

/*
NewDeltaEvaluationController : initializer
*/
func NewDeltaEvaluationController() DeltaEvaluationController {
	obj := new(deltaEvaluationController)

	obj.deltaEResults = make([]float64, 0)
	obj.refData = make([][]float64, 0)
	obj.compData = make([][]float64, 0)

	return obj
}

func (dc *deltaEvaluationController) makeDataArray(rawdata [][]string) [][]float64 {
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

/*
SetData
	in	:refDataPath, compDataPath string
	out	:bool
*/
func (dc *deltaEvaluationController) SetData(refDataPath, compDataPath string) bool {
	// check data path
	if refDataPath == "" || compDataPath == "" {
		return false
	}

	// read csv file
	iohandler := util.NewIOUtil()
	if refData, ok := iohandler.ReadCSVFile(refDataPath); ok {
		dc.refData = dc.makeDataArray(refData)
	} else {
		return false
	}
	if compData, ok := iohandler.ReadCSVFile(compDataPath); ok {
		dc.compData = dc.makeDataArray(compData)
	} else {
		return false
	}

	// return
	return true
}

/*
RGBtoLabConversion
	in	:
*/
func (dc *deltaEvaluationController) RGBtoLabConversion(colorSpace models.ColorSpace, rgb []float64) []float64 {

	spaceChanger := NewColorSpaceChange()

	// whitpoint file path
	dirhandler := util.NewDirectoryHandler()
	path := dirhandler.GetCurrentDirectoryPath() + "/json/whitepoint.json"

	// set white point
	if !spaceChanger.ReadWhitePoint(path) {
		return []float64{}
	}

	// calculate
	xyz := spaceChanger.SpaceChangeRGBtoXYZ(colorSpace, rgb)
	lab := spaceChanger.SpaceChangeXYZtoLab(colorSpace, xyz)

	return lab
}

/*
RunDeltaEEvaluation
	in	:refData, compData [][]float64, kvalues []float64
	out	:[]float64, float64
*/
func (dc *deltaEvaluationController) RunDeltaEEvaluation(colorSpace models.ColorSpace, refData, compData [][]float64, kvalues []float64) ([]float64, float64) {
	var (
		ref  [][]float64
		comp [][]float64
	)

	// default value
	if len(refData)*len(compData) == 0 {
		ref = dc.refData
		comp = dc.compData
	}

	// RGBtoLab Conversion
	convertRGBtoLab := func(data [][]float64) [][]float64 {
		labs := make([][]float64, 0)

		for _, rgb := range data {

			// r / 255.0
			rgbNorm := make([]float64, 0)
			for _, signal := range rgb {
				rgbNorm = append(rgbNorm, signal/255.0)
			}

			lab := dc.RGBtoLabConversion(colorSpace, rgbNorm)
			labs = append(labs, lab)
		}

		return labs
	}

	refLabs := convertRGBtoLab(ref)
	compLabs := convertRGBtoLab(comp)

	// calculate deltaE
	results, average := dc.CalculateDeltaE(refLabs, compLabs, kvalues)

	// return results and average
	return results, average
}

/*
CalculateDeltaE
	in	:ref, comp [][]float64, kvalues []float64
	out	;[]float64, float64
*/
func (dc *deltaEvaluationController) CalculateDeltaE(ref, comp [][]float64, kvalues []float64) ([]float64, float64) {
	results := make([]float64, 0)

	// initialize deltaE calculator
	calculator := tools.NewDeltaLabCalculator()

	// calculate deltaE
	var sum float64
	for index := 0; index < 24; index++ {
		result := calculator.DeltaLab(ref[index], comp[index], kvalues)
		sum += result

		// stock the result
		results = append(results, result)
	}

	// average
	deltaEAve := sum / 24.0

	// update
	dc.deltaEResults = results

	// return
	return results, deltaEAve
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
