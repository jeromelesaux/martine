package animate

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/jeromelesaux/martine/assembly"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/compression"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/errors"
	"github.com/jeromelesaux/martine/gfx/transformation"
	"github.com/jeromelesaux/martine/log"
	zx0 "github.com/jeromelesaux/zx0/encode"
)

type DeltaExportFormat int

var (
	DeltaExportV1 DeltaExportFormat = 1
	DeltaExportV2 DeltaExportFormat = 2
)

// nolint: funlen
func DeltaPackingMemory(images []image.Image, cfg *config.MartineConfig, initialAddress uint16, mode uint8) ([]*transformation.DeltaCollection, [][]byte, color.Palette, error) {
	var isSprite bool = true
	maxImages := 22
	var pad int = 1
	var err error
	var palette color.Palette
	if !cfg.CustomDimension && cfg.ScrCfg.Type != config.SpriteHardFormat {
		isSprite = false
	}
	if len(images) <= 1 {
		return nil, nil, palette, fmt.Errorf("need more than one image to proceed")
	}
	if len(images) > maxImages {
		log.GetLogger().Error("Warning gif exceed 30 images. Will corrupt the number of images.")
		pad = len(images) / maxImages
	}
	rawImages := make([][]byte, 0)
	deltaData := make([]*transformation.DeltaCollection, 0)

	var raw []byte

	// now transform images as win or scr
	log.GetLogger().Info("Let's go transform images files in win or scr\n")

	_, _, palette, _, err = gfx.ApplyOneImage(images[0], cfg, int(mode), palette, mode)
	if err != nil {
		return nil, nil, palette, err
	}
	for i := 0; i < len(images); i += pad {
		in := images[i]
		raw, _, _, _, err = gfx.ApplyOneImage(in, cfg, int(mode), palette, mode)
		if err != nil {
			return nil, nil, palette, err
		}
		rawImages = append(rawImages, raw)
		log.GetLogger().Info("Image [%d] proceed\n", i)
	}

	lineOctetsWidth := cfg.LineWidth
	x0, y0, err := transformation.CpcCoordinates(initialAddress, 0xC000, lineOctetsWidth)
	if err != nil {
		log.GetLogger().Error("error while computing cpc coordinates :%v\n", err)
	}

	log.GetLogger().Info("Let's go deltapacking raw images\n")
	realSize := &constants.Size{Width: cfg.ScrCfg.Size.Width, Height: cfg.ScrCfg.Size.Height}
	if isSprite {
		realSize.Width = realSize.ModeWidth(mode)
	}
	var lastImage []byte
	for i := 0; i < len(rawImages)-1; i++ {
		log.GetLogger().Info("Compare image [%d] with [%d] ", i, i+1)
		d1 := rawImages[i]
		d2 := rawImages[i+1]
		if len(d1) != len(d2) {
			return nil, nil, palette, errors.ErrorSizeDiffers
		}
		lastImage = d2
		dc := transformation.Delta(d1, d2, isSprite, *realSize, mode, uint16(x0), uint16(y0), lineOctetsWidth)
		deltaData = append(deltaData, dc)
		log.GetLogger().Info("%d bytes differ from the both images\n", len(dc.Items))
	}
	log.GetLogger().Info("Compare image [%d] with [%d] ", len(rawImages)-1, 0)
	d1 := lastImage
	d2 := rawImages[0]
	dc := transformation.Delta(d1, d2, isSprite, *realSize, mode, uint16(x0), uint16(y0), lineOctetsWidth)
	deltaData = append(deltaData, dc)
	log.GetLogger().Info("%d bytes differ from the both images\n", len(dc.Items))
	return deltaData, rawImages, palette, nil
}

// nolint: funlen
func DeltaPacking(gitFilepath string, cfg *config.MartineConfig, initialAddress uint16, mode uint8, exportVersion DeltaExportFormat) error {
	isSprite := true
	maxImages := 22
	if !cfg.CustomDimension && cfg.ScrCfg.Type != config.SpriteHardFormat {
		isSprite = false
	}
	fr, err := os.Open(gitFilepath)
	if err != nil {
		return err
	}
	defer fr.Close()
	gifImages, err := gif.DecodeAll(fr)
	if err != nil {
		return err
	}
	images := ConvertToImage(*gifImages)
	var pad int = 1
	if len(images) <= 1 {
		return fmt.Errorf("need more than one image to proceed")
	}
	if len(images) > maxImages {
		log.GetLogger().Error("Warning gif exceed 30 images. Will corrupt the number of images.")
		pad = len(images) / maxImages
	}
	rawImages := make([][]byte, 0)
	deltaData := make([]*transformation.DeltaCollection, 0)
	var palette color.Palette
	var raw []byte

	// now transform images as win or scr
	log.GetLogger().Info("Let's go transform images files in win or scr\n")

	if cfg.FilloutGif {
		imgs := filloutGif(*gifImages)
		_, _, palette, _, err = gfx.ApplyOneImage(imgs[0], cfg, int(mode), palette, mode)
		if err != nil {
			return err
		}
		for i := 0; i < len(imgs); i += pad {
			in := imgs[i]
			raw, _, _, _, err = gfx.ApplyOneImage(in, cfg, int(mode), palette, mode)
			if err != nil {
				return err
			}
			rawImages = append(rawImages, raw)
			log.GetLogger().Info("Image [%d] proceed\n", i)
		}
	} else {
		_, _, palette, _, err = gfx.ApplyOneImage(images[0], cfg, int(mode), palette, mode)
		if err != nil {
			return err
		}
		for i := 0; i < len(images); i += pad {
			in := images[i]
			raw, _, _, _, err = gfx.ApplyOneImage(in, cfg, int(mode), palette, mode)
			if err != nil {
				return err
			}
			err = png.Png(cfg.ScrCfg.OutputPath+fmt.Sprintf("/%.2d.png", i), in)
			if err != nil {
				return err
			}
			rawImages = append(rawImages, raw)
			log.GetLogger().Info("Image [%d] proceed\n", i)
		}
	}
	lineOctetsWidth := cfg.LineWidth
	x0, y0, err := transformation.CpcCoordinates(initialAddress, 0xC000, lineOctetsWidth)
	if err != nil {
		log.GetLogger().Error("error while computing cpc coordinates :%v\n", err)
	}

	log.GetLogger().Info("Let's go deltapacking raw images\n")
	realSize := &constants.Size{Width: cfg.ScrCfg.Size.Width, Height: cfg.ScrCfg.Size.Height}
	realSize.Width = realSize.ModeWidth(mode)
	var lastImage []byte
	for i := 0; i < len(rawImages)-1; i++ {
		log.GetLogger().Info("Compare image [%d] with [%d] ", i, i+1)
		d1 := rawImages[i]
		d2 := rawImages[i+1]
		if len(d1) != len(d2) {
			return errors.ErrorSizeDiffers
		}
		lastImage = d2
		dc := transformation.Delta(d1, d2, isSprite, *realSize, mode, uint16(x0), uint16(y0), lineOctetsWidth)
		deltaData = append(deltaData, dc)
		log.GetLogger().Info("%d bytes differ from the both images\n", len(dc.Items))
	}
	log.GetLogger().Info("Compare image [%d] with [%d] ", len(rawImages)-1, 0)
	d1 := lastImage
	d2 := rawImages[0]
	dc := transformation.Delta(d1, d2, isSprite, *realSize, mode, uint16(x0), uint16(y0), lineOctetsWidth)
	deltaData = append(deltaData, dc)
	log.GetLogger().Info("%d bytes differ from the both images\n", len(dc.Items))
	filename := cfg.OsFilename(".asm")
	_, err = ExportDeltaAnimate(rawImages[0], deltaData, palette, isSprite, cfg, initialAddress, mode, filepath.Join(cfg.ScrCfg.OutputPath, filename), exportVersion)
	return err
}

func ConvertToImage(g gif.GIF) []*image.NRGBA {
	c := make([]*image.NRGBA, 0)
	imgRect := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: g.Config.Width, Y: g.Config.Height}}
	origImg := image.NewRGBA(imgRect)
	draw.Draw(origImg, g.Image[0].Bounds(), g.Image[0], g.Image[0].Bounds().Min, 0)
	c = append(c, (*image.NRGBA)(origImg))

	previousImg := origImg

	for i := 1; i < len(g.Image); i++ {
		img := image.NewRGBA(imgRect)
		draw.Draw(img, previousImg.Bounds(), previousImg, previousImg.Bounds().Min, draw.Over)
		currImg := g.Image[i]
		draw.Draw(img, currImg.Bounds(), currImg, currImg.Bounds().Min, draw.Over)
		c = append(c, (*image.NRGBA)(img))
		previousImg = img
	}
	return c
}

func filloutGif(g gif.GIF) []image.Image {
	c := make([]image.Image, 0)
	width := g.Image[0].Bounds().Max.X
	height := g.Image[0].Bounds().Max.Y
	reference := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(reference, reference.Bounds(), g.Image[0], image.Point{0, 0}, draw.Src)
	for i := 1; i < len(g.Image)-1; i++ {
		in := g.Image[i]
		draw.Draw(reference, reference.Bounds(), in, image.Point{0, 0}, draw.Over)
		img := image.NewNRGBA(image.Rect(0, 0, width, height))
		draw.Draw(img, img.Bounds(), reference, image.Point{0, 0}, draw.Over)
		c = append(c, img)
	}
	return c
}

// nolint: funlen, gocognit
func ExportDeltaAnimate(
	imageReference []byte,
	delta []*transformation.DeltaCollection,
	palette color.Palette,
	isSprite bool,
	cfg *config.MartineConfig,
	initialAddress uint16,
	mode uint8,
	filename string,
	exportVersion DeltaExportFormat,
) (string, error) {
	an := AnimateValues{
		InitialAddress: fmt.Sprintf("#%.4x", initialAddress),
		Palette:        palette,
		Mode:           int(mode),
		LigneLarge:     fmt.Sprintf("#%.4x", 0xC000+cfg.LineWidth),
		Haut:           fmt.Sprintf("%d", cfg.ScrCfg.Size.Height),
		Large:          fmt.Sprintf("%d", cfg.ScrCfg.Size.ModeWidth(mode)),
		Image:          imageReference,
		Type: AnimateExportType{
			Compress: cfg.ScrCfg.Compression != compression.NONE,
			IsSprite: isSprite,
			CPCPlus:  cfg.ScrCfg.IsPlus,
		},
	}
	data := make([][]byte, 0)
	for _, v := range delta {
		if v.OccurencePerFrame == 0 {
			continue
		}
		if exportVersion == DeltaExportV2 {
			v2 := &transformation.DeltaCollectionV2{DeltaCollection: v}
			d, err := v2.Marshall()
			if err != nil {
				return "", err
			}
			data = append(data, d)
		} else {
			d, err := v.Marshall()
			if err != nil {
				return "", err
			}
			data = append(data, d)
		}
	}
	an.Delta = data

	var sourceCode string

	if !isSprite {
		if cfg.ScrCfg.Compression != compression.NONE {
			sourceCode = deltaScreenCompressCodeDelta
			if cfg.ScrCfg.IsPlus {
				sourceCode = deltaScreenCompressCodeDeltaPlus
			} else {
				if exportVersion == DeltaExportV2 {
					sourceCode = deltaScreenCompressCodeDeltaV2
				}
			}
		} else {
			sourceCode = deltaScreenCodeDelta
			if cfg.ScrCfg.IsPlus {
				sourceCode = deltaScreenCodeDeltaPlus
			} else {
				if exportVersion == DeltaExportV2 {
					sourceCode = deltaScreenCodeDeltaV2
				}
			}
		}
	} else {
		sourceCode = deltaCodeDelta
	}
	var buf bytes.Buffer
	temp := template.Must(template.New("code").Parse(sourceCode))
	err := temp.Execute(&buf, an)
	if err != nil {
		return "", err
	}
	fmt.Println(buf.String())

	code := buf.String()
	if cfg.ScrCfg.Compression != compression.NONE {
		code += "\nbuffer:\n"
	}
	code += "\nend\n"
	code += "\nsave'disc.bin',#200, end - start,DSK,'delta.dsk'"

	if filename != "" {
		err = amsdos.SaveStringOSFile(filename, code)
		if err != nil {
			return "", err
		}
		return code, nil
	}

	return code, nil
}

type AnimateExportType struct {
	Compress bool
	CPCPlus  bool
	IsSprite bool
}

var AnimateTemplates = map[AnimateExportType]string{
	{Compress: false, CPCPlus: false, IsSprite: false}: depackRoutine,
}

type AnimateValues struct {
	Type           AnimateExportType
	InitialAddress string
	Large          string
	Haut           string
	LigneLarge     string
	Mode           int
	Image          []byte
	Delta          [][]byte
	Palette        color.Palette
}

func (a AnimateValues) DisplayCode() string {
	var code string
	var mu sync.Mutex
	if a.Type.Compress {
		log.GetLogger().Info("Using Zx0 cruncher")
		d := zx0.Encode(a.Image)
		mu.Lock()
		code += "sprite:\n"
		code += ascii.FormatAssemblyDatabyte(d, "\n")
		mu.Unlock()
	} else {
		code += "sprite:\n"
		code += ascii.FormatAssemblyDatabyte(a.Image, "\n")
	}
	ascii.ByteToken = "db"
	if a.Type.Compress {
		var wg sync.WaitGroup

		for i, v := range a.Delta {
			wg.Add(1)
			go func(indice int, data []byte) {
				defer wg.Done()
				log.GetLogger().Info("Using Zx0 cruncher")
				d := zx0.Encode(data)
				mu.Lock()
				code += fmt.Sprintf("delta%.2d:\n", indice) // add semaphore on code string
				code += ascii.FormatAssemblyDatabyte(d, "\n")
				mu.Unlock()
			}(i, v)

		}
		wg.Wait()
	} else {
		for i, v := range a.Delta {
			code += fmt.Sprintf("delta%.2d:\n", i)
			code += ascii.FormatAssemblyDatabyte(v, "\n")
		}
	}

	return code
}

func (a AnimateValues) TableDelta() string {
	var code string
	deltaIndexes := make([]string, 0)
	for i := range a.Delta {
		deltaIndexes = append(deltaIndexes, fmt.Sprintf("delta%.2d", i))
	}
	ascii.ByteToken = "dw"
	code += ascii.FormatAssemblyString(deltaIndexes, "\n")
	return code
}

func (a AnimateValues) DisplayPalette() string {
	var code string
	ascii.ByteToken = "db"
	code += "db "
	if a.Type.CPCPlus {
		switch a.Mode {
		case 0:
			code += "#8C\ndb "
		case 1:
			code += "#8D\ndb "
		case 2:
			code += "#8E\ndb "
		}
		code += ascii.FormatAssemblyCPCPlusPalette(a.Palette, "\n")
	} else {
		code += ascii.FormatAssemblyBasicPalette(a.Palette, "\n")
	}
	return code
}

var deltaScreenCodeDeltaV2 = `
;--- dimensions du sprite ----
large equ {{ .Large }}
haut equ {{ .Haut }} 
loadingaddress equ #200
linewidth equ {{ .LigneLarge }}
nbdelta equ {{ .Delta }}
nbcolors equ {{ .Palette }}
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,{{ .Mode }}
	call #BC0E
;-------------------------------

;--- gestion de la palette ----
	call palettefirmware
;------------------------------

call xvbl

;--- affichage du sprite initiale --
	; affichage du premier ecran
	ld de,#C000
	ld hl,sprite
	ldir
;------------------------------------

mainloop    ; routine pour afficher les deltas provenant de martine

;all #bb06

call xvbl
call next_delta

jp mainloop



;--- routine de deltapacking --------------------------
next_delta:
table_index:
	ld a,-1
	inc a
	cp nbdelta
	jr c, table_next
	xor a
table_next:
	ld (table_index+1),a
	add a
	ld e,a
	ld d,0
	ld hl,table_delta
	add hl,de
	ld a,(hl)
	inc hl
	ld h,(hl)
	ld l,a

delta
	ld c,(hl) ; nombre de frame
	inc hl
	ld b,(hl)
	ld (nbdeltas),bc

init
	inc hl
	ld a,(hl) ; octet a poker
	ld (pixel),a
	inc hl
	ld a,(hl) ;

	ld (highbyte_value+1),a ; valeur du HighByte
	inc hl
	ld c,(hl) ; nbfois
	inc hl
	ld b,(hl)
	ld (nblb), bc ; nombre de LowByte

iter_lowbytes
	;
	inc hl
	ld e,(hl) ; recuperation du lowbyte
highbyte_value 	ld d,0

	ld a,(pixel)
	push hl ; on ajoute l'adresse ecran
	ld hl,#c000
	add hl,de
	ld d,h
	ld e,l
	pop hl
	ld (de),a ; poke a l'adresse dans de

	ld bc,(nblb)
	dec bc
	ld (nblb),bc
	ld a,b ; test a t'on poke toutes les lowbytes
	or a
	jr nz, iter_lowbytes
	ld a,c
	or a
	jr nz, iter_lowbytes

	ld bc,(nbdeltas)
	dec bc
	ld (nbdeltas),bc
	ld a,b
	or a
	jr nz,init
	ld a,c
	or a
	jr nz, init

	; a t'on encore des frames a traite


	ret


;---------------------------------------------------------------
;
; attente de plusieurs vbl
;
xvbl ld e,50
	call waitvbl
	dec e
	jr nz,xvbl+2
	ret
;-----------------------------------

` + assembly.WaitVbl + `
` + assembly.FirmwarePalette + `
` + assembly.BC26 + `

;--- variables memoires -----
pixel db 0
nblb dw 0
nbdeltas dw 0
;----------------------------


{{ .DisplayCode }} 

table_delta 
{{ .TableDelta }}

Palette:
{{ .DisplayPalette }}
`

var deltaScreenCompressCodeDeltaPlus = `
;--- dimensions du sprite ----
loadingaddress equ #200
linewidth equ {{ .LigneLarge }}
nbdelta equ {{ len .Delta }}
nbcolors equ {{ len .Palette }}
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,{{ .Mode }}
	call #BC0E
;--- gestion de la palette ----
	DI
	LD	BC,#BC11
	LD	HL,UnlockAsic
Unlock:
	LD	A,(HL)
	OUT	(C),A
	INC	HL
	DEC	C
	JR	NZ,Unlock
	LD BC,#7FB8
	OUT	(C),C
	LD	A,(palette)
	OUT	(C),A

	LD	HL,palette+1
	LD	DE,#6400
	LD	BC,{{ len .Palette }}+{{ len .Palette }}
	LDIR
	LD	BC,#7FA0
	OUT	(C),C
	EI
;------------------------------
jp display_screen


` + assembly.UnlockAsicSequence + `

;---------------------------------------------

Palette:
{{ .DisplayPalette }}

display_screen:
LD HL,#b7f9
call #bcdd

call xvbl

;--- affichage du sprite initiale --
	; affichage du premier ecran
	ld de,#C000
	ld hl,sprite
	call Depack
;------------------------------------
mainloop    ; routine pour afficher les deltas provenant de martine

call xvbl
call next_delta

jp mainloop


;--- routine de deltapacking --------------------------
next_delta:
table_index:
	ld a,-1
	inc a
	cp nbdelta
	jr c, table_next
	xor a
table_next:
	ld (table_index+1),a
	add a
	ld e,a
	ld d,0
	ld hl,table_delta
	add hl,de
	ld a,(hl)
	inc hl
	ld h,(hl)
	ld l,a
	ld de,buffer

	call Depack
	ld hl,buffer ; utilisation de la structure delta décompactée

delta:
	ld a,(hl) ; nombre de byte a poker
	push af   ; stockage en mémoire
	inc hl
init
	ld a,(hl) ; octet a poker
	ld (pixel),a
	inc hl
	ld c,(hl) ; nbfois
	inc hl
	ld b,(hl)
	inc hl
;
poke_octet
	ld e,(hl)
	inc hl
	ld d,(hl) ; de=adresse
	inc hl
	ld a,(pixel)
	push hl ; on ajoute l'adresse ecran
	ld hl,#c000
	add hl,de
	ld d,h
	ld e,l
	pop hl
	ld (de),a ; poke a l'adresse dans de
	dec bc
	ld a,b ; test a t'on poke toutes les adresses compteur bc
	or a
	jr nz, poke_octet
	ld a,c
	or a
	jr nz, poke_octet
	pop af
; reste t'il d'autres bytes a poker ?
	dec a
	push af
	jr nz,init
	pop af
	ret

` + assembly.DeltapackRoutine + `

;---------------------------------------------------------------
;
; attente de plusieurs vbl
;
xvbl ld e,50
	call waitvbl
	dec e
	jr nz,xvbl+2
	ret
;-----------------------------------

` + assembly.WaitVbl + `
` + assembly.BC26 + `


;--- variables memoires -----
pixel db 0
nblb dw 0
nbdeltas dw 0
;----------------------------

{{ .DisplayCode }} 

table_delta 
{{ .TableDelta }}

`

var deltaScreenCodeDeltaPlus = `
;--- dimensions du sprite ----
loadingaddress equ #200
linewidth equ {{ .LigneLarge }}
nbdelta equ {{ len .Delta }}
nbcolors equ {{ len .Palette }}
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,{{ .Mode }}
	call #BC0E

;--- gestion de la palette / unlock asic ----
	DI
	LD	BC,#BC11
	LD	HL,UnlockAsic
Unlock:
	LD	A,(HL)
	OUT	(C),A
	INC	HL
	DEC	C
	JR	NZ,Unlock
	LD BC,#7FB8
	OUT	(C),C
	LD	A,(palette)
	OUT	(C),A

	LD	HL,palette+1
	LD	DE,#6400
	LD	BC,{{ len .Palette }}+{{ len .Palette }}
	LDIR
	LD	BC,#7FA0
	OUT	(C),C
	EI
;------------------------------
jp display_screen

` + assembly.UnlockAsicSequence + `


Palette:
{{ .DisplayPalette }}

display_screen:

call xvbl

;--- affichage du sprite initiale --
	; affichage du premier ecran
	ld de,#C000
	ld hl,sprite
	ldir
;------------------------------------

mainloop    ; routine pour afficher les deltas provenant de martine

call xvbl
call next_delta

jp mainloop


;--- routine de deltapacking --------------------------
next_delta:
table_index:
	ld a,-1
	inc a
	cp nbdelta
	jr c, table_next
	xor a
table_next:
	ld (table_index+1),a
	add a
	ld e,a
	ld d,0
	ld hl,table_delta
	add hl,de
	ld a,(hl)
	inc hl
	ld h,(hl)
	ld l,a
delta
	ld a,(hl) ; nombre de byte a poker
	push af   ; stockage en mémoire
	inc hl
init
	ld a,(hl) ; octet a poker
	ld (pixel),a
	inc hl
	ld c,(hl) ; nbfois
	inc hl
	ld b,(hl)
	inc hl
;
poke_octet
	ld e,(hl)
	inc hl
	ld d,(hl) ; de=adresse
	inc hl
	ld a,(pixel)
	push hl ; on ajoute l'adresse ecran
	ld hl,#c000
	add hl,de
	ld d,h
	ld e,l
	pop hl
	ld (de),a ; poke a l'adresse dans de
	dec bc
	ld a,b ; test a t'on poke toutes les adresses compteur bc
	or a
	jr nz, poke_octet
	ld a,c
	or a
	jr nz, poke_octet
	pop af
; reste t'il d'autres bytes a poker ?
	dec a
	push af
	jr nz,init
	pop af
	ret

;---------------------------------------------------------------
;
; attente de plusieurs vbl
;
xvbl ld e,50
	call waitvbl
	dec e
	jr nz,xvbl+2
	ret
;-----------------------------------

` + assembly.WaitVbl + `
` + assembly.BC26 + `

;--- variables memoires -----
pixel db 0

;----------------------------

{{ .DisplayCode }} 

table_delta 
{{ .TableDelta }}

`

var deltaScreenCompressCodeDeltaV2 = `
;--- dimensions du sprite ----
large equ {{ .Large }}
haut equ {{ .Haut }}
loadingaddress equ #200
linewidth equ {{ .LigneLarge }}
nbdelta equ {{ len .Delta }} 
nbcolors equ {{ len .Palette }}
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,{{ .Mode }}
	call #BC0E
;-------------------------------

;--- gestion de la palette ----
	call palettefirmware
;------------------------------

call xvbl

;--- affichage du sprite initiale --
	; affichage du premier ecran
	ld de,#C000
	ld hl,sprite
	call Depack
;------------------------------------

mainloop    ; routine pour afficher les deltas provenant de martine

;all #bb06

call xvbl
call next_delta

jp mainloop


;--- routine de deltapacking --------------------------
next_delta:
table_index:
	ld a,-1
	inc a
	cp nbdelta
	jr c, table_next
	xor a
table_next:
	ld (table_index+1),a
	add a
	ld e,a
	ld d,0
	ld hl,table_delta
	add hl,de
	ld a,(hl)
	inc hl
	ld h,(hl)
	ld l,a
	ld de,buffer

	call Depack

	ld hl,buffer ; utilisation de la structure delta décompactée
	delta
	ld c,(hl) ; nombre de frame
	inc hl
	ld b,(hl)
	ld (nbdeltas),bc

init
	inc hl
	ld a,(hl) ; octet a poker
	ld (pixel),a
	inc hl
	ld a,(hl) ;

	ld (highbyte_value+1),a ; valeur du HighByte
	inc hl
	ld c,(hl) ; nbfois
	inc hl
	ld b,(hl)
	ld (nblb), bc ; nombre de LowByte

iter_lowbytes
	;
	inc hl
	ld e,(hl) ; recuperation du lowbyte
highbyte_value 	ld d,0

	ld a,(pixel)
	push hl ; on ajoute l'adresse ecran
	ld hl,#c000
	add hl,de
	ld d,h
	ld e,l
	pop hl
	ld (de),a ; poke a l'adresse dans de

	ld bc,(nblb)
	dec bc
	ld (nblb),bc
	ld a,b ; test a t'on poke toutes les lowbytes
	or a
	jr nz, iter_lowbytes
	ld a,c
	or a
	jr nz, iter_lowbytes

	ld bc,(nbdeltas)
	dec bc
	ld (nbdeltas),bc
	ld a,b
	or a
	jr nz,init
	ld a,c
	or a
	jr nz, init

	; a t'on encore des frames a traite


	ret
` + assembly.DeltapackRoutine + `

;---------------------------------------------------------------
;
; attente de plusieurs vbl
;
xvbl ld e,50
	call waitvbl
	dec e
	jr nz,xvbl+2
	ret
;-----------------------------------

` + assembly.WaitVbl + `
` + assembly.FirmwarePalette + `
` + assembly.BC26 + `


;--- variables memoires -----
pixel db 0
nblb dw 0
nbdeltas dw 0
;----------------------------

{{ .DisplayCode }} 

table_delta 
{{ .TableDelta }}

Palette:
{{ .DisplayPalette }}
`

var deltaScreenCodeDelta = `
;--- dimensions du sprite ----
large equ {{ .Large }}
haut equ {{ .Haut }}
loadingaddress equ #200
linewidth equ {{ .LigneLarge }}
nbdelta equ {{ len .Delta }} 
nbcolors equ {{ len .Palette }}
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,{{ .Mode }}
	call #BC0E
;-------------------------------

;--- gestion de la palette ----
	call palettefirmware
;------------------------------

call xvbl

;--- affichage du sprite initiale --
	; affichage du premier ecran
	ld de,#C000
	ld hl,sprite
	ldir
;------------------------------------

mainloop    ; routine pour afficher les deltas provenant de martine

;all #bb06

call xvbl
call next_delta

jp mainloop


;--- routine de deltapacking --------------------------
next_delta:
table_index:
	ld a,-1
	inc a
	cp nbdelta
	jr c, table_next
	xor a
table_next:
	ld (table_index+1),a
	add a
	ld e,a
	ld d,0
	ld hl,table_delta
	add hl,de
	ld a,(hl)
	inc hl
	ld h,(hl)
	ld l,a
delta
	ld a,(hl) ; nombre de byte a poker
	push af   ; stockage en mémoire
	inc hl
init
	ld a,(hl) ; octet a poker
	ld (pixel),a
	inc hl
	ld c,(hl) ; nbfois
	inc hl
	ld b,(hl)
	inc hl
;
poke_octet
	ld e,(hl)
	inc hl
	ld d,(hl) ; de=adresse
	inc hl
	ld a,(pixel)
	push hl ; on ajoute l'adresse ecran
	ld hl,#c000
	add hl,de
	ld d,h
	ld e,l
	pop hl
	ld (de),a ; poke a l'adresse dans de
	dec bc
	ld a,b ; test a t'on poke toutes les adresses compteur bc
	or a
	jr nz, poke_octet
	ld a,c
	or a
	jr nz, poke_octet
	pop af
; reste t'il d'autres bytes a poker ?
	dec a
	push af
	jr nz,init
	pop af
	ret

;---------------------------------------------------------------
;
; attente de plusieurs vbl
;
xvbl ld e,50
	call waitvbl
	dec e
	jr nz,xvbl+2
	ret
;-----------------------------------

` + assembly.WaitVbl + `
` + assembly.FirmwarePalette + `
` + assembly.BC26 + `


;--- variables memoires -----
pixel db 0

;----------------------------

{{ .DisplayCode }} 

table_delta 
{{ .TableDelta }}

Palette:
{{ .DisplayPalette }}
`

var deltaScreenCompressCodeDelta string = `
;--- dimensions du sprite ----
large equ {{ .Large }}
haut equ {{ .Haut }}
loadingaddress equ #200
linewidth equ {{ .LigneLarge }}
nbdelta equ {{ len .Delta }} 
nbcolors equ {{ len .Palette }}
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,{{ .Mode }}
	call #BC0E
;-------------------------------

;--- gestion de la palette ----
	call palettefirmware
;------------------------------

call xvbl

;--- affichage du sprite initiale --
	; affichage du premier ecran
	ld de,#C000
	ld hl,sprite
	call Depack
;------------------------------------

mainloop    ; routine pour afficher les deltas provenant de martine

;all #bb06

call xvbl
call next_delta

jp mainloop


;--- routine de deltapacking --------------------------
next_delta:
table_index:
	ld a,-1
	inc a
	cp nbdelta
	jr c, table_next
	xor a
table_next:
	ld (table_index+1),a
	add a
	ld e,a
	ld d,0
	ld hl,table_delta
	add hl,de
	ld a,(hl)
	inc hl
	ld h,(hl)
	ld l,a
	ld de,buffer

	call Depack

	ld hl,buffer ; utilisation de la structure delta décompactée

delta
	ld a,(hl) ; nombre de byte a poker
	push af   ; stockage en mémoire
	inc hl
init
	ld a,(hl) ; octet a poker
	ld (pixel),a
	inc hl
	ld c,(hl) ; nbfois
	inc hl
	ld b,(hl)
	inc hl
;
poke_octet
	ld e,(hl)
	inc hl
	ld d,(hl) ; de=adresse
	inc hl
	ld a,(pixel)
	push hl ; on ajoute l'adresse ecran
	ld hl,#c000
	add hl,de
	ld d,h
	ld e,l
	pop hl
	ld (de),a ; poke a l'adresse dans de
	dec bc
	ld a,b ; test a t'on poke toutes les adresses compteur bc
	or a
	jr nz, poke_octet
	ld a,c
	or a
	jr nz, poke_octet
	pop af
; reste t'il d'autres bytes a poker ?
	dec a
	push af
	jr nz,init
	pop af
	ret


` + assembly.DeltapackRoutine + `

;---------------------------------------------------------------
;
; attente de plusieurs vbl
;
xvbl ld e,50
	call waitvbl
	dec e
	jr nz,xvbl+2
	ret
;-----------------------------------

` + assembly.WaitVbl + `
` + assembly.FirmwarePalette + `

` + assembly.BC26 + `


;--- variables memoires -----
pixel db 0

;----------------------------


{{ .DisplayCode }} 

table_delta 
{{ .TableDelta }}

Palette:
{{ .DisplayPalette }}
`

var deltaCodeDelta string = `;--- dimensions du sprite ----
large equ {{ .Large }}
haut equ {{ .Haut }}
loadingaddress equ #200
linewidth equ {{ .LigneLarge }}
nbdelta equ {{ len .Delta }} 
nbcolors equ {{ len .Palette }}
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,{{ .Mode }}
	call #BC0E
;-------------------------------

;--- gestion de la palette ----
	call palettefirmware
;------------------------------

call xvbl

;--- affichage du sprite initiale --
	; affichage du premier sprite
	ld de,{{ .InitialAddress }} ; adresse de l'ecran
	ld hl,sprite ; pointeur sur l'image en memoire
	ld b, haut ; hauteur de l'image
	loop
	push bc ; sauve le compteur hauteur dans la pile
	push de ; sauvegarde de l'adresse ecran dans la pile
	ld bc, large ; largeur de l'image a afficher
	ldir ; remplissage de n * largeur octets a l'adresse dans de
	pop de ; recuperation de l'adresse d'origine
	ex de,hl ; echange des valeurs des adresses
	call bc26 ; calcul de l'adresse de la ligne suivante
	ex de,hl ; echange des valeurs des adresses
	pop bc ; retabli le compteur
	djnz loop
;------------------------------------

mainloop    ; routine pour afficher les deltas provenant de martine

;call #bb06

call xvbl
call next_delta

jp mainloop


;--- routine de deltapacking --------------------------
next_delta:
table_index:
	ld a,-1
	inc a
	cp nbdelta
	jr c, table_next
	xor a
table_next:
	ld (table_index+1),a
	add a
	ld e,a
	ld d,0
	ld hl,table_delta
	add hl,de
	ld a,(hl)
	inc hl
	ld h,(hl)
	ld l,a
delta
	ld a,(hl) ; nombre de byte a poker
	push af   ; stockage en mémoire
	inc hl
init
	ld a,(hl) ; octet a poker
	ld (pixel),a
	inc hl
	ld c,(hl) ; nbfois
	inc hl
	ld b,(hl)
	inc hl
;
poke_octet
	ld e,(hl)
	inc hl
	ld d,(hl) ; de=adresse
	inc hl
	ld a,(pixel)
	ld (de),a ; poke a l'adresse dans de
	dec bc
	ld a,b ; test a t'on poke toutes les adresses compteur bc
	or a
	jr nz, poke_octet
	ld a,c
	or a
	jr nz, poke_octet
	pop af
; reste t'il d'autres bytes a poker ?
	dec a
	push af
	jr nz,init
	pop af
	ret

;---------------------------------------------------------------
;
; attente de plusieurs vbl
;
xvbl ld e,50
	call waitvbl
	dec e
	jr nz,xvbl+2
	ret
;-----------------------------------

` + assembly.WaitVbl + `
` + assembly.FirmwarePalette + `
` + assembly.BC26 + `


;--- variables memoires -----
pixel db 0
;----------------------------

{{ .DisplayCode }} 

table_delta 
{{ .TableDelta }}

Palette:
{{ .DisplayPalette }}
`

var depackRoutine = `
;--- dimensions du sprite ----
large equ {{ .Large }}
haut equ {{ .Haut }}
loadingaddress equ #200
linewidth equ {{ .LigneLarge }}
nbdelta equ {{ len .Delta }}
nbcolors equ {{ len .Palette }}
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,{{ .Mode }}
	call #BC0E
;-------------------------------

;--- gestion de la palette ----
	call palettefirmware
;------------------------------

call xvbl

;--- affichage du sprite initiale --
	; affichage du premier sprite
	ld de,buffer
	ld hl,sprite
	call Depack

	ld de, {{ .InitialAddress }} ; adresse de l'ecran
	ld hl,buffer ; pointeur sur l'image en memoire
	ld b, haut ; hauteur de l'image
	loop
	push bc ; sauve le compteur hauteur dans la pile
	push de ; sauvegarde de l'adresse ecran dans la pile
	ld bc, large ; largeur de l'image a afficher
	ldir ; remplissage de n * largeur octets a l'adresse dans de
	pop de ; recuperation de l'adresse d'origine
	ex de,hl ; echange des valeurs des adresses
	call bc26 ; calcul de l'adresse de la ligne suivante
	ex de,hl ; echange des valeurs des adresses
	pop bc ; retabli le compteur
	djnz loop
;------------------------------------

mainloop    ; routine pour afficher les deltas provenant de martine

;call #bb06

call xvbl
call next_delta

jp mainloop


;--- routine de deltapacking --------------------------
next_delta:
table_index:
	ld a,-1
	inc a
	cp nbdelta
	jr c, table_next
	xor a
table_next:
	ld (table_index+1),a
	add a
	ld e,a
	ld d,0
	ld hl,table_delta
	add hl,de
	ld a,(hl)
	inc hl
	ld h,(hl)
	ld l,a
	ld de,buffer

	call Depack

	ld hl,buffer ; utilisation de la structure delta décompactée

delta
	ld a,(hl) ; nombre de byte a poker
	push af   ; stockage en mémoire
	inc hl
init
	ld a,(hl) ; octet a poker
	ld (pixel),a
	inc hl
	ld c,(hl) ; nbfois
	inc hl
	ld b,(hl)
	inc hl
;
poke_octet
	ld e,(hl)
	inc hl
	ld d,(hl) ; de=adresse
	inc hl
	ld a,(pixel)
	ld (de),a ; poke a l'adresse dans de
	dec bc
	ld a,b ; test a t'on poke toutes les adresses compteur bc
	or a
	jr nz, poke_octet
	ld a,c
	or a
	jr nz, poke_octet
	pop af
; reste t'il d'autres bytes a poker ?
	dec a
	push af
	jr nz,init
	pop af
	ret



` + assembly.DeltapackRoutine + `

;---------------------------------------------------------------
;
; attente de plusieurs vbl
;
xvbl ld e,50
	call waitvbl
	dec e
	jr nz,xvbl+2
	ret
;-----------------------------------

` + assembly.WaitVbl + `
` + assembly.FirmwarePalette + `
` + assembly.BC26 + `


;--- variables memoires -----
pixel db 0

;----------------------------
{{ .DisplayCode }}
{{ .TableDelta }}
{{ .DisplayPalette }}
buffer: 
`
