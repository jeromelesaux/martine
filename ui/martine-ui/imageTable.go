package ui

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type ImageTable struct {
	widget.Table
	images                *[][]canvas.Image
	ImageCallbackFunc     func(*canvas.Image)
	IndexCallbackFunc     func(int, int)
	SetImagesCallbackFunc func(*[][]canvas.Image)
	imageSize             fyne.Size
	rowsNumber            int
	colsNumber            int
}

func NewImageTable(
	images *[][]canvas.Image,
	imageSize fyne.Size,
	nbRows, nbCols int,
	imageSelected func(*canvas.Image),
	indexSelected func(int, int),
	setImages func(*[][]canvas.Image)) *ImageTable {
	if len(*images) != nbRows || len((*images)[0]) != nbCols {
		panic("images matrix must corresponds to number of rows and columns")
	}
	imageTable := &ImageTable{}
	imageTable.images = images
	imageTable.ImageCallbackFunc = imageSelected
	imageTable.IndexCallbackFunc = indexSelected
	imageTable.SetImagesCallbackFunc = setImages
	imageTable.CreateCell = imageTable.ImageCreate
	imageTable.Length = imageTable.ImagesLength
	imageTable.UpdateCell = imageTable.ImageUpdate
	imageTable.OnSelected = imageTable.ImageSelect
	imageTable.imageSize = imageSize
	imageTable.rowsNumber = nbRows
	imageTable.colsNumber = nbCols
	imageTable.ExtendBaseWidget(imageTable)

	return imageTable
}

func (i *ImageTable) SubstitueImage(row, col int, newImage canvas.Image) {
	if row < 0 || row > len(*i.images) {
		return
	}
	if col < 0 || col > len((*i.images)[0]) {
		return
	}
	(*i.images)[row][col] = newImage
	i.UpdateCell(widget.TableCellID{Row: row, Col: col}, &newImage)
}

func (i *ImageTable) UpdateAll() {
	for x := 0; x < i.rowsNumber; x++ {
		for y := 0; y < i.colsNumber; y++ {
			i.UpdateCell(widget.TableCellID{Row: x, Col: y}, &(*i.images)[x][y])
		}
	}
	i.Refresh()
}

func (i *ImageTable) ImagesLength() (row int, col int) {
	return i.rowsNumber, i.colsNumber
}

func (i *ImageTable) ImageCreate() fyne.CanvasObject {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(i.Size().Width), int(i.Size().Width)}})
	r := canvas.NewImageFromImage(img)
	r.SetMinSize(i.imageSize)
	return r
}

func (i *ImageTable) ImageUpdate(id widget.TableCellID, o fyne.CanvasObject) {
	c := (*i.images)[id.Row][id.Col]
	o.(*canvas.Image).Image = c.Image
}

func (i *ImageTable) ImageSelect(id widget.TableCellID) {
	c := (*i.images)[id.Row][id.Col]
	if i.ImageCallbackFunc != nil {
		i.ImageCallbackFunc(&c)
	}
	if i.IndexCallbackFunc != nil {
		i.IndexCallbackFunc(id.Row, id.Col)
	}
	if i.SetImagesCallbackFunc != nil {
		i.SetImagesCallbackFunc(i.images)
	}
}
