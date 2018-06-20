package views

import (
	"PixelTool_RC1/controllers"
	"PixelTool_RC1/models"
	"PixelTool_RC1/util"
	"PixelTool_RC1/viewcontrollers"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/therecipe/qt/widgets"
)

/*
sideWin:settingInfo
main:message
*/

/*
TopWindow :top window structure
*/
type TopWindow struct {
	sideWin *SideWindow // side window
	mainWin *MainWindow // main window

	eventBus EventBus.Bus // Notification

	Cell *widgets.QWidget
}

/*
NewTopWindow :initializer of top window
*/
func NewTopWindow(bus EventBus.Bus) *TopWindow {
	obj := new(TopWindow)

	obj.Cell = widgets.NewQWidget(nil, 0)
	obj.eventBus = bus

	// initialize both windows
	obj.sideWin = NewSideWindow(bus)
	obj.mainWin = NewMainWindow(bus)

	// resize
	//obj.sideWin.Cell.SetMaximumWidth(460)

	// layout
	layout := widgets.NewQHBoxLayout()
	layout.AddWidget(obj.sideWin.Cell, 0, 0)
	layout.AddWidget(obj.mainWin.Cell, 0, 0)

	// apply layout
	obj.Cell.SetLayout(layout)

	// event bus subscribe
	bus.Subscribe("sideWin:settingInfo", obj.settingInfoReciever)
	bus.Subscribe("main:message", obj.messageReciever)

	return obj
}

// --- Subscriber ---
func (tw *TopWindow) settingInfoReciever(info *models.SettingInfo) {
	// generate standard macbeth color charts
	tw.generateStdMacbethColorChart(info)

	// generate device macbeth color charts
	tw.generateDevMacbethColorChart(info)
}

/*
Genrate standard macbeth color chart
*/
func (tw *TopWindow) generateStdMacbethColorChart(info *models.SettingInfo) bool {
	status := false

	dirHandler := util.NewDirectoryHandler()
	csvFilePath := dirHandler.GetCurrentDirectoryPath() + "/data/" + "Macbeth_Patch_Code.csv"

	// standard Macbeth color chart generate
	stdChartVC := viewcontrollers.NewColorCheckerViewController()
	state := stdChartVC.CreateColorCodePatch(
		csvFilePath,
		info.StdPatchSavePath,
		info.StdPatchSaveDirName,
		info.PatchSize.H,
		info.PatchSize.V,
	)
	if state {

		// -- create 24 chart
		// path setting
		path := info.StdPatchSavePath + info.StdPatchSaveDirName + "/"
		filename := "std_24_ColorCahrt"
		imageController := controllers.NewImageController()

		// 24 color chart
		imageController.Create24MacbethChart(path, filename)
		tw.mainWin.messageBox.Append("Successed to generate standard Macbeth color pathc images" + "  :  " + time.Now().Format(time.ANSIC))

		// update status
		status = true
	} else {
		tw.mainWin.messageBox.Append("Faild to generate standard Macbeth color pathc images" + "  :  " + time.Now().Format(time.ANSIC))
	}

	return status
}

/*
Generate device macbeth chart
*/
func (tw *TopWindow) generateDevMacbethColorChart(info *models.SettingInfo) bool {
	status := false

	const (
		start         int   = 400
		stop          int   = 700
		step          int   = 5
		refPatchNo    int   = 19
		refPatchLevel uint8 = 243
	)

	devChartVC := viewcontrollers.NewDeviceResponseViewController()

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

// message reciever
func (tw *TopWindow) messageReciever(message string) {
	tw.mainWin.messageBox.Append(message + "  :  " + time.Now().Format(time.ANSIC))
	tw.mainWin.messageBox.Repaint()
}

// ---
