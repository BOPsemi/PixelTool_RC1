package views

import (
	"PixelTool_RC1/models"
	"fmt"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/therecipe/qt/widgets"
)

const (
	nodataImagePath     = "./data/NoDataImage.png"
	std24ColorChartName = "std_24_ColorChart"
	dev24ColorChartName = "dev_24_ColorChart"
)

/*
ImageViewIdentifier : vierwer identifier
*/
type ImageViewIdentifier int

/*
StdImageViewer :For Standard Macbeth chart imageViewer
*/
const (
	StdImageViewer ImageViewIdentifier = iota // StdImageViewer
	DevImageViewer
)

/*
MainWindow :main window
*/
type MainWindow struct {

	// --- Image Viewer ---
	stdCCImageView *ImageViewer // standard Macbeth color chart
	devCCImageView *ImageViewer // device Macbeth color chart

	// --- load buttons ---
	stdImageLoadButton *widgets.QPushButton // image load button for standard
	devImageLoadButton *widgets.QPushButton // image load button for device

	// --- additional action buttons ---
	calcDeltaEButton *widgets.QPushButton // reload linear matrix element
	optimizeButton   *widgets.QPushButton // show deltaE
	saveLogButton    *widgets.QPushButton // save log button

	// --- message box ---
	messageBox *widgets.QTextEdit // message box

	// --- current path ---
	settingInfo *models.SettingInfo

	// widget
	Cell *widgets.QWidget
}

/*
NewMainWindow : initializer of main window
*/
func NewMainWindow(bus EventBus.Bus) *MainWindow {
	obj := new(MainWindow)

	obj.Cell = widgets.NewQWidget(nil, 0)

	// imageViewer initialize
	obj.stdCCImageView = NewImageViewer(nodataImagePath, 0.5)
	obj.devCCImageView = NewImageViewer(nodataImagePath, 0.5)

	// button
	obj.stdImageLoadButton = widgets.NewQPushButton2("Image Load", obj.Cell)
	obj.devImageLoadButton = widgets.NewQPushButton2("Image Load", obj.Cell)
	obj.calcDeltaEButton = widgets.NewQPushButton2("Calculate DeltaE", obj.Cell)
	obj.optimizeButton = widgets.NewQPushButton2("Optimize Linear Matrix", obj.Cell)
	obj.saveLogButton = widgets.NewQPushButton2("Save Log", obj.Cell)

	// message box setup
	initlog := "Logging started" + "  :  " + time.Now().Format(time.ANSIC)
	obj.messageBox = widgets.NewQTextEdit2(initlog, obj.Cell)

	// group
	stdGroup := obj.setupStdGroup()
	devGroup := obj.setupDevGroup()
	optGroup := obj.setupOptGroup()

	// layout
	layout := widgets.NewQGridLayout2()
	layout.AddWidget(stdGroup, 0, 0, 0)
	layout.AddWidget(devGroup, 0, 1, 0)
	layout.AddWidget3(obj.messageBox, 2, 0, 1, 2, 0)
	layout.AddWidget3(optGroup, 1, 0, 1, 2, 0)

	// apply layout
	obj.Cell.SetLayout(layout)

	// bus subsriber
	bus.Subscribe("sideWin:settingInfo", obj.settingInfoReciever)

	// action connection
	// --- image load buttons ---

	// for standard image
	obj.stdImageLoadButton.ConnectClicked(func(checked bool) {
		if obj.settingInfo != nil {
			// file path
			filepath := obj.settingInfo.StdPatchSavePath + obj.settingInfo.StdPatchSaveDirName + "/" + std24ColorChartName

			// load image
			obj.reloadImage(filepath, 0.5, StdImageViewer)
			bus.Publish("main:message", "Standard Macbeth Color Chart was reloded")

		} else {
			bus.Publish("main:message", "Image load error")
		}

	})

	// for device image
	obj.devImageLoadButton.ConnectClicked(func(checked bool) {
		if obj.settingInfo != nil {
			// file path
			filepath := obj.settingInfo.DevPatchSavePath + obj.settingInfo.DevPatchSaveDirName + "/" + dev24ColorChartName

			// load image
			obj.reloadImage(filepath, 0.5, DevImageViewer)
			bus.Publish("main:message", "Device Macbeth Color Chart was reloded")
		} else {
			bus.Publish("main:message", "Image load error")
		}

	})

	// delta-E calculation
	obj.calcDeltaEButton.ConnectClicked(func(checked bool) {
		if obj.settingInfo != nil {
			bus.Publish("main:calculateDeltaE", obj.settingInfo)
		} else {
			bus.Publish("main:message", "No target data")
		}
	})

	// linear matrix optimization
	obj.optimizeButton.ConnectClicked(func(checked bool) {
		if obj.settingInfo != nil {
			bus.Publish("main:optimizeLinearMat", obj.settingInfo)
		} else {
			bus.Publish("main:message", "No target data")
		}

	})

	// save
	obj.saveLogButton.ConnectClicked(func(checked bool) {
		bus.Publish("main:message", "Log Save")

		pathInputDialog := NewTextInputDialog("Log Save", "Path")
		pathInputDialog.Cell.Show()
		pathInputDialog.Cell.ConnectAccepted(func() {
			log := obj.messageBox.ToPlainText()
			fmt.Println(log)
		})
	})

	return obj
}

func (mm *MainWindow) settingInfoReciever(info *models.SettingInfo) {
	mm.settingInfo = info
}

// std group setting
func (mm *MainWindow) setupStdGroup() *widgets.QGroupBox {
	stdGroup := widgets.NewQGroupBox2("Standard Macbeth Color Chart", mm.Cell)
	stdLayout := widgets.NewQVBoxLayout()
	stdLayout.AddWidget(mm.stdCCImageView.Cell, 0, 0)
	stdLayout.AddWidget(mm.stdImageLoadButton, 0, 0)
	stdGroup.SetLayout(stdLayout)

	return stdGroup
}

// dev group
func (mm *MainWindow) setupDevGroup() *widgets.QGroupBox {
	devGroup := widgets.NewQGroupBox2("Device Macbeth Color Chart", mm.Cell)
	devLayout := widgets.NewQVBoxLayout()
	devLayout.AddWidget(mm.devCCImageView.Cell, 0, 0)
	devLayout.AddWidget(mm.devImageLoadButton, 0, 0)
	devGroup.SetLayout(devLayout)

	return devGroup
}

// opt group
func (mm *MainWindow) setupOptGroup() *widgets.QGroupBox {
	optGroup := widgets.NewQGroupBox(mm.Cell)
	optLayout := widgets.NewQHBoxLayout()
	optLayout.AddWidget(mm.calcDeltaEButton, 0, 0)
	optLayout.AddWidget(mm.optimizeButton, 0, 0)
	optLayout.AddWidget(mm.saveLogButton, 0, 0)
	optGroup.SetLayout(optLayout)

	return optGroup
}

// reloadImage
func (mm *MainWindow) reloadImage(path string, scale float64, identifier ImageViewIdentifier) {
	switch identifier {
	case StdImageViewer:
		mm.stdCCImageView.SetImageView(path, scale)
		mm.stdCCImageView.Cell.Repaint()

	case DevImageViewer:
		mm.devCCImageView.SetImageView(path, scale)
		mm.devCCImageView.Cell.Repaint()
	}
}
