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
	GenerateMacbethColorChart(dev bool, info *models.SettingInfo) bool

	EvaluateDeltaE(colorSpace models.ColorSpace, refDataPath, compDataPath string, kvalues []float64) ([]float64, bool)
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
EvaluateDeltaE
	in	:refDataPath, compDataPath string, kvalues []float64
	out	:bool
*/
func (vc *topViewViewController) EvaluateDeltaE(colorSpace models.ColorSpace, refDataPath, compDataPath string, kvalues []float64) ([]float64, bool) {
	// --- Stage 1 ---
	// set Data
	vc.deltaEval.SetData(refDataPath, compDataPath)

	// --- Stage 2 ---
	// calculate
	deltaE, _ := vc.deltaEval.RunDeltaEEvaluation(colorSpace, [][]float64{}, [][]float64{}, kvalues)

	// --- Stage 3 ---
	// check
	if len(deltaE) == 0 {
		return []float64{}, false
	}

	// return
	return deltaE, true
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

/*
GenerateMacbethColorChart
	in	:dev bool, info *models.SettingInfo
	out	:bool
*/
func (vc *topViewViewController) GenerateMacbethColorChart(dev bool, info *models.SettingInfo) bool {
	// initalize color chart controller
	cccontroller := controllers.NewColorChartController()

	// initalize dir handler
	dirhandler := util.NewDirectoryHandler()

	// --- make csv file path
	path := dirhandler.GetCurrentDirectoryPath()
	datapath := path + "/data/"
	filepath := datapath + macbethColorChartCodeFileName

	if dev {
		// device

		// --- Light source identification ---
		var lightSource models.IlluminationCode
		switch info.LightSource {
		case "D65":
			lightSource = models.D65
		case "IllA":
			lightSource = models.IllA
		default:
			lightSource = models.D65
		}

		// --- calculte response ---
		data := cccontroller.RunDevice(
			info.LinearMatrixDataPath,
			datapath,
			info.DeiceQEDataPath,
			lightSource,
			info.Gamma,
			[]float64{},
		)

		// --- generate color path ---
		if !cccontroller.GenerateColorPatchImage(
			data,
			info.DevPatchSavePath,
			info.DevPatchSaveDirName,
			info.PatchSize.H,
			info.PatchSize.V,
		) {
			return false
		}

		// --- generate 24 Macbeth patch chart ----
		if !cccontroller.Generate24MacbethColorChart(
			true,
			info.DevPatchSavePath+info.DevPatchSaveDirName+"/",
		) {
			return false
		}

		// --- save csv file ---

		if !cccontroller.SaveColorPatchCSVData(
			data,
			info.DevPatchSavePath+info.DevPatchSaveDirName+"/",
			dev24ColorChartName,
		) {
			return false
		}

	} else {
		// standard

		// --- get color code ---
		data := cccontroller.RunStandad(filepath)

		// --- generate color path ---
		if !cccontroller.GenerateColorPatchImage(
			data,
			info.StdPatchSavePath,
			info.StdPatchSaveDirName,
			info.PatchSize.H,
			info.PatchSize.V,
		) {
			return false
		}

		// --- generate 24 Macbeth patch chart ----
		if !cccontroller.Generate24MacbethColorChart(
			false,
			info.StdPatchSavePath+info.StdPatchSaveDirName+"/",
		) {
			return false
		}

		// --- save csv file ---
		if !cccontroller.SaveColorPatchCSVData(
			data,
			info.StdPatchSavePath+info.StdPatchSaveDirName+"/",
			std24ColorChartName,
		) {
			return false
		}
	}

	return true
}
