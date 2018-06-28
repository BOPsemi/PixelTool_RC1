package views

import (
	"PixelTool_RC1/models"
	"PixelTool_RC1/viewcontrollers"
	"strconv"
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

	// view controller
	viewController viewcontrollers.TopViewViewController

	// state indicator
	stdChartReady bool
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

	// initialize view controller
	obj.viewController = viewcontrollers.NewTopViewViewController()

	// resize
	obj.sideWin.Cell.SetMaximumWidth(460)

	// layout
	layout := widgets.NewQHBoxLayout()
	layout.AddWidget(obj.sideWin.Cell, 0, 0)
	layout.AddWidget(obj.mainWin.Cell, 0, 0)

	// apply layout
	obj.Cell.SetLayout(layout)

	// event bus subscribe
	bus.Subscribe("sideWin:settingInfo", obj.settingInfoReciever)
	bus.Subscribe("main:calculateDeltaE", obj.calculateDeltaE)
	bus.Subscribe("main:optimizeLinearMat", obj.optimizeLinearMat)

	bus.Subscribe("main:message", obj.messageReciever)

	// state indicator
	obj.stdChartReady = false

	return obj
}

// --- Subscriber ---
// Setting Infor reciever
func (tw *TopWindow) settingInfoReciever(info *models.SettingInfo) {
	// initalize view controller
	vc := viewcontrollers.NewTopViewViewController()

	// make Standard Macbeth patch
	if !tw.stdChartReady {
		if vc.GenerateMacbethColorChart(false, info) {
			tw.mainWin.messageBox.Append("Successed to generate standard Macbeth color patch images" + "  :  " + time.Now().Format(time.ANSIC))
			tw.stdChartReady = true
		} else {
			tw.mainWin.messageBox.Append("Faild to generate standard Macbeth color patch images" + "  :  " + time.Now().Format(time.ANSIC))
		}
	}

	// make Device Macbeth patch
	if vc.GenerateMacbethColorChart(true, info) {
		tw.mainWin.messageBox.Append("Successed to generate device Macbeth color patch images" + "  :  " + time.Now().Format(time.ANSIC))
	} else {
		tw.mainWin.messageBox.Append("Faild to generate device Macbeth color patch images" + "  :  " + time.Now().Format(time.ANSIC))
	}
}

// calculateDeltaE
func (tw *TopWindow) calculateDeltaE(info *models.SettingInfo) {
	if info != nil {
		stdDataPath := info.StdPatchSavePath + info.StdPatchSaveDirName + "/" + std24ColorChartName + ".csv"
		devDataPath := info.DevPatchSavePath + info.DevPatchSaveDirName + "/" + dev24ColorChartName + ".csv"

		// kvalue definition
		kvalues := []float64{1.0, 1.0, 1.0}

		if results, ok := tw.viewController.EvaluateDeltaE(stdDataPath, devDataPath, kvalues); ok {
			if tw.viewController.SaveDeltaEResultData() {

				// output the result to message box
				tw.mainWin.messageBox.Append("Delta-E Calculation result")
				for index, data := range results {
					str := strconv.Itoa(index+1) + " : " + strconv.FormatFloat(data, 'f', 4, 64)
					tw.mainWin.messageBox.Append(str)
				}
			}
		}
	}
}

// optimizeLinearMat
func (tw *TopWindow) optimizeLinearMat(info *models.SettingInfo) {

}

// message reciever
func (tw *TopWindow) messageReciever(message string) {
	tw.mainWin.messageBox.Append(message + "  :  " + time.Now().Format(time.ANSIC))
	tw.mainWin.messageBox.Repaint()
}

// ---
