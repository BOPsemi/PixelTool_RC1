package controllers

import (
	"PixelTool_RC1/models"
	"PixelTool_RC1/util"
	"math"
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
	std24ColorChartName         = "std_24_ColorChart"
	dev24ColorChartName         = "dev_24_ColorChart"
)

/*
ColorChartController :linear matrix optimizer
*/
type ColorChartController interface {
	RunDevice(linearMatDataPath, dataPath, deviceQEDataPath string, lightSource models.IlluminationCode, gamma float64, linearMat []float64) []models.ColorCode
	RunStandad(filepath string) []models.ColorCode
	RunDeiceBatch(gamma float64, linearMat []float64) []models.ColorCode

	GenerateColorPatchImage(data []models.ColorCode, filesavepath, dirname string, width, height int) bool
	Generate24MacbethColorChart(dev bool, filesavepath string) bool
	SaveColorPatchCSVData(data []models.ColorCode, savepath, filename string) bool
}

type colorChartController struct {
	resController ResponseController

	lithSource             models.IlluminationCode  // light source info
	stdColorCodes          []models.ColorCode       // color code stocker for standar CC
	devColorCodes          []models.ColorCode       // color code stocker for device CC
	devResponses           []models.ChannelResponse // original device response
	rawResponses           []models.ChannelResponse // raw data of channel response
	rawLinearizedResponses [][]float64              // raw data after linear matrix calculation
	rawWBData              [][]float64              // raw data after white balance
	linearMatrixElm        []float64                // linear matrix elemets
	filepath               map[string]string        // file paths for calculation
}

/*
NewColorChartController :initializer
*/
func NewColorChartController() ColorChartController {
	obj := new(colorChartController)

	// initialize properties
	obj.stdColorCodes = make([]models.ColorCode, 0)
	obj.devColorCodes = make([]models.ColorCode, 0)
	obj.rawResponses = make([]models.ChannelResponse, 0)
	obj.rawLinearizedResponses = make([][]float64, 0)
	obj.rawWBData = make([][]float64, 0)

	obj.linearMatrixElm = make([]float64, 0)
	obj.filepath = make(map[string]string, 0)

	return obj
}

func (cg *colorChartController) cleanUp() {
	cg.stdColorCodes = make([]models.ColorCode, 0)
	cg.devColorCodes = make([]models.ColorCode, 0)
	cg.rawResponses = make([]models.ChannelResponse, 0)
	cg.rawLinearizedResponses = make([][]float64, 0)
	cg.rawWBData = make([][]float64, 0)
}

func (cg *colorChartController) RunDevice(linearMatDataPath, dataPath, deviceQEDataPath string, lightSource models.IlluminationCode, gamma float64, linearMat []float64) []models.ColorCode {

	// setup simulation enviroment
	cg.setEnv(linearMatDataPath, dataPath, deviceQEDataPath, lightSource)

	// run and retrun result
	cg.cleanUp()
	return cg.runDevice(gamma, linearMat)
}

func (cg *colorChartController) RunDeiceBatch(gamma float64, linearMat []float64) []models.ColorCode {
	cg.cleanUp()
	return cg.runDevice(gamma, linearMat)
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
func (cg *colorChartController) setReadingFile(dataPath, deviceQEpath string) map[string]string {
	filepath := make(map[string]string, 0)

	var path string
	if dataPath == "" {
		// get current path
		dirHandler := util.NewDirectoryHandler()
		path = dirHandler.GetCurrentDirectoryPath() + "/data/"
	} else {
		path = dataPath
	}

	// illumination data path
	filepath["D65"] = path + illuminationD65FileName
	filepath["IllA"] = path + illuminationIllAFileName

	// raw data
	filepath["DeviceQE"] = deviceQEpath
	filepath["ColorChecker"] = path + macbethColorCheckerFileName

	return filepath
}

/*
SetEnv	:setup enviroment
	in	:info
	out
*/
func (cg *colorChartController) setEnv(linearMatDataPath, dataPath, deviceQEDataPath string, ill models.IlluminationCode) {
	// initalize response controller
	cg.resController = NewResponseController()

	// light source setup
	cg.lithSource = ill

	// read Linear Matrix elements
	cg.linearMatrixElm = cg.readLinearMatElmFromCSV(linearMatDataPath)

	// set file path
	cg.filepath = cg.setReadingFile(dataPath, deviceQEDataPath)

	// read filepath
	cg.resController.ReadResponseData(cg.filepath)

	// clean up
	cg.cleanUp()

}

// calculate device response
func (cg *colorChartController) calculateDeviceResponse(gamma float64) bool {

	if cg.resController == nil {
		return false
	}

	// calculate channel response
	if ok, responses := cg.resController.CalculateChannelResponse(
		cg.lithSource,
		start,
		stop,
		step,
		refPathNo,
	); ok {

		// update
		cg.devResponses = responses

		// gamma
		for _, data := range responses {
			if status, result := cg.resController.CalculateGammaCorrection(gamma, &data); status {
				cg.rawResponses = append(cg.rawResponses, *result)
			}
		}
	}

	// check result
	if len(cg.devResponses) == 0 {
		return false
	}
	if len(cg.rawResponses) == 0 {
		return false
	}

	return true
}

// calculate linear matrix
func (cg *colorChartController) calculateLinearMatrix(linerMat []float64) bool {
	/*
		for _, data := range cg.rawResponses {
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
			cg.rawLinearizedResponses = append(cg.rawLinearizedResponses, response)
		}
	*/

	for _, data := range cg.devResponses {
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
		cg.rawLinearizedResponses = append(cg.rawLinearizedResponses, response)
	}

	// check
	if len(cg.rawLinearizedResponses) == 0 {
		return false
	}

	return true

}

// calculate white balance
func (cg *colorChartController) calculateWhiteBalance(wbRefPatchNo int) bool {
	redGain, blueGain := cg.resController.CalculateWhiteBalanceGain(cg.rawLinearizedResponses[wbRefPatchNo-1])
	for _, data := range cg.rawLinearizedResponses {

		// calculate raw data
		red := data[0] * redGain
		green := data[1]
		blue := data[2] * blueGain

		// make raw data
		rawdata := []float64{red, green, blue}
		cg.rawWBData = append(cg.rawWBData, rawdata)

	}

	// check
	if len(cg.rawWBData) == 0 {
		return false
	}

	return true
}

// calculate gamma correlection
func (cg *colorChartController) calculateGammaCorrection(gamma float64) bool {

	powFunc := func(base, gamma float64) float64 {
		return math.Pow(base, gamma)
	}

	// stocker
	results := make([][]float64, 0)

	// calculate gamma correction
	for _, data := range cg.rawWBData {
		/*
			Red		:data[0]
			Green 	:data[1]
			Blue	:data[2]
		*/
		red := powFunc(data[0], gamma)
		green := powFunc(data[1], gamma)
		blue := powFunc(data[2], gamma)

		// stock the result
		results = append(results, []float64{red, green, blue})
	}

	// update
	cg.rawWBData = results

	// check
	if len(cg.rawWBData) == 0 {
		return false
	}

	return true
}

// convert 8bit data
func (cg *colorChartController) convert8BitData() bool {
	// init degitizer
	digitizer := util.NewDigitizer()

	for index, data := range cg.rawWBData {
		red8bit := digitizer.D8bitDigitizeData(data[0], refPatchLevel)
		green8bit := digitizer.D8bitDigitizeData(data[1], refPatchLevel)
		blue8bit := digitizer.D8bitDigitizeData(data[2], refPatchLevel)

		// patch name
		pname := models.MacbethColorCode(index).String()

		// create color code model
		colorcode := models.SetColorCode(index+1, pname, red8bit, green8bit, blue8bit, 255)

		// update
		cg.devColorCodes = append(cg.devColorCodes, *colorcode)
	}

	// check
	if len(cg.devColorCodes) == 0 {
		return false
	}

	return true
}

func (cg *colorChartController) runDevice(gamma float64, linearMat []float64) []models.ColorCode {

	// --- 1st stage ---
	// calculate device response
	cg.calculateDeviceResponse(gamma)

	// --- 2nd stage ---
	// calculate linear matrix
	if len(linearMat) == 0 {
		// use outside csv file
		cg.calculateLinearMatrix(cg.linearMatrixElm)
	} else {
		// use inputted linear matrix elements
		cg.calculateLinearMatrix(linearMat)
	}

	// --- 3rd stage ---
	// white balance
	cg.calculateWhiteBalance(refPathNoForWB)

	// --- 4th stage ---
	// gamma correction
	cg.calculateGammaCorrection(gamma)

	// --- 5th stage ---
	// degitize
	cg.convert8BitData()

	return cg.devColorCodes
}

/*
GenerateColorPatchImage
	in	:data []models.ColorCode, filesavepath, dirname string, width, height int
	out	:bool
*/
func (cg *colorChartController) GenerateColorPatchImage(data []models.ColorCode, filesavepath, dirname string, width, height int) bool {

	// initialize image controller
	imgcontroller := NewImageController()

	// initialize directory handler
	dirhandler := util.NewDirectoryHandler()
	dirhandler.MakeDirectory(filesavepath, dirname)

	// create path string for file save
	path := filesavepath + dirname + "/"

	// initalize IO handler
	iohandler := util.NewIOUtil()

	// generate color patches
	for _, rawdata := range data {
		rawimage := imgcontroller.CreateSolidImage(*rawdata.GenerateColorRGBA(), height, width)
		if !iohandler.StreamOutPNGFile(path, rawdata.GetName(), rawimage) {
			break
		}
	}

	// return
	return true
}

/*
Generate24MacbethColorChart
	in	:dev bool, filesavepath string
	out	:bool
*/
func (cg *colorChartController) Generate24MacbethColorChart(dev bool, filesavepath string) bool {
	imageContrller := NewImageController()
	var fileName string
	if dev {
		fileName = dev24ColorChartName
	} else {
		fileName = std24ColorChartName
	}

	// create 24 Macbeth Chart image
	if !imageContrller.Create24MacbethChart(filesavepath, fileName) {
		return false
	}

	return true
}

/*
RunStandad
	in	:filepath string
	out	:[]models.ColorCode
*/
func (cg *colorChartController) RunStandad(filepath string) []models.ColorCode {
	return models.ReadColorCode(filepath)
}

/*
SaveColorPatchCSVData
	in	:savepath, filenam string
	out	:bool
*/
func (cg *colorChartController) SaveColorPatchCSVData(data []models.ColorCode, savepath, filename string) bool {
	// data stcoker
	dataArray := make([][]string, 0)

	// serialize data
	for _, rawdata := range data {
		str := rawdata.SerializeData()
		dataArray = append(dataArray, str)
	}

	// save file as CSV file
	iohandler := util.NewIOUtil()
	if !iohandler.WriteCSVFile(savepath, filename, dataArray) {
		return false
	}

	// return
	return true
}
