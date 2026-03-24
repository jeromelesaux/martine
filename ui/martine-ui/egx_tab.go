package ui

import (
	"fmt"
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	wgt "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/log"

	"github.com/jeromelesaux/martine/convert/screen"
	"github.com/jeromelesaux/martine/convert/screen/overscan"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/effect"
	"github.com/jeromelesaux/martine/ui/martine-ui/menu"
)

func (m *MartineUI) newEgxTab(di *menu.DoubleImageMenu) *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Image 1", m.newImageTransfertTab(di.LeftImage)),
		container.NewTabItem("Image 2", m.newImageTransfertTab(di.RightImage)),
		container.NewTabItem("Egx", m.newEgxTabItem(di)),
	)
}

// nolint:funlen
func (m *MartineUI) MergeImages(di *menu.DoubleImageMenu) {
	if di.RightImage.Cfg.ScrCfg.Mode == di.LeftImage.Cfg.ScrCfg.Mode {
		dialog.ShowError(fmt.Errorf("mode between the images must differ"), m.window)
		return
	}
	if di.RightImage.Cfg.ScrCfg.IsPlus != di.LeftImage.Cfg.ScrCfg.IsPlus {
		dialog.ShowError(fmt.Errorf("plus mode between the images differs, set the same mode"), m.window)
		return
	}
	if di.RightImage.Cfg.ScrCfg.Type != di.LeftImage.Cfg.ScrCfg.Type {
		dialog.ShowError(fmt.Errorf("the size of both image differs, set the same size for both"), m.window)
		return
	}
	var im *menu.ImageMenu
	var im2 *menu.ImageMenu
	var palette2 color.Palette
	if di.LeftImage.Cfg.ScrCfg.Mode < di.RightImage.Cfg.ScrCfg.Mode {
		im = di.LeftImage
		palette2 = append(palette2, di.LeftImage.Palette()...)
		im2 = di.RightImage
	} else {
		im = di.RightImage
		palette2 = append(palette2, di.RightImage.Palette()...)
		im2 = di.LeftImage
	}

	di.ResultImage.Cfg.ScrCfg.Mode = im.Cfg.ScrCfg.Mode
	di.ResultImage.Cfg.ScrCfg.Size = im.Cfg.ScrCfg.Size
	di.ResultImage.Cfg.ScrCfg.Type = im.Cfg.ScrCfg.Type
	di.ResultImage.Cfg.ScrCfg.Export = append(di.ResultImage.Cfg.ScrCfg.Export, im.Cfg.ScrCfg.Export...)

	out, downgraded, pal, _, err := gfx.ApplyOneImage(im.CpcImage().Image, im.Cfg, im.Palette(), im.Cfg.ScrCfg.Mode)
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	im.Cfg.PalCfg.Palette = pal
	im.Data = out
	im.SetCpcImage(downgraded)
	im.SetPalette(im.Palette())
	im.SetPaletteImage(png.PalToImage(im.Palette()))

	out, downgraded, pal, _, err = gfx.ApplyOneImage(im2.CpcImage().Image, im2.Cfg, palette2, im2.Cfg.ScrCfg.Mode)
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	im2.Cfg.PalCfg.Palette = pal
	im2.Data = out
	im2.SetCpcImage(downgraded)
	im2.SetPalette(im2.Palette())
	im2.SetPaletteImage(png.PalToImage(im.Palette()))

	pi := wgt.NewProgressInfinite("Computing, Please wait.", m.window)
	pi.Show()
	di.ResultImage.Cfg.PalCfg.Palette = im.Palette()
	di.ResultImage.Cfg.ScrCfg.ResetExport()
	di.ResultImage.Cfg.ScrCfg.Type = im.Cfg.ScrCfg.Type
	di.ResultImage.Cfg.ScrCfg.Export = append(di.ResultImage.Cfg.ScrCfg.Export, im.Cfg.ScrCfg.Export...)
	di.ResultImage.Cfg.ScrCfg.Size = im.Cfg.ScrCfg.Size

	res, pal, egxType, err := effect.EgxRaw(im.CpcImage().Image, im2.CpcImage().Image, di.ResultImage.Cfg.PalCfg.Palette, int(im.Cfg.ScrCfg.Mode), int(im2.Cfg.ScrCfg.Mode), di.ResultImage.Cfg)
	pi.Hide()
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}
	di.ResultImage.EgxType = egxType
	di.ResultImage.Data = res
	di.ResultImage.Cfg.PalCfg.Palette = pal
	var img image.Image
	if di.ResultImage.Cfg.ScrCfg.Type == config.FullscreenFormat {
		img, err = overscan.OverscanRawToImg(di.ResultImage.Data, im.Cfg.ScrCfg.Mode, di.ResultImage.Cfg.PalCfg.Palette)
		if err != nil {
			dialog.ShowError(err, m.window)
			return
		}
	} else {
		img, err = screen.ScrRawToImg(di.ResultImage.Data, 0, di.ResultImage.Cfg.PalCfg.Palette)
		if err != nil {
			dialog.ShowError(err, m.window)
			return
		}
	}
	if di.ResultImage.Cfg.ScrCfg.Type.IsFullScreen() {
		if im.Cfg.ScrCfg.Mode == 2 || im2.Cfg.ScrCfg.Mode == 2 {
			di.ResultImage.Cfg.ScrCfg.Type = config.Egx2FullscreenFormat
		} else {
			di.ResultImage.Cfg.ScrCfg.Type = config.Egx1FullscreenFormat
		}
	} else {
		if im.Cfg.ScrCfg.Mode == 2 || im2.Cfg.ScrCfg.Mode == 2 {
			di.ResultImage.Cfg.ScrCfg.Type = config.Egx2Format
		} else {
			di.ResultImage.Cfg.ScrCfg.Type = config.Egx1Format
		}
	}
	di.ResultImage.CpcResultImage.Image = img
	di.ResultImage.CpcResultImage.Refresh()
	di.ResultImage.PaletteImage.Image = png.PalToImage(im.Palette())
	di.ResultImage.PaletteImage.Refresh()
}

// nolint: funlen
func (m *MartineUI) newEgxTabItem(di *menu.DoubleImageMenu) *fyne.Container {
	di.ResultImage = menu.NewMergedImageMenu()
	di.ResultImage.CpcLeftImage = di.LeftImage.CpcImage()
	di.ResultImage.CpcRightImage = di.RightImage.CpcImage()
	di.ResultImage.LeftPalette = di.LeftImage.Palette()
	di.ResultImage.RightPalette = di.RightImage.Palette()
	di.ResultImage.LeftPaletteImage = di.LeftImage.PaletteImage()
	di.ResultImage.RightPaletteImage = di.RightImage.PaletteImage()

	return container.New(
		layout.NewGridLayoutWithRows(3),
		container.New(
			layout.NewGridLayoutWithColumns(2),
			di.LeftImage.CpcImage(),
			di.RightImage.CpcImage(),
		),

		container.New(
			layout.NewVBoxLayout(),

			container.New(
				layout.NewGridLayoutWithColumns(2),
				container.New(
					layout.NewGridLayoutWithColumns(1),
					di.LeftImage.PaletteImage(),
				),
				container.New(
					layout.NewGridLayoutWithColumns(1),
					di.RightImage.PaletteImage(),
				),

				widget.NewButtonWithIcon("Merge image", theme.MediaPlayIcon(), func() {
					m.MergeImages(m.egx)
				}),
				widget.NewButtonWithIcon("Refresh UI", theme.ComputerIcon(), func() {
					di.ResultImage.CpcLeftImage.Image = di.LeftImage.CpcImage().Image
					di.ResultImage.CpcRightImage.Image = di.RightImage.CpcImage().Image
					di.ResultImage.CpcLeftImage.Refresh()
					di.ResultImage.CpcResultImage.Refresh()
					di.ResultImage.LeftPalette = di.LeftImage.Palette()
					di.ResultImage.RightPalette = di.RightImage.Palette()
					di.ResultImage.LeftPaletteImage.Image = di.LeftImage.PaletteImage().Image
					di.ResultImage.LeftPaletteImage.Refresh()
					di.ResultImage.RightPaletteImage.Image = di.RightImage.PaletteImage().Image
					di.ResultImage.LeftPaletteImage.Refresh()
					di.ResultImage.PaletteImage.Image = png.PalToImage(di.ResultImage.Cfg.PalCfg.Palette)
					di.ResultImage.PaletteImage.Refresh()
				}),
				widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
					// export the egx image
					di.ResultImage.Cfg.ScrCfg.IsPlus = di.LeftImage.Cfg.ScrCfg.IsPlus
					m.exportEgxDialog(di.ResultImage.Cfg, m.window)
				}),
				widget.NewButton("show cmd", func() {
					e := widget.NewMultiLineEntry()
					e.SetText(di.CmdLine())
					d := dialog.NewCustom("Command line generated",
						"Ok",
						e,
						m.window)
					log.GetLogger().Info("%s\n", di.CmdLine())
					size := m.window.Content().Size()
					size = fyne.Size{Width: size.Width / 2, Height: size.Height / 2}
					d.Resize(size)
					d.Show()
				}),
			),
		),
		container.New(
			layout.NewGridLayoutWithColumns(2),
			di.ResultImage.CpcResultImage,
			container.New(
				layout.NewGridLayoutWithRows(3),
				di.ResultImage.PaletteImage,
			),
		),
	)
}
