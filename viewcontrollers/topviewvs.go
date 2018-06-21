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
}

// topViewViewController
type topViewViewController struct {
}

/*
NewTopViewViewController :initializer
*/
func NewTopViewViewController() TopViewViewController {
	obj := new(topViewViewController)

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
	csvFilePath := dirHandler.GetCurrentDirectoryPath() + "/data/" + "Macbeth_Patch_Code.csv"

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
		filename := "std_24_ColorCahrt"
		imageController := controllers.NewImageController()

		// 24 color chart
		imageController.Create24MacbethChart(path, filename)

		// update status
		status = true
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

	const (
		start         int   = 400
		stop          int   = 700
		step          int   = 5
		refPatchNo    int   = 19
		refPatchLevel uint8 = 243
	)

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
			redGain, blueGain := devChartVC.CalculateWhiteBalanceGain(refPatchNo)
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
			filename := "dev_24_ColorCahrt"
			imageController := controllers.NewImageController()

			// 24 color chart
			imageController.Create24MacbethChart(path, filename)

			// update status
			status = true
		}
	}
	return status
}
