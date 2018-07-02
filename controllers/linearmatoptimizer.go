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
VariableSet :var set
*/
type VariableSet struct {
	Index     int
	InitValue float64 // initial value
	MaxValue  float64 // sweep stop value
	MinValue  float64 // min value
	Variables []float64
}

/*
LinearMatrixOptimizer :linear mat optimizer
*/
type LinearMatrixOptimizer interface {
	SetEnv(linearMatDataPath, dataPath, deviceQEDataPath string, lightSource models.IlluminationCode, gamma float64) bool
	SetRefColorCode(filepath string) bool
	Run(trial int, linearMatElm []float64)
}

//
type linearMatrixOptimizer struct {
	orgElm       []float64
	devColorCode []models.ColorCode
	refColorCode []models.ColorCode

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
	obj.deltaEvalController = NewDeltaEvaluationController()
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

func (lo *linearMatrixOptimizer) makeVariableSet(elm []float64) {
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
		step := (targetMaxValue - targetMinValue) / 100.0

		// set initial value
		variable := targetMinValue

		// make variable set
		for j := 0; j < 100; j++ {
			// stocker
			stocker := elm

			// calculate new value
			variable += step

			// upfate slice
			stocker = append(stocker[:pos+1], stocker[pos:]...)
			stocker[pos] = variable

			fmt.Println(stocker)
		}
	}
}

/*
MakeDataSet :
*/
func (lo *linearMatrixOptimizer) makeVariable(paramIndex int, elm []float64) *VariableSet {
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

	// make variable data set
	varSet := new(VariableSet)
	varSet.Index = paramIndex
	varSet.InitValue = elm[paramIndex]
	varSet.Variables = parameters

	if (varSet.InitValue + 0.5*varSet.InitValue) > 1.0 {
		varSet.MaxValue = 1.0
	} else {
		varSet.MaxValue = varSet.InitValue + 0.5*varSet.InitValue
	}

	if (varSet.InitValue - 0.5*varSet.InitValue) < 0.0 {
		varSet.MinValue = 0.0
	} else {
		varSet.MinValue = varSet.InitValue - 0.5*varSet.InitValue
	}

	// return
	return varSet
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

	// update flog
	lo.numOfTrial++

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
Run :
*/
func (lo *linearMatrixOptimizer) Run(trial int, linearMatElm []float64) {

	// --- Step-1 ----
	// make variable set
	lo.makeVariableSet(linearMatElm)

}
