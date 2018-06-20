package viewcontrollers

import (
	"PixelTool_RC1/controllers"
	"PixelTool_RC1/util"
	"image"
	"image/color"
)

/*
NoiseAdditionViewController :interface of noise addtion vc
*/
type NoiseAdditionViewController interface {
	SetImageDataForWhitePixelAddition(basefilepath, noisefilepath string) bool
	CreateImageWithWhitePixel(darklevel int, filename, filesavepath, dirname string) bool
}

// definition of structure
type noiseAdditionViewController struct {
	base  []color.RGBA // base image data
	noise []color.RGBA // noise image data

	width  int // image width
	height int // image height
}

/*
NewNoiseAdditionViewController :initializer of noise addtion vc
*/
func NewNoiseAdditionViewController() NoiseAdditionViewController {
	obj := new(noiseAdditionViewController)

	// init properties
	obj.base = make([]color.RGBA, 0)
	obj.noise = make([]color.RGBA, 0)

	return obj
}

/*
SetImageDataForWhitePixelAddition :
	in	;basefilepath, noisefilepath string
	out	;bool
*/
func (vc *noiseAdditionViewController) SetImageDataForWhitePixelAddition(basefilepath, noisefilepath string) bool {
	status := false

	if basefilepath != "" && noisefilepath != "" {

		// read base image data
		baseData := vc.imageFileOpen(basefilepath)
		if len(baseData) > 0 {
			vc.base = baseData
		}

		// read noise image data
		noiseData := vc.imageFileOpen(noisefilepath)
		if len(noiseData) > 0 {
			vc.noise = noiseData
		}

		// update status
		if len(baseData)*len(noiseData) > 0 {
			status = true
		}

		// debug
		//fmt.Println(vc.base)
	}

	return status
}

func (vc *noiseAdditionViewController) imageFileOpen(filepath string) []color.RGBA {
	data := make([]color.RGBA, 0)

	if filepath != "" {
		iohandler := util.NewIOUtil()
		imageData := iohandler.ReadImageFile(filepath)

		if imageData != nil {
			// check image width and height
			vc.height = imageData.Bounds().Size().Y
			vc.width = imageData.Bounds().Size().X

			// type change
			if img, ok := imageData.(*image.RGBA); ok {

				// call iamge controller
				imagecontroller := controllers.NewImageController()
				rgbaImagaData := imagecontroller.SerializeImage(img)

				// check data size
				if len(rgbaImagaData) > 0 {
					// update data
					data = rgbaImagaData
				}
			}
		}
	}

	return data
}

/*
ImageWithWhitePixel :
	in	;darklevel int
	out	;bool
*/
func (vc *noiseAdditionViewController) CreateImageWithWhitePixel(darklevel int, filename, filesavepath, dirname string) bool {
	status := false

	if darklevel > 0 {
		if len(vc.base)*len(vc.noise) > 0 {
			noiseController := controllers.NewNoiseController()
			rawImageData := noiseController.AddWhitePixelNoise(vc.base, vc.noise, darklevel)

			// check image size
			if len(rawImageData) == vc.width*vc.height {

				// create image from row data
				imageController := controllers.NewImageController()
				img := imageController.CreateImage(rawImageData, vc.height, vc.width)

				// stream out image data
				path := filesavepath + dirname + "/"

				// directory handling
				dirhandler := util.NewDirectoryHandler()
				if dirhandler.MakeDirectory(filesavepath, dirname) {

					// stream out file
					status = vc.streamOutImageData(path, filename, img)
				} else {
					// stream out file
					status = vc.streamOutImageData(path, filename, img)
				}
			}
		}
	}

	return status
}

func (vc *noiseAdditionViewController) streamOutImageData(path, filename string, data *image.RGBA) bool {
	iohandler := util.NewIOUtil()
	return iohandler.StreamOutPNGFile(path, filename, data)
}
