package views

import (
	"PixelTool_RC1/models"
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
	bus.Subscribe("main:message", obj.messageReciever)

	// state indicator
	obj.stdChartReady = false

	return obj
}

// --- Subscriber ---
func (tw *TopWindow) settingInfoReciever(info *models.SettingInfo) {
	// initalize view controller
	vc := viewcontrollers.NewTopViewViewController()

	// generate standard Macbeth color charts
	if !tw.stdChartReady {
		if vc.GenerateStdMacbethColorChart(info) {
			tw.mainWin.messageBox.Append("Successed to generate standard Macbeth color pathc images" + "  :  " + time.Now().Format(time.ANSIC))
			tw.stdChartReady = true
		} else {
			tw.mainWin.messageBox.Append("Faild to generate standard Macbeth color pathc images" + "  :  " + time.Now().Format(time.ANSIC))
		}
	}

	// generate device Macbeth color charts
	if vc.GenerateDevMacbethColorChart(info) {
		tw.mainWin.messageBox.Append("Successed to generate device Macbeth color pathc images" + "  :  " + time.Now().Format(time.ANSIC))
	} else {
		tw.mainWin.messageBox.Append("Faild to generate device Macbeth color pathc images" + "  :  " + time.Now().Format(time.ANSIC))
	}
}

// message reciever
func (tw *TopWindow) messageReciever(message string) {
	tw.mainWin.messageBox.Append(message + "  :  " + time.Now().Format(time.ANSIC))
	tw.mainWin.messageBox.Repaint()
}

// ---
