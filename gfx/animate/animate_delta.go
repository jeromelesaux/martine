package animate

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export/amsdos"
	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/export/png"
	"github.com/jeromelesaux/martine/gfx"
	"github.com/jeromelesaux/martine/gfx/errors"
	"github.com/jeromelesaux/martine/gfx/transformation"
	zx0 "github.com/jeromelesaux/zx0/encode"
)

type DeltaExportFormat int

var (
	DeltaExportV1 DeltaExportFormat = 1
	DeltaExportV2 DeltaExportFormat = 2
)

func DeltaPackingMemory(images []image.Image, ex *config.MartineConfig, initialAddress uint16, mode uint8) ([]*transformation.DeltaCollection, [][]byte, color.Palette, error) {
	var isSprite bool = true
	var maxImages = 22
	var pad int = 1
	var err error
	var palette color.Palette
	if !ex.CustomDimension && !ex.SpriteHard {
		isSprite = false
	}
	if len(images) <= 1 {
		return nil, nil, palette, fmt.Errorf("need more than one image to proceed")
	}
	if len(images) > maxImages {
		fmt.Fprintf(os.Stderr, "Warning gif exceed 30 images. Will corrupt the number of images.")
		pad = len(images) / maxImages
	}
	rawImages := make([][]byte, 0)
	deltaData := make([]*transformation.DeltaCollection, 0)

	var raw []byte

	// now transform images as win or scr
	fmt.Printf("Let's go transform images files in win or scr\n")

	_, _, palette, _, err = gfx.ApplyOneImage(images[0], ex, int(mode), palette, mode)
	if err != nil {
		return nil, nil, palette, err
	}
	for i := 0; i < len(images); i += pad {
		in := images[i]
		raw, _, _, _, err = gfx.ApplyOneImage(in, ex, int(mode), palette, mode)
		if err != nil {
			return nil, nil, palette, err
		}
		rawImages = append(rawImages, raw)
		fmt.Printf("Image [%d] proceed\n", i)
	}

	lineOctetsWidth := ex.LineWidth
	x0, y0, err := transformation.CpcCoordinates(initialAddress, 0xC000, lineOctetsWidth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while computing cpc coordinates :%v\n", err)
	}

	fmt.Printf("Let's go deltapacking raw images\n")
	realSize := &constants.Size{Width: ex.Size.Width, Height: ex.Size.Height}
	if isSprite {
		realSize.Width = realSize.ModeWidth(mode)
	}
	var lastImage []byte
	for i := 0; i < len(rawImages)-1; i++ {
		fmt.Printf("Compare image [%d] with [%d] ", i, i+1)
		d1 := rawImages[i]
		d2 := rawImages[i+1]
		if len(d1) != len(d2) {
			return nil, nil, palette, errors.ErrorSizeDiffers
		}
		lastImage = d2
		dc := transformation.Delta(d1, d2, isSprite, *realSize, mode, uint16(x0), uint16(y0), lineOctetsWidth)
		deltaData = append(deltaData, dc)
		fmt.Printf("%d bytes differ from the both images\n", len(dc.Items))
	}
	fmt.Printf("Compare image [%d] with [%d] ", len(rawImages)-1, 0)
	d1 := lastImage
	d2 := rawImages[0]
	dc := transformation.Delta(d1, d2, isSprite, *realSize, mode, uint16(x0), uint16(y0), lineOctetsWidth)
	deltaData = append(deltaData, dc)
	fmt.Printf("%d bytes differ from the both images\n", len(dc.Items))
	return deltaData, rawImages, palette, nil
}

func DeltaPacking(gitFilepath string, ex *config.MartineConfig, initialAddress uint16, mode uint8, exportVersion DeltaExportFormat) error {
	var isSprite = true
	var maxImages = 22
	if !ex.CustomDimension && !ex.SpriteHard {
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
		fmt.Fprintf(os.Stderr, "Warning gif exceed 30 images. Will corrupt the number of images.")
		pad = len(images) / maxImages
	}
	rawImages := make([][]byte, 0)
	deltaData := make([]*transformation.DeltaCollection, 0)
	var palette color.Palette
	var raw []byte

	// now transform images as win or scr
	fmt.Printf("Let's go transform images files in win or scr\n")

	if ex.FilloutGif {
		imgs := filloutGif(*gifImages, ex)
		_, _, palette, _, err = gfx.ApplyOneImage(imgs[0], ex, int(mode), palette, mode)
		if err != nil {
			return err
		}
		for i := 0; i < len(imgs); i += pad {
			in := imgs[i]
			raw, _, _, _, err = gfx.ApplyOneImage(in, ex, int(mode), palette, mode)
			if err != nil {
				return err
			}
			rawImages = append(rawImages, raw)
			fmt.Printf("Image [%d] proceed\n", i)
		}
	} else {
		_, _, palette, _, err = gfx.ApplyOneImage(images[0], ex, int(mode), palette, mode)
		if err != nil {
			return err
		}
		for i := 0; i < len(images); i += pad {
			in := images[i]
			raw, _, _, _, err = gfx.ApplyOneImage(in, ex, int(mode), palette, mode)
			if err != nil {
				return err
			}
			png.Png(ex.OutputPath+fmt.Sprintf("/%.2d.png", i), in)

			rawImages = append(rawImages, raw)
			fmt.Printf("Image [%d] proceed\n", i)
		}
	}
	lineOctetsWidth := ex.LineWidth
	x0, y0, err := transformation.CpcCoordinates(initialAddress, 0xC000, lineOctetsWidth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while computing cpc coordinates :%v\n", err)
	}

	fmt.Printf("Let's go deltapacking raw images\n")
	realSize := &constants.Size{Width: ex.Size.Width, Height: ex.Size.Height}
	realSize.Width = realSize.ModeWidth(mode)
	var lastImage []byte
	for i := 0; i < len(rawImages)-1; i++ {
		fmt.Printf("Compare image [%d] with [%d] ", i, i+1)
		d1 := rawImages[i]
		d2 := rawImages[i+1]
		if len(d1) != len(d2) {
			return errors.ErrorSizeDiffers
		}
		lastImage = d2
		dc := transformation.Delta(d1, d2, isSprite, *realSize, mode, uint16(x0), uint16(y0), lineOctetsWidth)
		deltaData = append(deltaData, dc)
		fmt.Printf("%d bytes differ from the both images\n", len(dc.Items))
	}
	fmt.Printf("Compare image [%d] with [%d] ", len(rawImages)-1, 0)
	d1 := lastImage
	d2 := rawImages[0]
	dc := transformation.Delta(d1, d2, isSprite, *realSize, mode, uint16(x0), uint16(y0), lineOctetsWidth)
	deltaData = append(deltaData, dc)
	fmt.Printf("%d bytes differ from the both images\n", len(dc.Items))
	filename := string(ex.OsFilename(".asm"))
	return exportDeltaAnimate(rawImages[0], deltaData, palette, isSprite, ex, initialAddress, mode, ex.OutputPath+string(filepath.Separator)+filename, exportVersion)
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

func filloutGif(g gif.GIF, ex *config.MartineConfig) []image.Image {
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

func ExportDeltaAnimate(imageReference []byte, delta []*transformation.DeltaCollection, palette color.Palette, isSprite bool, ex *config.MartineConfig, initialAddress uint16, mode uint8, exportVersion DeltaExportFormat) (string, error) {
	var sourceCode string
	var dataCode string
	var deltaIndex []string
	var code string
	var nbDelta = len(delta)
	if !isSprite {
		if ex.Compression != -1 {

			sourceCode = deltaScreenCompressCodeDelta
			if exportVersion == DeltaExportV2 {
				sourceCode = deltaScreenCompressCodeDeltaV2
			}
			if ex.CpcPlus {
				sourceCode = deltaScreenCompressCodeDeltaPlus
			}
		} else {
			sourceCode = deltaScreenCodeDelta
			if exportVersion == DeltaExportV2 {
				sourceCode = deltaScreenCodeDeltaV2
			}
			if ex.CpcPlus {
				sourceCode = deltaScreenCodeDeltaPlus
			}
		}
	} else {
		sourceCode = deltaCodeDelta
	}
	// copy of the sprite
	dataCode += "\nsprite:\n"
	if ex.Compression != -1 {
		if isSprite {
			sourceCode = depackRoutine
		}
		fmt.Fprintf(os.Stdout, "Using Zx0 cruncher")
		data := zx0.Encode(imageReference)
		dataCode += ascii.FormatAssemblyDatabyte(data, "\n")
	} else {
		dataCode += ascii.FormatAssemblyDatabyte(imageReference, "\n")
	}
	// copy of all delta
	for i := 0; i < len(delta); i++ {
		dc := delta[i]
		var data []byte
		var err error
		if exportVersion == DeltaExportV1 {
			data, err = dc.Marshall()
			if err != nil {
				return "", err
			}
		} else {
			v2 := transformation.DeltaCollectionV2{DeltaCollection: dc}
			data, err = v2.Marshall()
			if err != nil {
				return "", err
			}
		}
		name := fmt.Sprintf("delta%.2d", i)
		dataCode += name + ":\n"
		if ex.Compression != -1 {
			fmt.Fprintf(os.Stdout, "Using Zx0 cruncher")
			if dc.OccurencePerFrame != 0 {
				d := zx0.Encode(data)
				dataCode += ascii.FormatAssemblyDatabyte(d, "\n")
			} else {
				nbDelta--
			}
		} else {
			dataCode += ascii.FormatAssemblyDatabyte(data, "\n")
		}
		deltaIndex = append(deltaIndex, name)
	}
	dataCode += "table_delta:\n"
	ascii.ByteToken = "dw"
	dataCode += ascii.FormatAssemblyString(deltaIndex, "\n")

	ascii.ByteToken = "db"
	dataCode += "palette:\n" + ascii.ByteToken + " "
	dataCode += ascii.FormatAssemblyBasicPalette(palette, "\n")

	header := sourceCode
	// replace the initial address
	if isSprite {
		address := fmt.Sprintf("#%.4x", initialAddress)
		header = strings.Replace(sourceCode, "$INITIALADDRESS$", address, 1)
	}
	// replace number of colors
	nbColors := fmt.Sprintf("%d", len(palette))
	header = strings.Replace(header, "$NBCOLORS$", nbColors, 1)

	// replace the number of delta
	nbDeltaLabel := fmt.Sprintf("%d", nbDelta)
	header = strings.Replace(header, "$NBDELTA$", nbDeltaLabel, 1)

	// replace char large for the screen
	charLarge := fmt.Sprintf("#%.4x", 0xC000+ex.LineWidth)
	header = strings.Replace(header, "$LIGNELARGE$", charLarge, 1)

	// replace heigth
	height := fmt.Sprintf("%d", ex.Size.Height)
	header = strings.Replace(header, "$HAUT$", height, 1)

	// replace width
	var width string = fmt.Sprintf("%d", ex.Size.ModeWidth(mode))
	header = strings.Replace(header, "$LARGE$", width, 1)

	var modeSet string
	switch mode {
	case 0:
		modeSet = "0"
	case 1:
		modeSet = "1"
	case 2:
		modeSet = "2"
	}

	// replace mode
	header = strings.Replace(header, "$SETMODE$", modeSet, 1)

	code += header
	code += dataCode
	if ex.Compression != -1 {
		code += "\nbuffer:\n"
	}
	code += "\nend\n"
	code += "\nsave'disc.bin',#200, end - start,DSK,'martine-animate.dsk'"

	return code, nil
}

func exportDeltaAnimate(
	imageReference []byte,
	delta []*transformation.DeltaCollection,
	palette color.Palette,
	isSprite bool,
	ex *config.MartineConfig,
	initialAddress uint16,
	mode uint8,
	filename string,
	exportVersion DeltaExportFormat) error {

	var sourceCode string = deltaCodeDelta
	var dataCode string
	var deltaIndex []string
	var code string
	nbDelta := len(delta)
	if exportVersion == DeltaExportV2 {
		sourceCode = deltaScreenCodeDeltaV2
	}
	if !isSprite {
		if ex.Compression != -1 {
			sourceCode = deltaScreenCompressCodeDelta
			if exportVersion == DeltaExportV2 {
				sourceCode = deltaScreenCompressCodeDeltaV2
			}
		} else {
			sourceCode = deltaScreenCodeDelta
			if exportVersion == DeltaExportV2 {
				sourceCode = deltaScreenCodeDeltaV2
			}
		}
	}
	// copy of the sprite
	dataCode += "sprite:\n"
	if ex.Compression != -1 {
		if isSprite {
			sourceCode = depackRoutine
		}
		fmt.Fprintf(os.Stdout, "Using Zx0 cruncher")
		data := zx0.Encode(imageReference)
		dataCode += ascii.FormatAssemblyDatabyte(data, "\n")
	} else {
		dataCode += ascii.FormatAssemblyDatabyte(imageReference, "\n")
	}
	// copy of all delta
	for i := 0; i < len(delta); i++ {
		dc := delta[i]
		var data []byte
		var err error
		if exportVersion == DeltaExportV2 {
			v2 := &transformation.DeltaCollectionV2{DeltaCollection: dc}
			data, err = v2.Marshall()
			if err != nil {
				return err
			}
		} else {
			data, err = dc.Marshall()
			if err != nil {
				return err
			}
		}
		name := fmt.Sprintf("delta%.2d", i)
		dataCode += name + ":\n"
		if ex.Compression != -1 {
			fmt.Fprintf(os.Stdout, "Using Zx0 cruncher")
			if dc.OccurencePerFrame != 0 {
				d := zx0.Encode(data)
				dataCode += ascii.FormatAssemblyDatabyte(d, "\n")
				dataCode += ascii.FormatAssemblyDatabyte(d, "\n")
			} else {
				nbDelta--
			}
		} else {
			dataCode += ascii.FormatAssemblyDatabyte(data, "\n")
		}
		deltaIndex = append(deltaIndex, name)
	}
	dataCode += "table_delta:\n"
	ascii.ByteToken = "dw"
	dataCode += ascii.FormatAssemblyString(deltaIndex, "\n")

	ascii.ByteToken = "db"
	dataCode += "palette:\n" + ascii.ByteToken + " "
	dataCode += ascii.FormatAssemblyBasicPalette(palette, "\n")

	var header string = sourceCode
	if isSprite {
		// replace the initial address
		address := fmt.Sprintf("#%.4x", initialAddress)
		header = strings.Replace(sourceCode, "$INITIALADDRESS$", address, 1)
	}
	// replace number of colors
	nbColors := fmt.Sprintf("%d", len(palette))
	header = strings.Replace(header, "$NBCOLORS$", nbColors, 1)

	// replace the number of delta
	// replace the number of delta
	nbDeltaLabel := fmt.Sprintf("%d", nbDelta)
	header = strings.Replace(header, "$NBDELTA$", nbDeltaLabel, 1)

	// replace char large for the screen
	charLarge := fmt.Sprintf("#%.4x", 0xC000+ex.LineWidth)
	header = strings.Replace(header, "$LIGNELARGE$", charLarge, 1)

	// replace heigth
	height := fmt.Sprintf("%d", ex.Size.Height)
	header = strings.Replace(header, "$HAUT$", height, 1)

	// replace width
	var width string = fmt.Sprintf("%d", ex.Size.ModeWidth(mode))
	header = strings.Replace(header, "$LARGE$", width, 1)

	var modeSet string
	switch mode {
	case 0:
		modeSet = "0"
	case 1:
		modeSet = "1"
	case 2:
		modeSet = "2"
	}

	// replace mode
	header = strings.Replace(header, "$SETMODE$", modeSet, 1)

	code += header
	code += dataCode
	code += "\nend\n"
	code += "\nsave'disc.bin',#200, end - start,DSK,'delta.dsk'"
	if ex.Compression != -1 {
		code += "\nbuffer dw 0\n"
	}

	return amsdos.SaveStringOSFile(filename, code)
}

var deltaScreenCodeDeltaV2 = `
;--- dimensions du sprite ----
large equ $LARGE$
haut equ $HAUT$
loadingaddress equ #200
linewidth equ $LIGNELARGE$
nbdelta equ $NBDELTA$
nbcolors equ $NBCOLORS$
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,$SETMODE$
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

;---- attente vbl ----------
waitvbl
	ld b,#f5 ; attente vbl
vbl     
	in a,(c)
	rra
	jp nc,vbl
	ret
;---------------------------

;--- application palette firmware -------------
palettefirmware ; hl pointe sur les valeurs de la palette
ld e,nbcolors
ld a,0
ld hl,palette

paletteloop
ld b,(hl)
ld c,b
push af
push de
push hl
call #bc32 ; af, de, hl corrupted
pop hl
pop de
pop af
inc a
inc hl
dec e
jr nz,paletteloop
ret
;---------------------------------------------

;---------------------------------------------

;---- recuperation de l'adresse de la ligne en dessous ------------
bc26 
ld a,h
add a,8 
ld h,a ; <---- le fameux que tu as oublié !
ret nc 
ld bc,linewidth ; on passe en 96 colonnes
add hl,bc
res 3,h
ret
;-----------------------------------------------------------------

;--- variables memoires -----
pixel db 0
nblb dw 0
nbdeltas dw 0
;----------------------------



`

var deltaScreenCompressCodeDeltaPlus = `
;--- dimensions du sprite ----
large equ $LARGE$
haut equ $HAUT$
loadingaddress equ #200
linewidth equ $LIGNELARGE$
nbdelta equ $NBDELTA$
nbcolors equ $NBCOLORS$
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,$SETMODE$
	call #BC0E
;-------------------------------

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
	LD	BC,#7FA0
	LD	A,(palette)
	OUT	(C),A
	OUT	(C),C
	LD	BC,#7FB8
	OUT	(C),C
	LD	HL,palette+1
	LD	DE,#6400
	LD	BC,#0022
	LDIR
	EI
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

	;
	; Decompactage ZX0
	; HL = source
	; DE = destination
	;
	Depack:
		ld    bc,#ffff        ; preserve default offset 1
		push    bc
		inc    bc
		ld    a,#80
	dzx0s_literals:
		call    dzx0s_elias        ; obtain length
		ldir                ; copy literals
		add    a,a            ; copy from last offset or new offset?
		jr    c,dzx0s_new_offset
		call    dzx0s_elias        ; obtain length
	dzx0s_copy:
		ex    (sp),hl            ; preserve source,restore offset
		push    hl            ; preserve offset
		add    hl,de            ; calculate destination - offset
		ldir                ; copy from offset
		pop    hl            ; restore offset
		ex    (sp),hl            ; preserve offset,restore source
		add    a,a            ; copy from literals or new offset?
		jr    nc,dzx0s_literals
	dzx0s_new_offset:
		call    dzx0s_elias        ; obtain offset MSB
		ld b,a
		pop    af            ; discard last offset
		xor    a            ; adjust for negative offset
		sub    c
		RET    Z            ; Plus d'octets a traiter = fini
	
		ld    c,a
		ld    a,b
		ld    b,c
		ld    c,(hl)            ; obtain offset LSB
		inc    hl
		rr    b            ; last offset bit becomes first length bit
		rr    c
		push    bc            ; preserve new offset
		ld    bc,1            ; obtain length
		call    nc,dzx0s_elias_backtrack
		inc    bc
		jr    dzx0s_copy
	dzx0s_elias:
		inc    c            ; interlaced Elias gamma coding
	dzx0s_elias_loop:
		add    a,a
		jr    nz,dzx0s_elias_skip
		ld    a,(hl)            ; load another group of 8 bits
		inc    hl
		rla
	dzx0s_elias_skip:
		ret     c
	dzx0s_elias_backtrack:
		add    a,a
		rl    c
		rl    b
		jr    dzx0s_elias_loop
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

;---- attente vbl ----------
waitvbl
	ld b,#f5 ; attente vbl
vbl     
	in a,(c)
	rra
	jp nc,vbl
	ret
;---------------------------

UnlockAsic:
	DB	#FF,#00,#FF,#77,#B3,#51,#A8,#D4,#62,#39,#9C,#46,#2B,#15,#8A,#CD,#EE

;---------------------------------------------

;---- recuperation de l'adresse de la ligne en dessous ------------
bc26 
ld a,h
add a,8 
ld h,a ; <---- le fameux que tu as oublié !
ret nc 
ld bc,linewidth ; on passe en 96 colonnes
add hl,bc
res 3,h
ret
;-----------------------------------------------------------------


;--- variables memoires -----
pixel db 0
nblb dw 0
nbdeltas dw 0
;--`

var deltaScreenCodeDeltaPlus = `
;--- dimensions du sprite ----
large equ $LARGE$
haut equ $HAUT$
loadingaddress equ #200
linewidth equ $LIGNELARGE$
nbdelta equ $NBDELTA$
nbcolors equ $NBCOLORS$
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,$SETMODE$
	call #BC0E
;-------------------------------

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
	LD	BC,#7FA0
	LD	A,(palette)
	OUT	(C),A
	OUT	(C),C
	LD	BC,#7FB8
	OUT	(C),C
	LD	HL,palette+1
	LD	DE,#6400
	LD	BC,#0022
	LDIR
	EI
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

;---- attente vbl ----------
waitvbl
	ld b,#f5 ; attente vbl
vbl     
	in a,(c)
	rra
	jp nc,vbl
	ret
;---------------------------

;---------------------------------------------

;---- recuperation de l'adresse de la ligne en dessous ------------
bc26 
ld a,h
add a,8 
ld h,a ; <---- le fameux que tu as oublié !
ret nc 
ld bc,linewidth ; on passe en 96 colonnes
add hl,bc
res 3,h
ret
;-----------------------------------------------------------------


UnlockAsic:
	DB	#FF,#00,#FF,#77,#B3,#51,#A8,#D4,#62,#39,#9C,#46,#2B,#15,#8A,#CD,#EE

;--- variables memoires -----
pixel db 0 

;----------------------------

`

var deltaScreenCompressCodeDeltaV2 = `
;--- dimensions du sprite ----
large equ $LARGE$
haut equ $HAUT$
loadingaddress equ #200
linewidth equ $LIGNELARGE$
nbdelta equ $NBDELTA$
nbcolors equ $NBCOLORS$
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,$SETMODE$
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

	;
	; Decompactage ZX0
	; HL = source
	; DE = destination
	;
	Depack:
		ld    bc,#ffff        ; preserve default offset 1
		push    bc
		inc    bc
		ld    a,#80
	dzx0s_literals:
		call    dzx0s_elias        ; obtain length
		ldir                ; copy literals
		add    a,a            ; copy from last offset or new offset?
		jr    c,dzx0s_new_offset
		call    dzx0s_elias        ; obtain length
	dzx0s_copy:
		ex    (sp),hl            ; preserve source,restore offset
		push    hl            ; preserve offset
		add    hl,de            ; calculate destination - offset
		ldir                ; copy from offset
		pop    hl            ; restore offset
		ex    (sp),hl            ; preserve offset,restore source
		add    a,a            ; copy from literals or new offset?
		jr    nc,dzx0s_literals
	dzx0s_new_offset:
		call    dzx0s_elias        ; obtain offset MSB
		ld b,a
		pop    af            ; discard last offset
		xor    a            ; adjust for negative offset
		sub    c
		RET    Z            ; Plus d'octets a traiter = fini
	
		ld    c,a
		ld    a,b
		ld    b,c
		ld    c,(hl)            ; obtain offset LSB
		inc    hl
		rr    b            ; last offset bit becomes first length bit
		rr    c
		push    bc            ; preserve new offset
		ld    bc,1            ; obtain length
		call    nc,dzx0s_elias_backtrack
		inc    bc
		jr    dzx0s_copy
	dzx0s_elias:
		inc    c            ; interlaced Elias gamma coding
	dzx0s_elias_loop:
		add    a,a
		jr    nz,dzx0s_elias_skip
		ld    a,(hl)            ; load another group of 8 bits
		inc    hl
		rla
	dzx0s_elias_skip:
		ret     c
	dzx0s_elias_backtrack:
		add    a,a
		rl    c
		rl    b
		jr    dzx0s_elias_loop
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

;---- attente vbl ----------
waitvbl
	ld b,#f5 ; attente vbl
vbl     
	in a,(c)
	rra
	jp nc,vbl
	ret
;---------------------------

;--- application palette firmware -------------
palettefirmware ; hl pointe sur les valeurs de la palette
ld e,nbcolors
ld a,0
ld hl,palette

paletteloop
ld b,(hl)
ld c,b
push af
push de
push hl
call #bc32 ; af, de, hl corrupted
pop hl
pop de
pop af
inc a
inc hl
dec e
jr nz,paletteloop
ret
;---------------------------------------------

;---------------------------------------------

;---- recuperation de l'adresse de la ligne en dessous ------------
bc26 
ld a,h
add a,8 
ld h,a ; <---- le fameux que tu as oublié !
ret nc 
ld bc,linewidth ; on passe en 96 colonnes
add hl,bc
res 3,h
ret
;-----------------------------------------------------------------


;--- variables memoires -----
pixel db 0
nblb dw 0
nbdeltas dw 0
;----------------------------
`
var deltaScreenCodeDelta = `
;--- dimensions du sprite ----
large equ $LARGE$
haut equ $HAUT$
loadingaddress equ #200
linewidth equ $LIGNELARGE$
nbdelta equ $NBDELTA$
nbcolors equ $NBCOLORS$
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,$SETMODE$
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

;---- attente vbl ----------
waitvbl
	ld b,#f5 ; attente vbl
vbl     
	in a,(c)
	rra
	jp nc,vbl
	ret
;---------------------------

;--- application palette firmware -------------
palettefirmware ; hl pointe sur les valeurs de la palette
ld e,nbcolors
ld a,0
ld hl,palette

paletteloop
ld b,(hl)
ld c,b
push af
push de
push hl
call #bc32 ; af, de, hl corrupted
pop hl
pop de
pop af
inc a
inc hl
dec e
jr nz,paletteloop
ret
;---------------------------------------------

;---------------------------------------------

;---- recuperation de l'adresse de la ligne en dessous ------------
bc26 
ld a,h
add a,8 
ld h,a ; <---- le fameux que tu as oublié !
ret nc 
ld bc,linewidth ; on passe en 96 colonnes
add hl,bc
res 3,h
ret
;-----------------------------------------------------------------


;--- variables memoires -----
pixel db 0 

;----------------------------

`

var deltaScreenCompressCodeDelta string = `
;--- dimensions du sprite ----
large equ $LARGE$
haut equ $HAUT$
loadingaddress equ #200
linewidth equ $LIGNELARGE$
nbdelta equ $NBDELTA$
nbcolors equ $NBCOLORS$
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,$SETMODE$
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



	;
	; Decompactage ZX0
	; HL = source
	; DE = destination
	;
	Depack:
		ld    bc,#ffff        ; preserve default offset 1
		push    bc
		inc    bc
		ld    a,#80
	dzx0s_literals:
		call    dzx0s_elias        ; obtain length
		ldir                ; copy literals
		add    a,a            ; copy from last offset or new offset?
		jr    c,dzx0s_new_offset
		call    dzx0s_elias        ; obtain length
	dzx0s_copy:
		ex    (sp),hl            ; preserve source,restore offset
		push    hl            ; preserve offset
		add    hl,de            ; calculate destination - offset
		ldir                ; copy from offset
		pop    hl            ; restore offset
		ex    (sp),hl            ; preserve offset,restore source
		add    a,a            ; copy from literals or new offset?
		jr    nc,dzx0s_literals
	dzx0s_new_offset:
		call    dzx0s_elias        ; obtain offset MSB
		ld b,a
		pop    af            ; discard last offset
		xor    a            ; adjust for negative offset
		sub    c
		RET    Z            ; Plus d'octets a traiter = fini
	
		ld    c,a
		ld    a,b
		ld    b,c
		ld    c,(hl)            ; obtain offset LSB
		inc    hl
		rr    b            ; last offset bit becomes first length bit
		rr    c
		push    bc            ; preserve new offset
		ld    bc,1            ; obtain length
		call    nc,dzx0s_elias_backtrack
		inc    bc
		jr    dzx0s_copy
	dzx0s_elias:
		inc    c            ; interlaced Elias gamma coding
	dzx0s_elias_loop:
		add    a,a
		jr    nz,dzx0s_elias_skip
		ld    a,(hl)            ; load another group of 8 bits
		inc    hl
		rla
	dzx0s_elias_skip:
		ret     c
	dzx0s_elias_backtrack:
		add    a,a
		rl    c
		rl    b
		jr    dzx0s_elias_loop
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

;---- attente vbl ----------
waitvbl
	ld b,#f5 ; attente vbl
vbl     
	in a,(c)
	rra
	jp nc,vbl
	ret
;---------------------------

;--- application palette firmware -------------
palettefirmware ; hl pointe sur les valeurs de la palette
ld e,nbcolors
ld a,0
ld hl,palette

paletteloop
ld b,(hl)
ld c,b
push af
push de
push hl
call #bc32 ; af, de, hl corrupted
pop hl
pop de
pop af
inc a
inc hl
dec e
jr nz,paletteloop
ret
;---------------------------------------------

;---------------------------------------------

;---- recuperation de l'adresse de la ligne en dessous ------------
bc26 
ld a,h
add a,8 
ld h,a ; <---- le fameux que tu as oublié !
ret nc 
ld bc,linewidth ; on passe en 96 colonnes
add hl,bc
res 3,h
ret
;-----------------------------------------------------------------


;--- variables memoires -----
pixel db 0 

;----------------------------

`

var deltaCodeDelta string = `;--- dimensions du sprite ----
large equ $LARGE$
haut equ $HAUT$
loadingaddress equ #200
linewidth equ $LIGNELARGE$
nbdelta equ $NBDELTA$
nbcolors equ $NBCOLORS$
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,$SETMODE$
	call #BC0E
;-------------------------------

;--- gestion de la palette ----
	call palettefirmware
;------------------------------

call xvbl

;--- affichage du sprite initiale --
	; affichage du premier sprite
	ld de,$INITIALADDRESS$ ; adresse de l'ecran 
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

;---- attente vbl ----------
waitvbl
	ld b,#f5 ; attente vbl
vbl     
	in a,(c)
	rra
	jp nc,vbl
	ret
;---------------------------

;--- application palette firmware -------------
palettefirmware ; hl pointe sur les valeurs de la palette
ld e,nbcolors
ld a,0
ld hl,palette

paletteloop
ld b,(hl)
ld c,b
push af
push de
push hl
call #bc32 ; af, de, hl corrupted
pop hl
pop de
pop af
inc a
inc hl
dec e
jr nz,paletteloop
ret
;---------------------------------------------

;---------------------------------------------

;---- recuperation de l'adresse de la ligne en dessous ------------
bc26 
ld a,h
add a,8 
ld h,a ; <---- le fameux que tu as oublié !
ret nc 
ld bc,linewidth ; on passe en 96 colonnes
add hl,bc
res 3,h
ret
;-----------------------------------------------------------------


;--- variables memoires -----
pixel db 0 
;----------------------------`

var depackRoutine = `
;--- dimensions du sprite ----
large equ $LARGE$
haut equ $HAUT$
loadingaddress equ #200
linewidth equ $LIGNELARGE$
nbdelta equ $NBDELTA$
nbcolors equ $NBCOLORS$
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,$SETMODE$
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

	ld de,$INITIALADDRESS$ ; adresse de l'ecran 
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



	;
	; Decompactage ZX0
	; HL = source
	; DE = destination
	;
	Depack:
		ld    bc,#ffff        ; preserve default offset 1
		push    bc
		inc    bc
		ld    a,#80
	dzx0s_literals:
		call    dzx0s_elias        ; obtain length
		ldir                ; copy literals
		add    a,a            ; copy from last offset or new offset?
		jr    c,dzx0s_new_offset
		call    dzx0s_elias        ; obtain length
	dzx0s_copy:
		ex    (sp),hl            ; preserve source,restore offset
		push    hl            ; preserve offset
		add    hl,de            ; calculate destination - offset
		ldir                ; copy from offset
		pop    hl            ; restore offset
		ex    (sp),hl            ; preserve offset,restore source
		add    a,a            ; copy from literals or new offset?
		jr    nc,dzx0s_literals
	dzx0s_new_offset:
		call    dzx0s_elias        ; obtain offset MSB
		ld b,a
		pop    af            ; discard last offset
		xor    a            ; adjust for negative offset
		sub    c
		RET    Z            ; Plus d'octets a traiter = fini
	
		ld    c,a
		ld    a,b
		ld    b,c
		ld    c,(hl)            ; obtain offset LSB
		inc    hl
		rr    b            ; last offset bit becomes first length bit
		rr    c
		push    bc            ; preserve new offset
		ld    bc,1            ; obtain length
		call    nc,dzx0s_elias_backtrack
		inc    bc
		jr    dzx0s_copy
	dzx0s_elias:
		inc    c            ; interlaced Elias gamma coding
	dzx0s_elias_loop:
		add    a,a
		jr    nz,dzx0s_elias_skip
		ld    a,(hl)            ; load another group of 8 bits
		inc    hl
		rla
	dzx0s_elias_skip:
		ret     c
	dzx0s_elias_backtrack:
		add    a,a
		rl    c
		rl    b
		jr    dzx0s_elias_loop
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

;---- attente vbl ----------
waitvbl
	ld b,#f5 ; attente vbl
vbl     
	in a,(c)
	rra
	jp nc,vbl
	ret
;---------------------------

;--- application palette firmware -------------
palettefirmware ; hl pointe sur les valeurs de la palette
ld e,nbcolors
ld a,0
ld hl,palette

paletteloop
ld b,(hl)
ld c,b
push af
push de
push hl
call #bc32 ; af, de, hl corrupted
pop hl
pop de
pop af
inc a
inc hl
dec e
jr nz,paletteloop
ret
;---------------------------------------------

;---------------------------------------------

;---- recuperation de l'adresse de la ligne en dessous ------------
bc26 
ld a,h
add a,8 
ld h,a ; <---- le fameux que tu as oublié !
ret nc 
ld bc,linewidth ; on passe en 96 colonnes
add hl,bc
res 3,h
ret
;-----------------------------------------------------------------


;--- variables memoires -----
pixel db 0 

;----------------------------
`
