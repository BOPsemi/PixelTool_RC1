package viewcontrollers

import (
	"PixelTool_RC1/controllers"
	"PixelTool_RC1/models"
	"PixelTool_RC1/util"
)

/*
TopViewViewController :
*/
type TopViewViewController interface {
	GenerateStdMacbethColorChart(info *models.SettingInfo) bool
	GenerateDevMacbethColorChart(info *models.SettingInfo) bool

	EvaluateDeltaE(refDataPath, compDataPath string, kvalues []float64) ([]float64, bool)
	SaveDeltaEResultData() bool
}

// topViewViewController
type topViewViewController struct {
	deltaEval controllers.DeltaEvaluationController
}

/*
NewTopViewViewController :initializer
*/
func NewTopViewViewController() TopViewViewController {
	obj := new(topViewViewController)

	obj.deltaEval = controllers.NewDeltaEvaluationController()

	return obj
}

/*
GenerateStdMacbethColorChart
	in	;info *models.SettingInfo
	out	;bool
*/
func (vc *topViewViewController) GenerateStdMacbethColorChart(info *models.SettingInfo) bool {
	status := false

	// directory handler
	dirHandler := util.NewDirectoryHandler()
	csvFilePath := dirHandler.GetCurrentDirectoryPath() + "/data/" + macbethColorChartCodeFileName

	// standard Macbeth color chart generate
	stdChartVC := NewColorCheckerViewController()
	if stdChartVC.CreateColorCodePatch(
		csvFilePath,
		info.StdPatchSavePath,
		info.StdPatchSaveDirName,
		info.PatchSize.H,
		info.PatchSize.V,
	) {
		// -- create 24 chart
		// path setting
		path := info.StdPatchSavePath + info.StdPatchSaveDirName + "/"
		imageController := controllers.NewImageController()

		// 24 color chart
		if imageController.Create24MacbethChart(path, std24ColorChartName) {
			// save csv file
			if stdChartVC.SaveColorCodePatchData(path, std24ColorChartName) {
				// update status
				status = true
			}
		}
	}

	return status
}

/*
GenerateDevMacbethColorChart
	in	:info *models.SettingInfo
	out	:bool
*/
func (vc *topViewViewController) GenerateDevMacbethColorChart(info *models.SettingInfo) bool {
	status := false

	devChartVC := NewDeviceResponseViewController()

	// read response Raw data
	if devChartVC.ReadResponseRawData(devChartVC.SetupReadingFile(info)) {

		// calculate device QE chart
		switch info.LightSource {
		case "D65":
			devChartVC.CalculateDeviceResponse(models.D65, start, stop, step, info.Gamma, refPatchNo)
		case "illA":
			devChartVC.CalculateDeviceResponse(models.IllA, start, stop, step, info.Gamma, refPatchNo)
		}

		// linear matrix calculation
		if devChartVC.CalculateLinearMatrix(devChartVC.ReadLinearMatrixElmData(info.LinearMatrixDataPath)) {
			// calculate red and blue gain for white balance
			redGain, blueGain := devChartVC.CalculateWhiteBalanceGain(refPatchNoForWB)
			for index, data := range devChartVC.LinearizedResponseData() {

				// generate patcheds
				code := devChartVC.Calculate8bitResponse(index, data, redGain, blueGain, refPatchLevel)
				devChartVC.CreateColorCodePatch(
					code,
					info.DevPatchSavePath,
					info.DevPatchSaveDirName,
					info.PatchSize.H,
					info.PatchSize.V,
				)
			}

			// create 24 patch
			path := info.DevPatchSavePath + info.DevPatchSaveDirName + "/"
			imageController := controllers.NewImageController()

			// 24 color chart
			if imageController.Create24MacbethChart(path, dev24ColorChartName) {
				if devChartVC.SaveColorCodePatchData(path, dev24ColorChartName) {
					// update status
					status = true
				}
			}
		}
	}
	return status
}

/*
EvaluateDeltaE
	in	:refDataPath, compDataPath string, kvalues []float64
	out	:bool
*/
func (vc *topViewViewController) EvaluateDeltaE(refDataPath, compDataPath string, kvalues []float64) ([]float64, bool) {
	//deltaEResults := make([]float64, 0)

	if results, status := vc.deltaEval.EvaluateDeltaE(refDataPath, compDataPath, kvalues); status {
		return results, status
	}

	return []float64{}, false
}

/*
SaveDeltaEResultData
	in	:savepath, filename string
	out	:bool
*/
func (vc *topViewViewController) SaveDeltaEResultData() bool {
	// directory handler
	dirHandler := util.NewDirectoryHandler()
	savepath := dirHandler.GetCurrentDirectoryPath() + "/data/"

	if vc.deltaEval.SaveDeltaEResultData(savepath, deltaEReulstFileName) {
		return true
	}

	return false
}
