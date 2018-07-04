package controllers

import (
	"PixelTool_RC1/models"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

/*
DataSet :definition of data set
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
	orgElm        []float64          // original linear matrix elements
	refBitCode    [][]float64        // reference patch bit codes
	dataElmSet    [][][]float64      // parametric elm data
	dataDeltaESet [][]models.DataSet // parametric deltaE data

	numOfTrial int

	// setting information struct
	settingInfo struct {
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

	// upload simualtio enviroment info
	lo.settingInfo.linearMat = linearMatDataPath
	lo.settingInfo.dataPath = dataPath
	lo.settingInfo.devQE = deviceQEDataPath
	lo.settingInfo.ill = lightSource
	lo.settingInfo.gamma = gamma

	// return
	return true
}

// serializer
func (lo *linearMatrixOptimizer) serializeData(data models.ColorCode) []float64 {
	// initialize buffer
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

// make deltaE data set
func (lo *linearMatrixOptimizer) makeDeltaEDataSet(dataSet [][][]float64) [][]models.DataSet {
	// init stocker
	dataElmSet := make([][]models.DataSet, 0)

	// start to make data set
	for _, elmSet := range dataSet {

		// init stocker
		dataElm := make([]models.DataSet, 0)

		for _, elm := range elmSet {
			// calculate device response by using elm (linea matrix elements)
			devBitCode := lo.calculateDevResponse(elm)

			// calculate deltaE
			deltaEEvalController := NewDeltaEvaluationController()
			kvalues := []float64{1.0, 1.0, 1.0}
			deltaE, deltaEAve := deltaEEvalController.RunDeltaEEvaluation(models.SRGB, lo.refBitCode, devBitCode, kvalues)

			// create dataSet object
			data := new(models.DataSet)
			data.DeltaE = deltaE
			data.DeltaEAve = deltaEAve
			data.Elm = make([]float64, 6)
			copy(data.Elm, elm)

			// stock
			dataElm = append(dataElm, *data)
		}

		dataElmSet = append(dataElmSet, dataElm)
	}

	return dataElmSet
}

func (lo *linearMatrixOptimizer) calculateDevResponse(linearMatElm []float64) [][]float64 {
	// initialize results
	ccController := NewColorChartController()

	// calculate device response
	// return all channel data
	devColorCodes := ccController.RunDevice(
		lo.settingInfo.linearMat,
		lo.settingInfo.dataPath,
		lo.settingInfo.devQE,
		lo.settingInfo.ill,
		lo.settingInfo.gamma,
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

// shuffle the data set
func (lo *linearMatrixOptimizer) shuffleDeltaEDataSet(deltaEDatSet []models.DataSet) []models.DataSet {
	newDeltaEDataSet := make([]models.DataSet, len(deltaEDatSet))
	copy(newDeltaEDataSet, deltaEDatSet)

	// random number
	rand.Seed(time.Now().UnixNano())

	number := len(newDeltaEDataSet)
	for index := number - 1; index >= 0; index-- {
		jump := rand.Intn(index + 1)
		newDeltaEDataSet[index], newDeltaEDataSet[jump] = newDeltaEDataSet[jump], newDeltaEDataSet[index]
	}

	return newDeltaEDataSet
}

// gradient
func (lo *linearMatrixOptimizer) gradient(elmIndex int, targetDeltaEAve float64, deltaEDataSet []models.DataSet, bachSize int) float64 {
	result := 0.0
	for _, data := range deltaEDataSet[0:bachSize] {
		//result += (data.DeltaEAve - targetDeltaEAve) * data.Elm[elmIndex]
		result += (data.DeltaEAve - targetDeltaEAve)
	}

	return result
}

/*
Run :
*/
func (lo *linearMatrixOptimizer) Run(splitNum, trial int, linearMatElm []float64) {

	// --- Step-0 ---
	// save original elm mat
	lo.orgElm = make([]float64, 6)
	copy(lo.orgElm, linearMatElm)

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
	lo.dataElmSet = dataSet

	// --- Step-3 ---
	// meke data set
	lo.dataDeltaESet = lo.makeDeltaEDataSet(dataSet)

	// --- Step-4 ---
	// optimization

	/*
		dataDeltaESetA := lo.dataDeltaESet[0] // weep a- parameter
		shuffledDataSet := lo.shuffleDeltaEDataSet(dataDeltaESetA)
	*/

	elmIndex := 0
	targetDeltaE := 2.0
	bachSize := 10
	learningRateC := 0.01
	epsilon := 0.0001

	grad2Integ := 0.0

	for trial := 0; trial < 100; trial++ {
		for index := 0; index < 6; index++ {
			if index != elmIndex {
				shuffledDataSet := lo.shuffleDeltaEDataSet(lo.dataDeltaESet[elmIndex])
				grad := lo.gradient(index, targetDeltaE, shuffledDataSet, bachSize)
				grad2 := grad * grad
				grad2Integ += grad2

				learningRate := learningRateC / (math.Sqrt(grad2Integ) + epsilon)
				update := -(learningRate * grad)

				updatedElm := lo.orgElm[index] + update

				fmt.Println(learningRate, lo.orgElm[index], updatedElm)

			}
		}
		fmt.Println("---", trial, "---")
	}

	/*
		resultA := lo.gradient(1, 2.0, shuffledDataSet, 5)

		fmt.Println(resultA)
	*/

}
