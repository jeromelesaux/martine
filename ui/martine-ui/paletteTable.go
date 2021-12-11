package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type PaletteTable struct {
	widget.Table
	Palette                color.Palette
	ColorCallbackFunc      func(color.Color)
	IndexCallbackFunc      func(int)
	SetPaletteCallbackFunc func(color.Palette)
}

func NewPaletteTable(p color.Palette, colorSelected func(color.Color), indexSelected func(int), setPalette func(color.Palette)) *PaletteTable {
	paletteTable := &PaletteTable{}
	paletteTable.Palette = p
	paletteTable.ColorCallbackFunc = colorSelected
	paletteTable.IndexCallbackFunc = indexSelected
	paletteTable.SetPaletteCallbackFunc = setPalette
	paletteTable.CreateCell = paletteTable.PaletteCreate
	paletteTable.Length = paletteTable.PaletteLength
	paletteTable.UpdateCell = paletteTable.PaletteUpdate
	paletteTable.OnSelected = paletteTable.PaletteSelect
	paletteTable.ExtendBaseWidget(paletteTable)

	return paletteTable
}

func (p *PaletteTable) SubstitueColor(index int, newColor color.Color) {
	if index < 0 || index > len(p.Palette) {
		return
	}
	p.Palette[index] = newColor
	p.UpdateCell(widget.TableCellID{Row: 0, Col: index}, canvas.NewRectangle(newColor))
}

func (p *PaletteTable) PaletteLength() (int, int) {
	return 1, len(p.Palette)
}

func (p *PaletteTable) PaletteCreate() fyne.CanvasObject {
	r := canvas.NewRectangle(color.White)
	r.SetMinSize(fyne.NewSize(30, 30))
	return r
}

func (p *PaletteTable) PaletteUpdate(id widget.TableCellID, o fyne.CanvasObject) {
	c := p.Palette[id.Col]
	o.(*canvas.Rectangle).FillColor = c
}

func (p *PaletteTable) PaletteSelect(id widget.TableCellID) {
	c := p.Palette[id.Col]
	if p.ColorCallbackFunc != nil {
		p.ColorCallbackFunc(c)
	}
	if p.IndexCallbackFunc != nil {
		p.IndexCallbackFunc(id.Col)
	}
	if p.SetPaletteCallbackFunc != nil {
		p.SetPaletteCallbackFunc(p.Palette)
	}
}
