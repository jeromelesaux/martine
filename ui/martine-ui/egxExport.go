package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/diskimage"
	"github.com/jeromelesaux/martine/export/m4"

	"github.com/jeromelesaux/martine/export/ocpartstudio"
	"github.com/jeromelesaux/martine/export/snapshot"
	"github.com/jeromelesaux/martine/gfx/common"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

var egxFilename = "aa.scr"

func (m *MartineUI) exportEgxDialog(ie *menu.ImageExport, w fyne.Window) {
	m2host := widget.NewEntry()
	m2host.SetPlaceHolder("Set your M2 IP here.")

	ie.Reset()
	cont := container.NewVBox(
		container.NewHBox(
			widget.NewCheck("import all file in Dsk", func(b bool) {
				ie.ExportDsk = b
			}),
			widget.NewCheck("export text file", func(b bool) {
				ie.ExportText = b
			}),
			widget.NewCheck("export Json file", func(b bool) {
				ie.ExportJson = b
			}),
			widget.NewCheck("add amsdos header", func(b bool) {
				ie.ExportWithAmsdosHeader = b
			}),
			widget.NewCheck("apply zigzag", func(b bool) {
				ie.ExportZigzag = b
			}),
			widget.NewCheck("export to M2", func(b bool) {
				ie.ExportToM2 = b
				ie.M2IP = m2host.Text
			}),
		),

		widget.NewLabel("Compression type:"),
		widget.NewSelect([]string{"none", "rle", "rle 16bits", "Lz4 Classic", "Lz4 Raw", "zx0 crunch"},
			func(s string) {
				switch s {
				case "none":
					ie.ExportCompression = 0
				case "rle":
					ie.ExportCompression = 1
				case "rle 16bits":
					ie.ExportCompression = 2
				case "Lz4 Classic":
					ie.ExportCompression = 3
				case "Lz4 Raw":
					ie.ExportCompression = 4
				case "zx0 crunch":
					ie.ExportCompression = 5
				}
			}),
		m2host,
		widget.NewButtonWithIcon("Export into folder", theme.DocumentSaveIcon(), func() {
			fo := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, m.window)
					return
				}
				if lu == nil {
					// cancel button
					return
				}
				ie.ExportFolderPath = lu.Path()
				m.egx.ResultImage.Path = lu.Path()
				fmt.Println(ie.ExportFolderPath)
				m.ExportEgxImage(m.egx)

				// apply and export
			}, m.window)
			fo.Show()
		}),
	)

	d := dialog.NewCustom("Export options", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
}

func (m *MartineUI) ExportEgxImage(me *menu.DoubleImageMenu) {
	pi := dialog.NewProgressInfinite("Saving....", "Please wait.", m.window)
	pi.Show()
	context := m.NewContext(&me.LeftImage, true)
	if context == nil {
		pi.Hide()
		return
	}
	context.OutputPath = me.ResultImage.Path
	context.Dsk = m.egxExport.ExportDsk
	if m.egxExport.ExportWithAmsdosHeader {
		context.NoAmsdosHeader = false
	} else {
		context.NoAmsdosHeader = true
	}
	if me.ResultImage.EgxType == 1 {
		context.EgxFormat = export.Egx1Mode
	} else {
		context.EgxFormat = export.Egx2Mode
	}
	context.EgxMode1 = uint8(me.LeftImage.Mode)
	context.EgxMode2 = uint8(me.RightImage.Mode)

	// palette export
	defer func() {
		os.Remove("temporary_palette.kit")
	}()
	if err := ocpartstudio.SaveKit("temporary_palette.kit", me.ResultImage.Palette, false); err != nil {
		pi.Hide()
		dialog.ShowError(err, m.window)
	}
	context.KitPath = "temporary_palette.kit"

	if !context.Overscan {
		if err := ocpartstudio.EgxLoader(me.ResultImage.Path+string(filepath.Separator)+egxFilename, me.ResultImage.Palette, uint8(me.LeftImage.Mode), uint8(me.RightImage.Mode), context); err != nil {
			pi.Hide()
			dialog.ShowError(err, m.window)
			return
		}
	} else {
		if err := common.Export(me.ResultImage.Path+string(filepath.Separator)+egxFilename, me.ResultImage.Data, me.ResultImage.Palette, uint8(me.ResultImage.EgxType), context); err != nil {
			pi.Hide()
			dialog.ShowError(err, m.window)
			return
		}
	}
	if m.egxExport.ExportDsk {
		if err := diskimage.ImportInDsk(me.ResultImage.Path, context); err != nil {
			dialog.NewError(err, m.window).Show()
			return
		}
	}
	if context.Sna {
		if context.Overscan {
			var gfxFile string
			for _, v := range context.DskFiles {
				if filepath.Ext(v) == ".SCR" {
					gfxFile = v
					break
				}
			}
			context.SnaPath = filepath.Join(me.ResultImage.Path, "test.sna")
			if err := snapshot.ImportInSna(gfxFile, context.SnaPath, 0); err != nil {
				dialog.NewError(err, m.window).Show()
				return
			}
		}
	}
	if m.egxExport.ExportToM2 {
		if err := m4.ImportInM4(context); err != nil {
			dialog.NewError(err, m.window).Show()
			fmt.Fprintf(os.Stderr, "Cannot send to M4 error :%v\n", err)
		}
	}
	pi.Hide()
	dialog.ShowInformation("Save", "Your files are save in folder \n"+m.egxExport.ExportFolderPath, m.window)

}
