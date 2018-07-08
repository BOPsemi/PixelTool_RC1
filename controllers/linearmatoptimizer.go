package controllers

import (
	"PixelTool_RC1/models"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	learning = 0.1
	epsilon  = 0.001
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

	Run(elm []float64, targetDeltaE float64, trialNum int, deltaP float64, bachSize int) bool

	OptimizedLinearMatrix() []float64
	FinalDeltaEInfo() (deltaE []float64, deltaEAve float64)
	//Run(splitNum, trial int, linearMatElm []float64)
}

//
type linearMatrixOptimizer struct {
	refBitCode   [][]float64 // reference patch bit codes
	optElm       []float64   // optimized Linear matrix
	optDeltaEAve float64     // optimized DeltaE ave
	optDeltaE    []float64   // optimized DeltaE slice

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

	obj.optElm = make([]float64, 6)
	obj.optDeltaE = make([]float64, 24)
	obj.optDeltaEAve = 0.0

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
OptimizedLinearMatrix :getter, return the optimized linear matrix
*/
func (lo *linearMatrixOptimizer) OptimizedLinearMatrix() []float64 {
	return lo.optElm
}

/*
FinalDeltaEInfo :getter, return the final deltaE and deltaE Ave.
*/
func (lo *linearMatrixOptimizer) FinalDeltaEInfo() (deltaE []float64, deltaEAve float64) {
	return lo.optDeltaE, lo.optDeltaEAve
}

/*
SetRefColorCode :
	in	;filepath
	out ;bool
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

// calculate device response
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

// deltaE calculator
func (lo *linearMatrixOptimizer) deltaECalculator(elm []float64) ([]float64, float64) {

	devBitCode := lo.calculateDevResponse(elm)

	deltaEEvalController := NewDeltaEvaluationController()
	kvalues := []float64{1.0, 1.0, 1.0}
	deltaE, deltaEAve := deltaEEvalController.RunDeltaEEvaluation(models.SRGB, lo.refBitCode, devBitCode, kvalues)

	return deltaE, deltaEAve

}

// calculate gradient at elm points
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

		// copy
		copy(plusElm, elm)
		copy(minusElm, elm)

		// update element value
		plusElm[elmNumber] = plusVal
		minusElm[elmNumber] = minusVal

		// calculate deltaE
		_, plusDeltaEAve := lo.deltaECalculator(plusElm)
		_, minusDeltaEAve := lo.deltaECalculator(minusElm)

		// calculate div
		divDeltaE := (plusDeltaEAve - minusDeltaEAve) / deltaVal

		// stock
		gradDeltaE = append(gradDeltaE, divDeltaE)
	}

	// return
	return gradDeltaE
}

// generate random variations for mini bach calculation
func (lo *linearMatrixOptimizer) randVarGenerator(rangePer, oriValue float64) float64 {
	rand.Seed(time.Now().UnixNano())

	// calculate max and min value
	max := oriValue + oriValue*rangePer*0.01
	min := oriValue - oriValue*rangePer*0.01

	// return value
	return rand.Float64()*(max-min) + min

}

/*
RunAdaGrad : run AdaGrad
	in	;elm []float64, targetDeltaE float64, trialNum int, deltaP float64, bachSize int
	out	;bool
*/
func (lo *linearMatrixOptimizer) Run(elm []float64, targetDeltaE float64, trialNum int, deltaP float64, bachSize int) bool {

	// --- Step-0 ---
	// make local elm slice
	divStocker := make([]float64, len(elm)) // gradient stocker
	localElm := make([]float64, len(elm))   // copy of elm
	copy(localElm, elm)

	// --- Step-1 ---
	// Start learning
	for trial := 0; trial < trialNum; trial++ {

		// dataSet stocker
		dataSetStocker := make([]models.DataSet, 0)

		// --- Step-2 ----
		// make data set
		for elmIndex := 0; elmIndex < len(localElm); elmIndex++ {
			// make new elm matrix
			newElm := make([]float64, len(localElm))
			copy(newElm, localElm)

			for trial := 0; trial < bachSize; trial++ {
				// randomize the data
				randElmValue := lo.randVarGenerator(deltaP*20, localElm[elmIndex])
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

		// --- Step-3 ---
		// randomize array order
		randomizedDataSet := make([]models.DataSet, len(dataSetStocker))
		copy(randomizedDataSet, dataSetStocker)

		rand.Seed(time.Now().UnixNano())
		n := len(randomizedDataSet)
		for i := n - 1; i >= 0; i-- {
			j := rand.Intn(i + 1)
			randomizedDataSet[i], randomizedDataSet[j] = randomizedDataSet[j], randomizedDataSet[i]
		}

		// --- Step-4 ---
		// calculate mini-bach
		bachSum := make([]float64, 6)

		for _, data := range randomizedDataSet {
			for index := 0; index < len(localElm); index++ {
				// for AdaGrad
				divStocker[index] += data.DivDeltaE[index] * data.DivDeltaE[index]

				// mini bach
				bachSum[index] += learning * data.DivDeltaE[index]
			}
		}

		// --- Step-5 ---
		// introduce next feedback
		nextElm := make([]float64, 6)
		copy(nextElm, localElm)
		for index := 0; index < len(nextElm); index++ {
			// for Ada Grad
			learningRate := learning / math.Sqrt(divStocker[index]+epsilon)
			nextElm[index] = localElm[index] - learningRate*bachSum[index]

			// check minus value
			if nextElm[index] < 0.0 {
				nextElm[index] = localElm[index]
			}
		}

		// --- Step-5 ---
		// update elm
		copy(localElm, nextElm)

		// --- Step-6 ---
		//update final value
		updatedDeltaE, updatedDeltaEAve := lo.deltaECalculator(localElm)

		lo.optDeltaEAve = updatedDeltaEAve
		lo.optDeltaE = updatedDeltaE
		lo.optElm = localElm

		//fmt.Println(localElm, updatedDeltaEAve)

	}
	/*
		fmt.Println("----")
		fmt.Println(lo.optElm, lo.optDeltaEAve, lo.optDeltaE)
	*/

	return true
}

/*
Run :
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
*/
