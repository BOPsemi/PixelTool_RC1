package viewcontrollers

import (
	"PixelTool_RC1/controllers"
	"PixelTool_RC1/models"
	"PixelTool_RC1/util"
	"strconv"
)

/*
DeviceResponseViewController :device response view controller
*/
type DeviceResponseViewController interface {
	SetupReadingFile(info *models.SettingInfo) map[string]string
	ReadLinearMatrixElmData(filepath string) []float64
	ReadResponseRawData(filepath map[string]string) bool

	CalculateDeviceResponse(ill models.IlluminationCode, start, stop, step int, gammma float64, refPatchNum int) bool
	CalculateLinearMatrix(elm []float64) bool
	CalculateWhiteBalanceGain(refPatchNumber int) (redGain, blueGain float64)
	Calculate8bitResponse(patchNumber int, data []float64, redGain, blueGain float64, refLevel uint8) *models.ColorCode

	// getters
	RawData() []models.ChannelResponse
	RawResponseData() []models.ChannelResponse
	LinearizedResponseData() [][]float64

	// steam out PNG patch image
	CreateColorCodePatch(data *models.ColorCode, filesavepath, dirname string, width, height int) bool

	// save the data as CSV file
	SaveColorCodePatchData(savepath, filename string) bool
}

// defintion of structure
type deviceResponseViewController struct {
	resCon controllers.ResponseController // response controller

	// stockers
	rawData         []models.ChannelResponse
	rawResponseData []models.ChannelResponse

	linearizedResData [][]float64

	colorCodes []models.ColorCode
}

/*
NewDeviceResponseViewController : initializer of VC
*/
func NewDeviceResponseViewController() DeviceResponseViewController {
	obj := new(deviceResponseViewController)

	// initialize properties
	obj.resCon = controllers.NewResponseController()

	// initialize stockers
	obj.rawData = make([]models.ChannelResponse, 0)
	obj.rawResponseData = make([]models.ChannelResponse, 0)
	obj.linearizedResData = make([][]float64, 0)
	obj.colorCodes = make([]models.ColorCode, 0)

	return obj
}

/*
SetupReadingFile
	in	;info *models.SettingInfo
	out	;(map[string]string, bool
*/
func (vc *deviceResponseViewController) SetupReadingFile(info *models.SettingInfo) map[string]string {
	filepath := make(map[string]string, 0)

	// get current path
	dirHandler := util.NewDirectoryHandler()
	path := dirHandler.GetCurrentDirectoryPath() + "/data/"

	// illumination data path
	filepath["D65"] = path + illuminationD65FileName
	filepath["IllA"] = path + illuminationIllAFileName

	// raw data
	filepath["DeviceQE"] = info.DeiceQEDataPath
	filepath["ColorChecker"] = path + macbethColorCheckerFileName

	return filepath
}

/*
ReadLinearMatrixElmData
	in	;filepath string
	out	;[]float64
*/
func (vc *deviceResponseViewController) ReadLinearMatrixElmData(filepath string) []float64 {
	// read linear matrix element data
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

/*
ReadResponseRawData	:
	in	;filepath map[string]string
	out	;bool
*/
func (vc *deviceResponseViewController) ReadResponseRawData(filepath map[string]string) bool {
	status := false

	// read response raw data
	if len(filepath) != 0 {
		status = vc.resCon.ReadResponseData(filepath)
	}

	return status
}

/*
CalculateDeviceResponse
	in	;
		ill models.Illumination,
		start 	;scan start wavelength
		stop	;scan stop wavelength
		step	;scan step wavelength
		refPatchNum
		refPatchLevel
	out	;bool
*/
func (vc *deviceResponseViewController) CalculateDeviceResponse(ill models.IlluminationCode, start, stop, step int, gamma float64, refPatchNum int) bool {
	status := false

	/*
		start to calculate channel response
			1. calculate channel response
			2. calculate gamma correction
			3. check slice size, and then update result to stocker
	*/
	if ok, responses := vc.resCon.CalculateChannelResponse(ill, start, stop, step, refPatchNum); ok {
		// stock rawdata
		vc.rawData = responses

		// initialize gamma correction buffer
		gammaCorrectedRes := make([]models.ChannelResponse, 0)

		// calculate gamma correction
		for _, data := range responses {
			flag, result := vc.resCon.CalculateGammaCorrection(gamma, &data)
			if flag {
				gammaCorrectedRes = append(gammaCorrectedRes, *result)
			}
		}

		// check slice size
		if len(gammaCorrectedRes) != 0 {
			vc.rawResponseData = gammaCorrectedRes

			//update status
			status = true
		}
	}

	return status
}

/*
CalculateLinearMatrix
	in	;elm []float64, grgbrb []float64
	out	;bool
*/
func (vc *deviceResponseViewController) CalculateLinearMatrix(elm []float64) bool {
	status := false

	if len(vc.rawResponseData) != 0 {
		// buffer
		responses := make([][]float64, 0)
		for _, data := range vc.rawResponseData {
			// change data format to linear matrix calculation
			grgbrb := []float64{
				data.Gr,
				data.Gb,
				data.R,
				data.B,
			}

			// calcualte linear matrix
			response := vc.resCon.CalculateLinearMatrix(elm, grgbrb)

			// stock
			responses = append(responses, response)
		}

		// check response
		if len(responses) != 0 {
			vc.linearizedResData = responses

			// update status
			status = true
		}
	}

	return status
}

/*
CalculateWhiteBalanceGain
	in	;refPtchNumber int
	out	;redGain, blueGain float64
*/
func (vc *deviceResponseViewController) CalculateWhiteBalanceGain(refPatchNumber int) (redGain, blueGain float64) {
	if refPatchNumber > -1 && refPatchNumber < 25 {
		return vc.resCon.CalculateWhiteBalanceGain(vc.linearizedResData[refPatchNumber-1])
	}
	return 0.0, 0.0
}

/*
Calculate8bitResponse
	in	;data []float64, redGain, blueGain float64, refLevel uint8
	out	;models.ColorCode
*/
func (vc *deviceResponseViewController) Calculate8bitResponse(patchNumber int, data []float64, redGain, blueGain float64, refLevel uint8) *models.ColorCode {

	// calculate raw data
	red := data[0] * redGain
	green := data[1]
	blue := data[2] * blueGain

	// digitize signal
	digitizer := util.NewDigitizer()
	red8bit := digitizer.D8bitDigitizeData(red, refLevel)
	green8bit := digitizer.D8bitDigitizeData(green, refLevel)
	blue8bit := digitizer.D8bitDigitizeData(blue, refLevel)

	// patch name
	pname := models.MacbethColorCode(patchNumber).String()

	// create color code model
	colorcode := models.SetColorCode(patchNumber+1, pname, red8bit, green8bit, blue8bit, 255)

	// stack the data into stocker
	vc.colorCodes = append(vc.colorCodes, *colorcode)

	return colorcode

}

/*
CreateColorCodePatch
	in	;filesavepath, dirname string, width, height int
	out	;bool
*/
func (vc *deviceResponseViewController) CreateColorCodePatch(data *models.ColorCode, filesavepath, dirname string, width, height int) bool {
	status := false

	if filesavepath != "" && dirname != "" && data != nil {

		// create file save path
		path := filesavepath + dirname + "/"

		// create solid image from data
		imgcontroller := controllers.NewImageController()
		rawimage := imgcontroller.CreateSolidImage(*data.GenerateColorRGBA(), height, width)

		dirhandler := util.NewDirectoryHandler()
		if dirhandler.MakeDirectory(filesavepath, dirname) {

			// --- Not exist save folder ---
			// save image file
			iohandler := util.NewIOUtil()
			if iohandler.StreamOutPNGFile(path, data.GetName(), rawimage) {
				status = true
			}
		} else {

			// --- exist save folder ---
			// save image file
			iohandler := util.NewIOUtil()
			if iohandler.StreamOutPNGFile(path, data.GetName(), rawimage) {
				status = true
			}
		}

	}

	return status
}

/*
SaveColorCodePatchData
	in	;savepath, filename string
	out	;bool
*/
func (vc *deviceResponseViewController) SaveColorCodePatchData(savepath, filename string) bool {
	status := false

	if len(vc.colorCodes) != 0 {
		if savepath != "" && filename != "" {
			// make string data from property
			data := make([][]string, 0)
			for _, obj := range vc.colorCodes {
				dataString := obj.SerializeData()
				data = append(data, dataString)
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

/*
RawData
	out	;[]models.ChannelResponse
*/
func (vc *deviceResponseViewController) RawData() []models.ChannelResponse {
	return vc.rawData
}

/*
RawResponseData
	out	;[]models.ChannelResponse
*/
func (vc *deviceResponseViewController) RawResponseData() []models.ChannelResponse {
	return vc.rawResponseData
}

/*
LinearizedResponseData
	out ;[][]float64
*/
func (vc *deviceResponseViewController) LinearizedResponseData() [][]float64 {
	return vc.linearizedResData
}
