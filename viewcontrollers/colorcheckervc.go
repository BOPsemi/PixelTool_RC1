package viewcontrollers

import (
	"PixelTool_RC1/controllers"
	"PixelTool_RC1/models"
	"PixelTool_RC1/util"
)

/*
ColorCheckerViewController :control view
	- generate color checker patches
*/
type ColorCheckerViewController interface {
	CreateColorCodePatch(csvfilepath, filesavepath, dirname string, width, height int) bool
	SaveColorCodePatchData(savepath, filename string) bool
}

// strcture definition
type colorCheckerViewController struct {
	imgcontroller controllers.ImageController
	dirhandler    util.DirectoryHandler
	iohandler     util.IOUtil

	// properties
	colorCodes []models.ColorCode
}

/*
NewColorCheckerViewController : initializer
*/
func NewColorCheckerViewController() ColorCheckerViewController {
	obj := new(colorCheckerViewController)

	// initialize instances
	obj.imgcontroller = controllers.NewImageController()
	obj.dirhandler = util.NewDirectoryHandler()
	obj.iohandler = util.NewIOUtil()

	return obj
}

/*
CreateColorCodePatch(csvfilepath, filesavepath, dirname string) bool
*/
func (cc *colorCheckerViewController) CreateColorCodePatch(csvfilepath, filesavepath, dirname string, width, height int) bool {
	status := false

	if (csvfilepath != "") && (filesavepath != "") && (dirname != "") {
		// initialize data directory
		if cc.dirhandler.MakeDirectory(filesavepath, dirname) {
			// initalize colorcodes
			cc.colorCodes = models.ReadColorCode(csvfilepath)
			path := filesavepath + dirname + "/"

			// create solid images
			if len(cc.colorCodes) > 0 {
				for _, data := range cc.colorCodes {
					rawimage := cc.imgcontroller.CreateSolidImage(*data.GenerateColorRGBA(), width, height)
					cc.iohandler.StreamOutPNGFile(path, data.GetName(), rawimage)
				}

				// status update
				status = true
			}
		} else {
			// initalize colorcodes
			// - just update
			cc.colorCodes = models.ReadColorCode(csvfilepath)

			// status update
			status = true
		}
	}
	return status
}

/*
SaveColorCodePatchData :save color code pathc data as CSV file
	in	;savepath, filename string
	out	;bool
*/
func (cc *colorCheckerViewController) SaveColorCodePatchData(savepath, filename string) bool {
	status := false

	if len(cc.colorCodes) != 0 {
		if savepath != "" && filename != "" {
			// make string data from property
			data := make([][]string, 0)
			for _, obj := range cc.colorCodes {
				dataString := obj.SerializeData()
				data = append(data, dataString)
			}

			// save data
			if cc.iohandler.WriteCSVFile(savepath, filename, data) {
				// status update
				status = true
			}
		}
	}

	return status
}
