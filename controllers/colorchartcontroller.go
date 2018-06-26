package controllers

import (
	"PixelTool_RC1/models"
	"PixelTool_RC1/util"
	"strconv"
)

const (
	start          int   = 400
	stop           int   = 700
	step           int   = 5
	refPathNo      int   = 19
	refPatchLevel  uint8 = 243
	refPathNoForWB int   = 22

	illuminationD65FileName     = "illumination_D65.csv"
	illuminationIllAFileName    = "illumination_A.csv"
	macbethColorCheckerFileName = "Macbeth_Color_Checker.csv"
)

var (
	lithSource             models.IlluminationCode  // light source info
	stdColorCodes          []models.ColorCode       // color code stocker for standar CC
	devColorCodes          []models.ColorCode       // color code stocker for device CC
	rawResponses           []models.ChannelResponse // raw data of channel response
	rawLinearizedResponses [][]float64              // raw data after linear matrix calculation
	rawWBData              [][]float64              // raw data after white balance
	linearMatrixElm        []float64                // linear matrix elemets
	filepath               map[string]string        // file paths for calculation
)

func init() {
	stdColorCodes = make([]models.ColorCode, 0)
	devColorCodes = make([]models.ColorCode, 0)
	rawResponses = make([]models.ChannelResponse, 0)
	rawLinearizedResponses = make([][]float64, 0)
	rawWBData = make([][]float64, 0)
	linearMatrixElm = make([]float64, 0)
	filepath = make(map[string]string, 0)
}

/*
ColorChartController :linear matrix optimizer
*/
type ColorChartController interface {
	RunDevice(linearMatDataPath, deviceQEDataPath string, lightSource models.IlluminationCode, initFlag bool, gamma float64, linearMat []float64) []models.ColorCode
}

type colorChartController struct {
	resController ResponseController
}

/*
NewLinearMatrixOptimizer :initializer
*/
func NewLinearMatrixOptimizer() ColorChartController {
	obj := new(colorChartController)

	// initialize response controller
	obj.resController = NewResponseController()

	return obj
}

func (cg *colorChartController) RunDevice(linearMatDataPath, deviceQEDataPath string, lightSource models.IlluminationCode, initFlag bool, gamma float64, linearMat []float64) []models.ColorCode {
	if initFlag {
		cg.setEnv(linearMatDataPath, deviceQEDataPath, lightSource)
		cg.resController.ReadResponseData(filepath)
	}

	return cg.runDevice(gamma, linearMat)
}

/*
SetEnv	:setup enviroment
	in	:info
	out
*/
func (cg *colorChartController) setEnv(linearMatDataPath, deviceQEDataPath string, ill models.IlluminationCode) {
	// light source setup
	lithSource = ill

	// read Linear Matrix elements
	linearMatrixElm = cg.readLinearMatElmFromCSV(linearMatDataPath)

	// set file path
	filepath = cg.setReadingFile(deviceQEDataPath)

}

// read csv
func (cg *colorChartController) readLinearMatElmFromCSV(filepath string) []float64 {
	iohandler := util.NewIOUtil()
	elm := make([]float64, 0)
	if data, ok := iohandler.ReadCSVFile(filepath); ok {
		for _, strs := range data {
			value, _ := strconv.ParseFloat(strs[0], 64)
			elm = append(elm, value)
		}
	}

	return elm
}

// setup reading files
func (cg *colorChartController) setReadingFile(deviceQEpath string) map[string]string {
	filepath := make(map[string]string, 0)

	// get current path
	dirHandler := util.NewDirectoryHandler()
	path := dirHandler.GetCurrentDirectoryPath() + "/data/"

	// illumination data path
	filepath["D65"] = path + illuminationD65FileName
	filepath["IllA"] = path + illuminationIllAFileName

	// raw data
	filepath["DeviceQE"] = deviceQEpath
	filepath["ColorChecker"] = path + macbethColorCheckerFileName

	return filepath
}

// calculate device response
func (cg *colorChartController) calculateDeviceResponse(gamma float64) {
	if ok, responses := cg.resController.CalculateChannelResponse(
		lithSource,
		start,
		stop,
		step,
		refPathNo,
	); ok {
		for _, data := range responses {
			if status, result := cg.resController.CalculateGammaCorrection(gamma, &data); status {
				rawResponses = append(rawResponses, *result)
			}
		}
	}
}

// calculate linear matrix
func (cg *colorChartController) calculateLinearMatrix(linerMat []float64) {
	for _, data := range rawResponses {
		// change data format to linear matrix calculation
		grgbrb := []float64{
			data.Gr,
			data.Gb,
			data.R,
			data.B,
		}

		// calcualte linear matrix
		response := cg.resController.CalculateLinearMatrix(linerMat, grgbrb)

		// stock
		rawLinearizedResponses = append(rawLinearizedResponses, response)
	}
}

// calculate white balance
func (cg *colorChartController) calculateWhiteBalance(wbRefPatchNo int) {
	redGain, blueGain := cg.resController.CalculateWhiteBalanceGain(rawLinearizedResponses[wbRefPatchNo-1])
	for _, data := range rawLinearizedResponses {

		// calculate raw data
		red := data[0] * redGain
		green := data[1]
		blue := data[2] * blueGain

		// make raw data
		rawdata := []float64{red, green, blue}
		rawWBData = append(rawWBData, rawdata)

	}

}

// convert 8bit data
func (cg *colorChartController) convert8BitData() {
	// init degitizer
	digitizer := util.NewDigitizer()

	for index, data := range rawWBData {
		red8bit := digitizer.D8bitDigitizeData(data[0], refPatchLevel)
		green8bit := digitizer.D8bitDigitizeData(data[1], refPatchLevel)
		blue8bit := digitizer.D8bitDigitizeData(data[2], refPatchLevel)

		// patch name
		pname := models.MacbethColorCode(index).String()

		// create color code model
		colorcode := models.SetColorCode(index+1, pname, red8bit, green8bit, blue8bit, 255)

		// update
		devColorCodes = append(devColorCodes, *colorcode)
	}
}

func (cg *colorChartController) runDevice(gamma float64, linearMat []float64) []models.ColorCode {

	// --- 2nd stage ---
	// calculate device response
	cg.calculateDeviceResponse(gamma)

	// --- 3rd stage ---
	// calculate linear matrix
	cg.calculateLinearMatrix(linearMat)

	// --- 4th stage ---
	// white balance
	cg.calculateWhiteBalance(refPathNoForWB)

	// --- 5th stage ---
	// degitize
	cg.convert8BitData()

	return devColorCodes

}
