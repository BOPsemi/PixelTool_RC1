package controllers

import (
	"PixelTool_RC1/models"
	"fmt"
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
	RunAdaGrad(elm []float64, targetDeltaE float64, deltaP float64, bachSize int)
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

func (lo *linearMatrixOptimizer) deltaECalculator(elm []float64) ([]float64, float64) {

	devBitCode := lo.calculateDevResponse(elm)

	deltaEEvalController := NewDeltaEvaluationController()
	kvalues := []float64{1.0, 1.0, 1.0}
	deltaE, deltaEAve := deltaEEvalController.RunDeltaEEvaluation(models.SRGB, lo.refBitCode, devBitCode, kvalues)

	return deltaE, deltaEAve

}

func (lo *linearMatrixOptimizer) gradient(elm []float64, deltaParcentage float64) []float64 {

	// stocker
	gradDeltaE := make([]float64, 0)

	// sweep elm number
	for elmNumber, elmValue := range elm {
		// calculate shift range
		plusVal := elmValue + elmValue*deltaParcentage*0.01
		minusVal := elmValue - elmValue*deltaParcentage*0.01
		deltaVal := plusVal - minusVal

		// make shifted elm
		plusElm := make([]float64, len(elm))
		minusElm := make([]float64, len(elm))

		copy(plusElm, elm)
		copy(minusElm, elm)

		plusElm[elmNumber] = plusVal
		minusElm[elmNumber] = minusVal

		// define deltaE calculator
		/*
			calculateDeltaE := func(linerElm []float64) (eachDeltaE []float64, deltaEAve float64) {
				devBitCode := lo.calculateDevResponse(linerElm)

				deltaEEvalController := NewDeltaEvaluationController()
				kvalues := []float64{1.0, 1.0, 1.0}
				deltaE, deltaEAve := deltaEEvalController.RunDeltaEEvaluation(models.SRGB, lo.refBitCode, devBitCode, kvalues)

				return deltaE, deltaEAve
			}
		*/

		// calculate gradient
		//_, plusDeltaEAve := calculateDeltaE(plusElm)
		//_, minusDeltaEAve := calculateDeltaE(minusElm)

		_, plusDeltaEAve := lo.deltaECalculator(plusElm)
		_, minusDeltaEAve := lo.deltaECalculator(minusElm)

		divDeltaE := (plusDeltaEAve - minusDeltaEAve) / deltaVal

		// stock
		gradDeltaE = append(gradDeltaE, divDeltaE)
	}

	// return
	return gradDeltaE
}

func (lo *linearMatrixOptimizer) randVarGenerator(rangePer, oriValue float64) float64 {
	rand.Seed(time.Now().UnixNano())

	max := oriValue + oriValue*rangePer*0.01
	min := oriValue - oriValue*rangePer*0.01

	return rand.Float64()*(max-min) + min

}

/*
RunAdaGrad : run AdaGrad
*/
func (lo *linearMatrixOptimizer) RunAdaGrad(elm []float64, targetDeltaE float64, deltaP float64, bachSize int) {
	// dataSet stocker
	dataSetStocker := make([]models.DataSet, 0)

	// --- Step-0 ---
	// make local elm slice
	localElm := make([]float64, 6)
	copy(localElm, elm)

	// --- Step-1 ----
	// make data set
	for elmIndex := 0; elmIndex < len(localElm); elmIndex++ {
		// make new elm matrix
		newElm := make([]float64, len(localElm))
		copy(newElm, localElm)

		for trial := 0; trial < bachSize; trial++ {
			// randomize the data
			randElmValue := lo.randVarGenerator(deltaP*10, localElm[elmIndex])
			newElm[elmIndex] = randElmValue

			// calculate deltaE with newElm
			deltaEArray, deltaEAve := lo.deltaECalculator(newElm)

			// calculate gradient
			gradDeltaE := lo.gradient(newElm, deltaP)

			// calculate div
			divDeltaE := make([]float64, 0)
			for _, gradData := range gradDeltaE {
				div := (deltaEAve - targetDeltaE) * gradData
				divDeltaE = append(divDeltaE, div)
			}

			// make data set
			dataSet := new(models.DataSet)
			dataSet.DeltaEAve = deltaEAve
			dataSet.DivDeltaE = divDeltaE
			dataSet.Elm = make([]float64, 6)
			copy(dataSet.Elm, newElm)

			dataSet.DeltaE = make([]float64, 24)
			copy(dataSet.DeltaE, deltaEArray)

			// stock
			dataSetStocker = append(dataSetStocker, *dataSet)
		}

	}

	// --- Step-2 ---
	// randomize array order

	// --- Step-3 ---
	// calculate mini-bach

	// --- Step-4 ---
	// introduce next feedback

	// --- Step-5 ---
	// update elm

}

/*
Run :
*/
func (lo *linearMatrixOptimizer) Run(splitNum, trial int, linearMatElm []float64) {

	// --- Step-0 ---
	// save original elm mat
	lo.orgElm = make([]float64, 6)
	copy(lo.orgElm, linearMatElm)

	// --- Step-1 ---
	// initialize elm update
	elm := make([]float64, 6)
	copy(elm, lo.orgElm)

	// starage of min deltaE condition
	minDeltaEdataSet := new(models.DataSet)
	minDeltaEAve := 10.0

	for loop := 0; loop < trial; loop++ {
		for elmIndex := 0; elmIndex < len(lo.orgElm); elmIndex++ {

			// calculate paramter sweep range and step
			initElmValue := elm[elmIndex]
			maxElmValue := initElmValue + initElmValue*0.8
			if maxElmValue > 1.0 {
				maxElmValue = 1.0
			}
			minElmValue := initElmValue - initElmValue*0.8
			if minElmValue < 0.0 {
				minElmValue = 0.0
			}
			step := (maxElmValue - minElmValue) / float64(splitNum)

			// generate data set
			deltaEDataSet := make([]models.DataSet, 0)
			variable := 0.0

			for index := 0; index < splitNum; index++ {

				// make new elm matrix for sweeping
				newElm := make([]float64, 6)
				copy(newElm, elm)

				variable += step
				newElm[elmIndex] = minElmValue + variable

				//  calculate deltaE
				if newElm[elmIndex] == 0.0 || newElm[elmIndex] == 1.0 {
					break
				}

				// calculate device response by using elm (linea matrix elements)
				devBitCode := lo.calculateDevResponse(newElm)

				// calculate deltaE
				deltaEEvalController := NewDeltaEvaluationController()
				kvalues := []float64{1.0, 1.0, 1.0}
				deltaE, deltaEAve := deltaEEvalController.RunDeltaEEvaluation(models.SRGB, lo.refBitCode, devBitCode, kvalues)

				// create dataSet object
				data := new(models.DataSet)
				data.DeltaE = deltaE
				data.DeltaEAve = deltaEAve
				data.Elm = make([]float64, 6)
				copy(data.Elm, newElm)

				deltaEDataSet = append(deltaEDataSet, *data) // stock elm data

				if data.DeltaEAve < minDeltaEAve {
					minDeltaEdataSet = data
					minDeltaEAve = minDeltaEdataSet.DeltaEAve
					elm = minDeltaEdataSet.Elm
				}
			}
		}

	}
	fmt.Println(elm, minDeltaEAve)
}
