package views

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

/*
ImageViewer :image viewer
*/
type ImageViewer struct {
	reader *gui.QImageReader            // image reader
	pixmap *gui.QPixmap                 // pixel map
	item   *widgets.QGraphicsPixmapItem // pixel item
	scene  *widgets.QGraphicsScene      // graphic scene

	Height int // image height
	Width  int // image width

	Cell *widgets.QGraphicsView
}

/*
NewImageViewer :initializer of image viewer
*/
func NewImageViewer(path string, scale float64) *ImageViewer {
	obj := new(ImageViewer)

	// initialize widgets
	obj.Cell = widgets.NewQGraphicsView(nil)

	if path != "" {
		obj.reader = gui.NewQImageReader3(path, core.NewQByteArray())
	} else {
		obj.reader = gui.NewQImageReader()
	}

	// make scene
	obj.makeScene(scale, false)

	return obj
}

/*
SetImageView : image viewer setting
	in	;string, scale
	out	;
*/
func (iv *ImageViewer) SetImageView(path string, scale float64) {
	iv.reader.SetFileName(path)
	iv.reader.SetFormat(core.NewQByteArray())

	// make scene
	iv.makeScene(scale, true)
}

// make scene
func (iv *ImageViewer) makeScene(scale float64, redraw bool) {
	iv.pixmap = gui.QPixmap_FromImageReader(iv.reader, core.Qt__AutoColor)
	iv.item = widgets.NewQGraphicsPixmapItem2(iv.pixmap, nil)

	// set scale
	var viewScale float64
	if scale < 0 {
		viewScale = 1.0
	} else {
		viewScale = scale
	}
	iv.item.SetScale(viewScale)

	// add item to scene
	iv.scene = widgets.NewQGraphicsScene(nil)
	iv.scene.AddItem(iv.item)

	if redraw {
		iv.Height = int(float64(iv.pixmap.Size().Height()) * scale)
		iv.Width = int(float64(iv.pixmap.Size().Width()) * scale)
	} else {
		iv.Height = 245
		iv.Width = 355
	}

	// create image view
	iv.Cell.SetScene(iv.scene)

}
