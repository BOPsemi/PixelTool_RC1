package controllers

import (
	"PixelTool_RC1/models"
	"fmt"
	"strconv"
)

/*
DataSet
*/
type DataSet struct {
	Index     int       // parameter index
	X         float64   // parameter value of a/b/c/d/e/f
	DeltaEAve float64   // deltaE Average
	Delta     []float64 // deltaE
}

/*
LinearMatrixOptimizer :linear mat optimizer
*/
type LinearMatrixOptimizer interface {
	SetEnv(linearMatDataPath, dataPath, deviceQEDataPath string, lightSource models.IlluminationCode, gamma float64) bool
	SetRefColorCode(filepath string) bool
	Run(splitNum, trial int, linearMatElm []float64)
}

//
type linearMatrixOptimizer struct {
	orgElm     []float64
	refBitCode [][]float64

	numOfTrial int

	deltaEvalController  DeltaEvaluationController
	colorChartController ColorChartController

	dataSet struct {
		linearMat string
		dataPath  string
		devQE     string
		ill       models.IlluminationCode
		gamma     float64
	}
}

/*
NewLinearMatrixOptimizer :initializer
*/
func NewLinearMatrixOptimizer() LinearMatrixOptimizer {
	obj := new(linearMatrixOptimizer)

	// initialize properties
	obj.numOfTrial = 0

	return obj
}

/*
SetEnv : set simulation enviroment
*/
func (lo *linearMatrixOptimizer) SetEnv(linearMatDataPath, dataPath, deviceQEDataPath string, lightSource models.IlluminationCode, gamma float64) bool {

	// make data set
	lo.dataSet.linearMat = linearMatDataPath
	lo.dataSet.dataPath = dataPath
	lo.dataSet.devQE = deviceQEDataPath
	lo.dataSet.ill = lightSource
	lo.dataSet.gamma = gamma

	// return
	return true
}

// serializer
func (lo *linearMatrixOptimizer) serializeData(data models.ColorCode) []float64 {
	rgbData := make([]float64, 0)

	// serialize
	rawRGBdata := data.SerializeData()

	// extract value and parse the data to float64
	for index, rgb := range rawRGBdata {
		if index > 1 && index < 5 {
			value, err := strconv.ParseFloat(rgb, 64)
			if err == nil {
				rgbData = append(rgbData, value)
			}
		}
	}

	// return
	return rgbData
}

/*
SetRefColorCode
*/
func (lo *linearMatrixOptimizer) SetRefColorCode(filepath string) bool {
	bitcodes := make([][]float64, 0)

	// calculate ref 8bit RGB code
	refColorCode := models.ReadColorCode(filepath)
	for _, rawdata := range refColorCode {
		rgb := lo.serializeData(rawdata)
		bitcodes = append(bitcodes, rgb)
	}

	// update
	lo.refBitCode = bitcodes

	return true
}

func (lo *linearMatrixOptimizer) makeVariableSet(splitNum int, elm []float64) [][]float64 {
	varDataSet := make([][]float64, 0)

	for pos := 0; pos < len(elm); pos++ {

		// calculate parameter sweep
		targetInitValue := elm[pos]
		targetMaxValue := targetInitValue + targetInitValue*0.5
		if targetMaxValue > 1.0 {
			targetMaxValue = 1.0
		}
		targetMinValue := targetInitValue - targetInitValue*0.5
		if targetMinValue < 0.0 {
			targetMinValue = 0.0
		}

		// calculate min value and step
		variables := targetMinValue
		step := (targetMaxValue - targetMinValue) / float64(splitNum)

		for index := 0; index < splitNum; index++ {

			// make stocker
			stocker := make([]float64, len(elm))
			copy(stocker, elm)

			variables += step
			stocker[pos] = variables

			//fmt.Println(stocker)
			varDataSet = append(varDataSet, stocker)
		}
	}

	return varDataSet
}

// make data set
func (lo *linearMatrixOptimizer) makeDataSet(index, splitNum int, variableSet [][]float64) [][]float64 {
	stocker := make([][]float64, splitNum)
	startPOS := index * splitNum
	endPOS := (index + 1) * splitNum
	copy(stocker, variableSet[startPOS:endPOS])

	return stocker
}

func (lo *linearMatrixOptimizer) calculateDevResponse(linearMatElm []float64) [][]float64 {
	// initialize results
	ccController := NewColorChartController()

	// calculate device response
	// return all channel data
	devColorCodes := ccController.RunDevice(
		lo.dataSet.linearMat,
		lo.dataSet.dataPath,
		lo.dataSet.devQE,
		lo.dataSet.ill,
		lo.dataSet.gamma,
		linearMatElm,
	)

	// serialize the data
	bitcodes := make([][]float64, 0)
	for _, rawdata := range devColorCodes {
		rgb := lo.serializeData(rawdata)
		bitcodes = append(bitcodes, rgb)
	}

	// return
	return bitcodes
}

/*
Run :
*/
func (lo *linearMatrixOptimizer) Run(splitNum, trial int, linearMatElm []float64) {

	// --- Step-1 ----
	// make variable set
	variableSet := lo.makeVariableSet(splitNum, linearMatElm)

	// --- Step-2 ---
	// make data set
	dataSet := make([][][]float64, 0)
	for index := 0; index < len(linearMatElm); index++ {
		data := lo.makeDataSet(index, splitNum, variableSet)
		dataSet = append(dataSet, data)
	}

	// --- Step-3 ---
	// calculate device Lab
	minDeltaEAve := 10.0
	minSetNum := 0
	minElmNum := 0
	minDeltaE := make([]float64, 0)
	minElm := make([]float64, 0)

	for setNum, elmSet := range dataSet {

		// iniit stocker

		for elmNum, elm := range elmSet {

			// calculate device response by using elm (linea matrix elements)
			devBitCode := lo.calculateDevResponse(elm)

			// calculate deltaE
			deltaEEvalController := NewDeltaEvaluationController()
			kvalues := []float64{1.0, 1.0, 1.0}
			deltaE, deltaEAve := deltaEEvalController.RunDeltaEEvaluation(models.SRGB, lo.refBitCode, devBitCode, kvalues)

			if deltaEAve < minDeltaEAve {
				minDeltaEAve = deltaEAve
				minSetNum = setNum
				minElmNum = elmNum
				minDeltaE = deltaE
				minElm = elm
			}

		}
	}

	fmt.Println(minDeltaE, minDeltaEAve, minElmNum, minSetNum, minElm)

}
