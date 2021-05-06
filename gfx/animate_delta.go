package gfx

import (
	"fmt"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
)

func DeltaPacking(gitFilepath string, ex *export.ExportType, initialAddress uint16, mode uint8) error {
	var isSprite = true
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
	var pad int = 1
	if len(gifImages.Image) > 30 {
		fmt.Fprintf(os.Stderr, "Warning gif exceed 30 images. Will corrupt the number of images.")
		pad = len(gifImages.Image) / 30
	}
	rawImages := make([][]byte, 0)
	deltaData := make([]*DeltaCollection, 0)
	var palette color.Palette
	var raw []byte
	// now transform images as win or scr
	fmt.Printf("Let's go transform images files in win or scr\n")
	for i := 0; i < len(gifImages.Image); i += pad {
		in := gifImages.Image[i]
		raw, palette, _, err = InternalApplyOneImage(in, ex, int(mode), mode)
		if err != nil {
			return err
		}
		rawImages = append(rawImages, raw)
		fmt.Printf("Image [%d] proceed\n", i)
	}
	lineOctetsWidth := ex.LineWidth
	x0, y0, err := CpcCoordinates(initialAddress, 0xC000, lineOctetsWidth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while computing cpc coordinates :%v\n", err)
	}

	fmt.Printf("Let's go deltapacking raw images\n")
	for i := 0; i < len(rawImages)-1; i++ {
		fmt.Printf("Compare image [%d] with [%d] ", i, i+1)
		d1 := rawImages[i]
		d2 := rawImages[i+1]
		if len(d1) != len(d2) {
			return ErrorSizeDiffers
		}
		dc := Delta(d1, d2, isSprite, ex.Size, mode, uint16(x0), uint16(y0), lineOctetsWidth)
		deltaData = append(deltaData, dc)
		fmt.Printf("%d bytes differ from the both images\n", len(dc.Items))
	}
	fmt.Printf("Compare image [%d] with [%d] ", 0, len(rawImages)-1)
	d1 := rawImages[len(rawImages)-1]
	d2 := rawImages[0]
	dc := Delta(d1, d2, isSprite, ex.Size, mode, uint16(x0), uint16(y0), lineOctetsWidth)
	deltaData = append(deltaData, dc)
	fmt.Printf("%d bytes differ from the both images\n", len(dc.Items))

	return exportDeltaAnimate(rawImages[0], deltaData, palette, ex, initialAddress, mode, ex.OutputPath+string(filepath.Separator)+"delta.asm")
}

func exportDeltaAnimate(imageReference []byte, delta []*DeltaCollection, palette color.Palette, ex *export.ExportType, initialAddress uint16, mode uint8, filename string) error {

	var dataCode string
	var deltaIndex []string
	var code string
	// copy of the sprite
	dataCode += "sprite:\n"
	dataCode += file.FormatAssemblyDatabyte(imageReference, "\n")
	// copy of all delta
	for i := 0; i < len(delta); i++ {
		dc := delta[i]
		data, err := dc.Marshall()
		if err != nil {
			return err
		}
		name := fmt.Sprintf("delta%.2d", i)
		dataCode += name + ":\n"
		dataCode += file.FormatAssemblyDatabyte(data, "\n")
		deltaIndex = append(deltaIndex, name)
	}
	dataCode += DeltaCodeDeltaTable
	file.ByteToken = "dw"
	dataCode += file.FormatAssemblyString(deltaIndex, "\n")
	file.ByteToken = "db"
	dataCode += "hardpalette:\n" + file.ByteToken + " "
	dataCode += file.FormatAssemblyCPCPalette(palette, "\n")

	// replace the color number in palette
	nbColor := fmt.Sprintf("%d", len(palette))
	header := strings.Replace(DeltaCodeHeader, "$NBCOLORS$", nbColor, 1)

	// replace the initial address
	address := fmt.Sprintf("#%.4x", initialAddress)
	header = strings.Replace(header, "$INITIALADDRESS$", address, 1)

	// replace the number of delta
	nbDelta := fmt.Sprintf("%d", len(delta))
	header = strings.Replace(header, "$NBDELTA$", nbDelta, 1)

	// replace char large for the screen
	charLarge := fmt.Sprintf("%d", ex.LineWidth)
	header = strings.Replace(header, "$LIGNELARGE$", charLarge, 1)

	// replace heigth
	height := fmt.Sprintf("%d", ex.Size.Height)
	header = strings.Replace(header, "$HAUT$", height, 1)

	// replace width
	var width string
	switch mode {
	case 0:
		width = fmt.Sprintf("%d", ex.Size.Width/2)
	case 1:
		width = fmt.Sprintf("%d", ex.Size.Width/4)
	case 2:
		width = fmt.Sprintf("%d", ex.Size.Width/8)
	}
	header = strings.Replace(header, "$LARGE$", width, 1)

	var modeSet string
	switch mode {
	case 0:
		modeSet = "#7f8c"
	case 1:
		modeSet = "#7f8d"
	case 2:
		modeSet = "#7f8e"
	}

	// replace mode
	header = strings.Replace(header, "$SETMODE$", modeSet, 1)

	code += header
	code += DeltaCodeNextDelta
	code += DeltaCodeDrawSprite
	code += DeltaCodePalette
	code += DeltaCodeBC26
	code += DeltaCodeVbl
	code += dataCode
	code += "\nend"

	fw, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fw.Close()
	fw.WriteString(code)
	return nil
}

var DeltaCodeDeltaTable string = "table_delta:\n"

var DeltaCodeNextDelta string = ";--- routine de deltapacking --------------------------\n" +
	"next_delta:\n" +
	"table_index:\n" +
	"    ld a,-1\n" +
	"    inc a\n" +
	"    and nbdelta\n" +
	"    ld (table_index+1),a\n" +
	"    add a\n" +
	"    ld e,a\n" +
	"    ld d,0\n" +
	"    ld hl,table_delta\n" +
	"    add hl,de\n" +
	"    ld a,(hl)\n" +
	"    inc hl\n" +
	"    ld h,(hl)\n" +
	"    ld l,a\n" +
	"delta\n" +
	"    ld a,(hl) ; nombre de byte a poker\n" +
	"    inc hl\n" +
	"init\n" +
	"    ex af,af'\n" +
	"    ld a,(hl) ; octet a poker\n" +
	"    ld (pixel1+1),a\n" +
	"    inc hl\n" +
	"    ld c,(hl) ; nbfois\n" +
	"    inc hl\n" +
	"    ld b,(hl)\n" +
	"    inc hl\n" +
	";\n" +
	"poke_octet\n" +
	"    ld e,(hl)\n" +
	"    inc hl\n" +
	"    ld d,(hl) ; de=adresse\n" +
	"    inc hl\n" +
	"pixel1:\n" +
	"     ld a,0\n" +
	"ld (de),a ; poke a l'adresse dans de\n" +
	";------------------\n" +
	"    dec bc\n" +
	"    ld a,b ; test a t'on poke toutes les adresses compteur bc\n" +
	"    or c  ; optimisation siko\n" +
	"    jr nz, poke_octet\n" +
	"    ex af,af'\n" +
	"    dec a ; reste t'il d'autres bytes a poker ?\n" +
	"    jr nz, init\n" +
	"    ret\n" +
	";---------------------------------------------------\n"

var DeltaCodeDrawSprite string = "drawSprite:\n" +
	".loop\n" +
	"push af ; sauve le compteur hauteur dans la pile\n" +
	"push de ; sauvegarde de l'adresse ecran dans la pile\n" +
	"push bc\n" +
	"ldir ; remplissage de n * largeur octets a l'adresse dans de\n" +
	"pop bc\n" +
	"pop de ; recuperation de l'adresse d'origine\n" +
	"ex de,hl ; echange des valeurs des adresses\n" +
	"call bc26 ; calcul de l'adresse de la ligne suivante\n" +
	"ex de,hl ; echange des valeurs des adresses\n" +
	"pop af ; retabli le compteur\n" +
	"dec a\n" +
	"jr nz, .loop\n" +
	"ret\n"

var DeltaCodeHeader string = ";--- dimensions du sprite ----\n" +
	"large equ $LARGE$\n" +
	"haut equ $HAUT$\n" +
	"loadingaddress equ #200\n" +
	"linewidth equ $LIGNELARGE$\n" +
	"nbdelta equ $NBDELTA$\n" +
	"nbcolors equ $NBCOLORS$\n" +
	";-----------------------------\n" +
	"org loadingaddress\n" +
	"run loadingaddress\n" +
	"start\n" +
	"di\n" +
	"ld bc,$SETMODE$ ; Mode 1\n" +
	"out (c),c\n" +
	"ld a,#c3\n" +
	"ld (#38),a\n" +
	"ld sp,loadingaddress\n" +
	"ld hl,hardpalette\n" +
	"call setpalette\n" +
	"call xvbl\n" +
	";--- affichage du sprite initial -\n" +
	"ld de,$INITIALADDRESS$ ; adresse de l'ecran\n" +
	"ld hl,sprite ; pointeur sur l'image en memoire\n" +
	"ld bc, large ; hauteur de l'image\n" +
	"ld a,haut\n" +
	"call drawSprite\n" +
	"call xvbl\n" +
	"ei\n" +
	";------------------------------------\n" +
	"mainloop\n" +
	"ld e,3\n" +
	"call xvbl.lp\n" +
	"call next_delta\n" +
	"jp mainloop\n"

var DeltaCodePalette string = ";--- application palette hardware ------------\n" +
	"setpalette\n" +
	"	ld b,#7F          ; gatearray pointer to ink 0\n" +
	"	xor a                ; ink number start with 0\n" +
	".loop\n" +
	"	ld e,(hl)\n" +
	"	out (c),a            ; on selectionne la couleur\n" +
	"	out (c),e            ; on envoie la couleur\n" +
	"	inc hl\n" +
	"	inc a\n" +
	"	cp nbcolors\n" +
	"	jr nz,.loop\n" +
	"	ld c,#10          ; meme chose pour le border avec la couleur 0\n" +
	"	out (c),c\n" +
	"	ld e,(hl)\n" +
	"	out (c),e\n" +
	"	ret\n" +
	";----------------------------------------------\n"

var DeltaCodeVbl string = ";---------------------------------------------------------------\n" +
	";\n" +
	"; attente de plusieurs vbl\n" +
	";\n" +
	"xvbl ld e,4\n" +
	".lp:\n" +
	"	call waitvbl\n" +
	"	dec e\n" +
	"	jr nz,.lp\n" +
	"	ret\n\n" +
	";-----------------------------------\n" +
	";---- attente vbl ----------\n" +
	"waitvbl\n" +
	"	ld b,#f5 ; attente vbl\n" +
	"vbl\n" +
	"	in a,(c)\n" +
	"	rra\n" +
	"	jr nc,vbl\n" +
	"	;ld b,#f5 ; attente vbl\n" +
	"ret\n"

var DeltaCodeBC26 string = ";---- recuperation de l'adresse de la ligne en dessous ------------\n" +
	"; CE: hl\n" +
	"; CS: hl contient l'adresse de la ligne suivante\n" +
	";     A,F modifiés\n" +
	"bc26\n" +
	"	ld   a,h\n" +
	"	add  8\n" +
	"	ld   h,a\n" +
	"	and  #38\n" +
	"	ret  nz\n" +
	"	ld   a,h\n" +
	"	sub  #40\n" +
	"	ld   h,a\n" +
	"	ld   a,l\n" +
	"	add  linewidth ; #60      ; 96 chars\n" +
	"	ld   l,a\n" +
	"	ret  nc\n" +
	"	inc  h\n" +
	"	ld   a,h\n" +
	"	and  7\n" +
	"	ret  nz\n" +
	"	ld   a,h\n" +
	"	sub  8\n" +
	"	ld   h,a\n" +
	"	ret\n" +
	";-----------------------------------------------------------------\n"
