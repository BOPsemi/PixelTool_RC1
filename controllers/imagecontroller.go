package controllers

import (
	"PixelTool_RC1/util"
	"image"
	"image/color"
	"image/draw"
)

var (
	LeftMargin   = 30
	RightMargin  = 30
	TopMargin    = 30
	BottomMargin = 30
	Spacing      = 10
)

/*
Image controller package
*/

/*
ImageController :interface of image controller
*/
type ImageController interface {
	// create raw image data
	CreateImage(data []color.RGBA, height, width int) *image.RGBA
	CreateSolidImage(data color.RGBA, height, width int) *image.RGBA

	// read pixel data from image file
	SerializeImage(img *image.RGBA) []color.RGBA

	// create 24 Macbeth chart
	Create24MacbethChart(path, filename string) bool
}

// definition of image controller
type imageContrller struct {
}

// NewImageController :initializer of image controller
func NewImageController() ImageController {
	obj := new(imageContrller)

	return obj
}

/*
CreateImage :create image from data
	data 				<-[]color.RGBA
	height, width 		<-int
*/
func (im *imageContrller) CreateImage(data []color.RGBA, height, width int) *image.RGBA {
	img := new(image.RGBA)

	if height > 0 && width > 0 {
		// check data size
		if (height * width) == len(data) {

			// create image
			canvas := image.NewRGBA(image.Rect(0, 0, width, height))
			for i := 0; i < width; i++ {
				for j := 0; j < height; j++ {
					index := width*i + j

					// raw data
					rawData := color.RGBA{
						R: data[index].R,
						G: data[index].G,
						B: data[index].B,
						A: 255,
					}

					// draw the raw data on canvas
					canvas.Set(i, j, rawData)
				}
			}

			// update image
			img = canvas
		}
	}
	return img
}

/*
CreateSolidImage(data color.RGBA, height, width int) *image.RGBA
*/
func (im *imageContrller) CreateSolidImage(data color.RGBA, height, width int) *image.RGBA {
	img := new(image.RGBA)

	if height > 0 && width > 0 {

		// create image
		canvas := image.NewRGBA(image.Rect(0, 0, width, height))
		for i := 0; i < width; i++ {
			for j := 0; j < height; j++ {

				// draw the raw data on canvas
				canvas.Set(i, j, data)
			}
		}

		// update image
		img = canvas

	}

	return img
}

/*
SerializeImage :serialize image data to color.RGBA slice
	img 	:*image.RGBA
*/
func (im *imageContrller) SerializeImage(img *image.RGBA) []color.RGBA {
	data := make([]color.RGBA, 0)

	if img != nil {
		for i := 0; i < img.Bounds().Size().X; i++ {
			for j := 0; j < img.Bounds().Size().Y; j++ {

				// extract point data
				rgba := img.At(i, j)

				// each channel data
				r, g, b, a := rgba.RGBA()

				// create raw data
				rawdata := color.RGBA{
					R: uint8(r),
					G: uint8(g),
					B: uint8(b),
					A: uint8(a),
				}

				// stack data
				data = append(data, rawdata)
			}
		}
	}

	return data
}

/*
Create24MacbethChart :
in	;path, filename string
out	;bool
*/
func (im *imageContrller) Create24MacbethChart(path, filename string) bool {
	status := false

	if path != "" && filename != "" {
		// initialize directory handler
		dirHandler := util.NewDirectoryHandler()
		ioHandler := util.NewIOUtil() // initialize IO handler

		// get file names from directory path
		files, names := dirHandler.GetFileListInDirectory(path)

		// make patch Data Map
		patchDataMap := make(map[string]image.Image, 0)
		for index := 0; index < len(files); index++ {
			patchDataMap[names[index]] = ioHandler.ReadImageFile(files[index])
		}

		width := patchDataMap["Black"].Bounds().Size().X
		height := patchDataMap["Black"].Bounds().Size().Y

		// canvas
		canvas := im.makeCanvas(width, height)

		// point calculator
		for index := 0; index < 24; index++ {

			// calculate start and end point
			start, end := im.pointCalculator(index, width, height)
			srcRect := image.Rectangle{start, end}

			// create image
			draw.Draw(canvas, srcRect, patchDataMap[im.returnPatchNameString(index)], image.Pt(0, 0), draw.Over)
		}

		// stream out final image
		if canvas != nil {
			if ioHandler.StreamOutPNGFile(path, filename, canvas) {
				status = true
			}
		}
	}

	return status
}

/*
Return Patch Name
*/
func (im *imageContrller) returnPatchNameString(index int) string {
	patchName := ""
	switch index {
	case 0:
		patchName = "DarkSkin"
	case 1:
		patchName = "LightSkin"
	case 2:
		patchName = "BlueSky"
	case 3:
		patchName = "Foliage"
	case 4:
		patchName = "BlueFlower"
	case 5:
		patchName = "BluishGreen"
	case 6:
		patchName = "Orange"
	case 7:
		patchName = "PurplishBlue"
	case 8:
		patchName = "ModerateRed"
	case 9:
		patchName = "Purple"
	case 10:
		patchName = "YellowGreen"
	case 11:
		patchName = "OrangeYellow"
	case 12:
		patchName = "Blue"
	case 13:
		patchName = "Green"
	case 14:
		patchName = "Red"
	case 15:
		patchName = "Yellow"
	case 16:
		patchName = "Magenta"
	case 17:
		patchName = "Cyan"
	case 18:
		patchName = "White"
	case 19:
		patchName = "Neutral8"
	case 20:
		patchName = "Neutral6p5"
	case 21:
		patchName = "Neutral5"
	case 22:
		patchName = "Neutral3p5"
	case 23:
		patchName = "Black"
	default:
		patchName = ""
	}

	return patchName
}

/*
Calculate each Patch point from index information
*/
func (im *imageContrller) pointCalculator(index, width, height int) (start image.Point, end image.Point) {
	// variables

	xpoint := 0
	ypoint := 0

	// point data
	if index < 6 {
		// 1st row
		i := index
		xpoint = LeftMargin + i*width + i*Spacing
		ypoint = TopMargin

	} else if index > 5 && index < 12 {
		// 2nd row
		i := (index - 6)
		xpoint = LeftMargin + i*width + i*Spacing
		ypoint = TopMargin + height + Spacing

	} else if index > 10 && index < 18 {
		// 3rd row
		i := (index - 12)
		xpoint = LeftMargin + i*width + i*Spacing
		ypoint = TopMargin + 2*height + 2*Spacing

	} else {
		// 4th row
		i := (index - 18)
		xpoint = LeftMargin + i*width + i*Spacing
		ypoint = TopMargin + 3*height + 3*Spacing
	}

	startPoint := image.Pt(xpoint, ypoint)
	endPoint := image.Pt(xpoint+width, ypoint+width)

	return startPoint, endPoint
}

/*
Make Blank canvas, color is pure black
*/
func (im *imageContrller) makeCanvas(width, height int) *image.RGBA {

	// calculate canvas size
	backgroundWidth := LeftMargin + 5*Spacing + RightMargin + width*6
	backgroundHeight := TopMargin + 3*Spacing + BottomMargin + height*4

	// create canvas
	canvas := image.NewRGBA(image.Rect(0, 0, backgroundWidth, backgroundHeight))

	for i := 0; i < backgroundWidth; i++ {
		for j := 0; j < backgroundHeight; j++ {
			canvas.Set(i, j, color.RGBA{0, 0, 0, 255})
		}
	}

	return canvas
}
