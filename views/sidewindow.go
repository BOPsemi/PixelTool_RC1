package views

import (
	"PixelTool_RC1/models"
	"PixelTool_RC1/util"

	"github.com/asaskevich/EventBus"
	"github.com/therecipe/qt/widgets"
)

/*
Default values
*/
const (
	GAMMA       = 0.24
	StdSavePath = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"
	StdSaveDir  = "std_path"

	DevSavePath = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"
	DevSaveDir  = "dev_path"

	DeltaSavePath = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"
	DeltaSaveDir  = "delta_e"

	DevQEInputPath  = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/device_QE.csv"
	WPInputPath     = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/white_pixel.csv"
	LinmatInputPath = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/linearmatrix_elm.csv"
)

/*
--- Bus signal information ---
Tag				:sideWin:settingInfo
*/

/*
SideWindow :Side window
*/
type SideWindow struct {
	// --- Enviroment Setting group ---
	gammaAdjuster       *SliderInput         // slider object for gamma correction
	lightSourceSelector *ComboBoxSelector    // light source selector
	patchSizeInput      *PixelSizeInputField // patch size setting

	// --- file save group ---
	stdPatchSave *SavePathField // standard Macbeth Patch save point
	devPatchSave *SavePathField // Simulated Macbeth Pathc save point
	deltaESave   *SavePathField // deltaE save pint

	// --- input file group ---
	deviceQEData   *InputField // device QE file input
	whitePixelData *InputField // white pixel file input
	linearMatData  *InputField // linear matrix elem file input

	// Apply button
	applybutton   *widgets.QPushButton // apply button
	defaultButton *widgets.QPushButton // defalut setting loading button

	// Cell
	Cell *widgets.QWidget
}

/*
NewSideWindow :initializer of SideWindow
*/
func NewSideWindow(bus EventBus.Bus) *SideWindow {
	obj := new(SideWindow)

	// initialize widgets
	obj.Cell = widgets.NewQWidget(nil, 0)

	// initalize button
	obj.applybutton = widgets.NewQPushButton2("Run", obj.Cell)
	obj.defaultButton = widgets.NewQPushButton2("Default Setting", obj.Cell)

	// initialize each gourp
	envGroup := obj.setupEnvGroup()
	fileSaveGroup := obj.setFileSaveGroup()
	inputFileGroup := obj.setInputFileGroup()

	// resize
	envGroup.SetMaximumHeight(180)
	fileSaveGroup.SetMaximumHeight(250)
	inputFileGroup.SetMaximumHeight(200)

	// layout
	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(envGroup, 0, 0)
	layout.AddWidget(fileSaveGroup, 0, 0)
	layout.AddWidget(inputFileGroup, 0, 0)
	layout.AddWidget(obj.applybutton, 0, 0)
	layout.AddWidget(obj.defaultButton, 0, 0)
	layout.SetSpacing(1)
	layout.SetContentsMargins(0, 0, 0, 0)

	// apply layout
	obj.Cell.SetLayout(layout)

	// action connection
	obj.applybutton.ConnectClicked(func(checked bool) {
		info := new(models.SettingInfo)
		//
		info.Gamma = obj.gammaAdjuster.Value
		info.LightSource = obj.lightSourceSelector.SelectedItem

		// patch size
		info.PatchSize.H = obj.patchSizeInput.HorizontalSize
		info.PatchSize.V = obj.patchSizeInput.VerticalSize

		// field
		info.StdPatchSavePath = obj.stdPatchSave.textFieldForPath.Text()
		info.StdPatchSaveDirName = obj.stdPatchSave.textFieldForDirName.Text()
		info.DevPatchSavePath = obj.devPatchSave.textFieldForPath.Text()
		info.DevPatchSaveDirName = obj.devPatchSave.textFieldForDirName.Text()
		info.DeltaESavePath = obj.deltaESave.textFieldForPath.Text()
		info.DeltaESaveDirName = obj.deltaESave.textFieldForDirName.Text()
		info.DeiceQEDataPath = obj.deviceQEData.textField.Text()
		info.WhitePixelDataPath = obj.whitePixelData.textField.Text()
		info.LinearMatrixDataPath = obj.linearMatData.textField.Text()

		// validation
		validationStatus := true
		if !(obj.validation(info.StdPatchSavePath) && obj.validation(info.DevPatchSavePath) && obj.validation(info.DeltaESavePath)) {
			errorMessage := "We found some empty fields in Save Path category"
			bus.Publish("main:message", errorMessage)
			validationStatus = false
		}

		if !obj.validation(info.DeiceQEDataPath) {
			errorMessage := "Device QE data is missing"
			bus.Publish("main:message", errorMessage)
			validationStatus = false
		}

		if !obj.validation(info.WhitePixelDataPath) {
			errorMessage := "White Pixel data is missing"
			bus.Publish("main:message", errorMessage)
			validationStatus = false
		}

		if !obj.validation(info.LinearMatrixDataPath) {
			errorMessage := "Linear matrix elements data is missing"
			bus.Publish("main:message", errorMessage)
			validationStatus = false
		}

		if validationStatus {
			bus.Publish("sideWin:settingInfo", info)
			bus.Publish("main:message", "Setted file information")
		}
	})

	obj.defaultButton.ConnectClicked(func(checked bool) {
		obj.defaultSetting()
	})

	return obj
}

// Enviroment setting group
func (sw *SideWindow) setupEnvGroup() *widgets.QGroupBox {
	sw.gammaAdjuster = NewSliderInput("Gamma", 0.42)
	sw.lightSourceSelector = NewComboBoxSelector("Light Source", []string{"D65", "D50", "illA"})
	sw.patchSizeInput = NewPixelSizeInputField("Patch Size", 100, 100)

	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(sw.gammaAdjuster.Cell, 0, 0)
	layout.AddWidget(sw.lightSourceSelector.Cell, 0, 0)
	layout.AddWidget(sw.patchSizeInput.Cell, 0, 0)
	layout.SetSpacing(0)
	layout.SetContentsMargins(0, 0, 0, 0)

	group := widgets.NewQGroupBox2("Simulation Enviroment Setting", nil)
	group.SetLayout(layout)

	return group
}

// file save group
func (sw *SideWindow) setFileSaveGroup() *widgets.QGroupBox {
	sw.stdPatchSave = NewSavePathField("Std Patch Save")
	sw.devPatchSave = NewSavePathField("Dev Patch Save")
	sw.deltaESave = NewSavePathField("DeltaE Data Save")

	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(sw.stdPatchSave.Cell, 0, 0)
	layout.AddWidget(sw.devPatchSave.Cell, 0, 0)
	layout.AddWidget(sw.deltaESave.Cell, 0, 0)
	layout.SetSpacing(0)
	layout.SetContentsMargins(0, 0, 0, 0)

	group := widgets.NewQGroupBox2("File Save Setting", nil)
	group.SetLayout(layout)

	return group
}

// input file group
func (sw *SideWindow) setInputFileGroup() *widgets.QGroupBox {
	sw.deviceQEData = NewInputField("Device QE", "Device QE raw data")
	sw.whitePixelData = NewInputField("White Pixel", "White pixel raw data")
	sw.linearMatData = NewInputField("Linear Matrix", "Linear Matrix element data")

	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(sw.deviceQEData.Cell, 0, 0)
	layout.AddWidget(sw.whitePixelData.Cell, 0, 0)
	layout.AddWidget(sw.linearMatData.Cell, 0, 0)
	layout.SetSpacing(0)
	layout.SetContentsMargins(0, 0, 0, 0)

	group := widgets.NewQGroupBox2("Input File information", nil)
	group.SetLayout(layout)

	return group
}

// func validation
func (sw *SideWindow) validation(str string) bool {
	if str == "" {
		return false
	}
	return true
}

// default setting
func (sw *SideWindow) defaultSetting() {
	dirHandler := util.NewDirectoryHandler()
	currentPath := dirHandler.GetCurrentDirectoryPath()

	// set default value
	sw.gammaAdjuster.slider.SetValue(45)
	sw.gammaAdjuster.textField.SetText("0.45")

	// -- standard patch
	sw.stdPatchSave.textFieldForPath.SetText(currentPath + "/data/")
	sw.stdPatchSave.textFieldForDirName.SetText("std_patch")

	// -- device patch
	sw.devPatchSave.textFieldForPath.SetText(currentPath + "/data/")
	sw.devPatchSave.textFieldForDirName.SetText("dev_patch")

	// -- delta E
	sw.deltaESave.textFieldForPath.SetText(currentPath + "/data/")
	sw.deltaESave.textFieldForDirName.SetText("deltaE_result")

	// -- Device QE
	sw.deviceQEData.textField.SetText(currentPath + "/data/" + "device_QE.csv")

	// -- white pixel
	sw.whitePixelData.textField.SetText(currentPath + "/data/" + "white_pixel.csv")

	// -- Linear Matrix data
	sw.linearMatData.textField.SetText(currentPath + "/data/" + "linearmatrix_elm.csv")

}
