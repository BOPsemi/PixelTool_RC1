package viewcontrollers

import (
	"PixelTool_RC1/controllers"
	"PixelTool_RC1/models"
	"PixelTool_RC1/util"
	"image/color"
	"math/rand"
	"time"
)

/*
WhitePixelCheckerViewController :interface of white pixel checker VC
- create white pixel raw image
*/
type WhitePixelCheckerViewController interface {
	CreateWhitePixelPatch(csvfilepath, filesavepath, dirname string, width, height int) bool
}

// definition of structure
type whitePixelCheckerViewController struct {
	dirhandler    util.DirectoryHandler       // directory handler
	iohandler     util.IOUtil                 // io handler for reading csv file
	randomizer    util.Randomizer             // white pixel randomizer
	imgcontroller controllers.ImageController // image controller

	// properties
	level     []int        // white pixel level
	count     []int        // white pixel count
	wprawdata []int        // serialized white pixel data
	imagedata []color.RGBA // image data
}

// NewWhitePixelCheckerViewController :initializer of struct
func NewWhitePixelCheckerViewController() WhitePixelCheckerViewController {
	obj := new(whitePixelCheckerViewController)

	//initialize handlers
	obj.dirhandler = util.NewDirectoryHandler()
	obj.iohandler = util.NewIOUtil()
	obj.randomizer = util.NewRandomizer()
	obj.imgcontroller = controllers.NewImageController()

	// initalizer properties
	obj.level = make([]int, 0)
	obj.count = make([]int, 0)
	obj.wprawdata = make([]int, 0)
	obj.imagedata = make([]color.RGBA, 0)

	return obj
}

/*
CreateWhitePixelPatch :
	in	;csvfilepath string
	out ;bool
*/
func (vc *whitePixelCheckerViewController) CreateWhitePixelPatch(csvfilepath, filesavepath, dirname string, width, height int) bool {
	status := false

	if csvfilepath != "" && width > 0 && height > 0 {
		// read csv file
		wp := models.ReadWhitePixel(csvfilepath)

		// check white pixel size
		// - fail len(wp) == 0
		if len(wp) > 0 {

			// extract white pixel level and count
			vc.level, vc.count = vc.extractWhitePixelCount(wp)

			// randamize white pixe according to patch size
			iteration := width * height
			rwprawdata := vc.randomizer.Run(vc.count, iteration)

			// serialize raw data
			for index, data := range rwprawdata {
				if data != 0 {
					for i := 0; i < data; i++ {
						vc.wprawdata = append(vc.wprawdata, vc.level[index])
					}
				}
			}

			// randamize seriarized data
			for i := iteration - 1; i >= 0; i-- {
				rand.Seed(time.Now().UnixNano())
				j := rand.Intn(i + 1)
				vc.wprawdata[i], vc.wprawdata[j] = vc.wprawdata[j], vc.wprawdata[i]
			}

			// create image data
			for _, data := range vc.wprawdata {
				rawdata := color.RGBA{
					R: uint8(data),
					G: uint8(data),
					B: uint8(data),
					A: 255,
				}
				vc.imagedata = append(vc.imagedata, rawdata)
			}

			// save WhitePixel image data
			if filesavepath != "" && dirname != "" {
				// initialize data directory
				if vc.dirhandler.MakeDirectory(filesavepath, dirname) {
					path := filesavepath + dirname + "/"

					// create white pixel image
					rawimage := vc.imgcontroller.CreateImage(vc.imagedata, height, width)
					vc.iohandler.StreamOutPNGFile(path, "white_pixel", rawimage)

					// status update
					status = true
				} else {
					path := filesavepath + dirname + "/"

					// create white pixel image
					rawimage := vc.imgcontroller.CreateImage(vc.imagedata, height, width)
					vc.iohandler.StreamOutPNGFile(path, "white_pixel", rawimage)

					// status update
					status = true

				}
			}
		}
	}

	return status
}

/*
extract level and count from white pixel structure
*/
func (vc *whitePixelCheckerViewController) extractWhitePixelCount(wp []models.WhitePixel) (level, count []int) {
	whitepixelcount := make([]int, 0) // buffer for count
	whitepixellevel := make([]int, 0) // buffer for level

	for _, data := range wp {
		whitepixelcount = append(whitepixelcount, data.GetCount())
		whitepixellevel = append(whitepixellevel, data.GetLevel())
	}

	return whitepixellevel, whitepixelcount
}
