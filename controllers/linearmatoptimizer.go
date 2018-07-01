package controllers

import (
	"PixelTool_RC1/models"
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
	Run(paramIndex, trial int, linearMatElm []float64)
}

//
type linearMatrixOptimizer struct {
	orgElm       []float64
	devColorCode []models.ColorCode
	refColorCode []models.ColorCode

	numOfTrial int

	deltaEEvalController DeltaEvaluationController
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
	obj.deltaEEvalController = NewDeltaEvaluationController()
	obj.colorChartController = NewColorChartController()

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

/*
SetRefColorCode
*/
func (lo *linearMatrixOptimizer) SetRefColorCode(filepath string) bool {
	lo.refColorCode = lo.colorChartController.RunStandad(filepath)

	if len(lo.refColorCode) == 0 {
		return false
	}

	return true
}

/*
MakeDataSet :
*/
func (lo *linearMatrixOptimizer) makeVals(paramIndex int, elm []float64) []float64 {
	parameters := make([]float64, 0)
	/*
		0 :a
		1 :b
		2 :c
		3 :d
		4 :e
		5 :f
	*/

	// skip the element if the index matched to paraIndex
	for index, data := range elm {
		if index != paramIndex {
			parameters = append(parameters, data)
		}
	}

	// return
	return parameters
}

func (lo *linearMatrixOptimizer) evaluateDeltaE(linearMatElm []float64) {
	// initialize results
	devColorCodes := make([]models.ColorCode, 0)

	// check number of trial
	if lo.numOfTrial == 0 {
		devColorCodes = lo.colorChartController.RunDevice(
			lo.dataSet.linearMat,
			lo.dataSet.dataPath,
			lo.dataSet.devQE,
			lo.dataSet.ill,
			lo.dataSet.gamma,
			linearMatElm,
		)
	}
	// not initial
	devColorCodes = lo.colorChartController.RunDeiceBatch(
		lo.dataSet.gamma,
		linearMatElm,
	)

	// store the value
	lo.devColorCode = devColorCodes

	/*
		// serialize data
		for _, data := range devColorCodes {
			code := data.SerializeData()

		}
	*/

	// update flog
	lo.numOfTrial++

}

/*
Run :
*/
func (lo *linearMatrixOptimizer) Run(paramIndex, trial int, linearMatElm []float64) {
	// --- Step-1 ----

}
