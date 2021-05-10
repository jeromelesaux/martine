package gfx

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeromelesaux/martine/constants"
	"github.com/jeromelesaux/martine/export"
	"github.com/jeromelesaux/martine/export/file"
)

func DeltaPacking(gitFilepath string, ex *export.ExportType, initialAddress uint16, mode uint8) error {
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
	images := convertToImage(*gifImages)
	var pad int = 1
	if len(images) > maxImages {
		fmt.Fprintf(os.Stderr, "Warning gif exceed 30 images. Will corrupt the number of images.")
		pad = len(images) / maxImages
	}
	rawImages := make([][]byte, 0)
	deltaData := make([]*DeltaCollection, 0)
	var palette color.Palette
	var raw []byte

	// now transform images as win or scr
	fmt.Printf("Let's go transform images files in win or scr\n")

	if ex.FilloutGif {
		imgs := filloutGif(*gifImages, ex)
		_, palette, _, err = InternalApplyOneImage(imgs[0], ex, int(mode), palette, mode)
		if err != nil {
			return err
		}
		for i := 0; i < len(imgs); i += pad {
			in := imgs[i]
			/*	fw, _ := os.Create(ex.OutputPath + fmt.Sprintf("/a%.2d.png", i))
				png.Encode(fw, in)
				fw.Close()*/
			raw, _, _, err = InternalApplyOneImage(in, ex, int(mode), palette, mode)
			if err != nil {
				return err
			}
			rawImages = append(rawImages, raw)
			fmt.Printf("Image [%d] proceed\n", i)
		}
	} else {
		_, palette, _, err = InternalApplyOneImage(images[0], ex, int(mode), palette, mode)
		if err != nil {
			return err
		}
		for i := 0; i < len(images); i += pad {
			in := images[i]
			raw, _, _, err = InternalApplyOneImage(in, ex, int(mode), palette, mode)
			if err != nil {
				return err
			}
			rawImages = append(rawImages, raw)
			fmt.Printf("Image [%d] proceed\n", i)
		}
	}
	lineOctetsWidth := ex.LineWidth
	x0, y0, err := CpcCoordinates(initialAddress, 0xC000, lineOctetsWidth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while computing cpc coordinates :%v\n", err)
	}

	fmt.Printf("Let's go deltapacking raw images\n")
	realSize := &constants.Size{Width: ex.Size.Width, Height: ex.Size.Height}
	realSize.Width = realSize.ModeWidth(mode)

	for i := 0; i < len(rawImages)-1; i++ {
		fmt.Printf("Compare image [%d] with [%d] ", i, i+1)
		d1 := rawImages[i]
		d2 := rawImages[i+1]
		if len(d1) != len(d2) {
			return ErrorSizeDiffers
		}
		dc := Delta(d1, d2, isSprite, *realSize, mode, uint16(x0), uint16(y0), lineOctetsWidth)
		deltaData = append(deltaData, dc)
		fmt.Printf("%d bytes differ from the both images\n", len(dc.Items))
	}
	fmt.Printf("Compare image [%d] with [%d] ", len(rawImages)-1, 0)
	d1 := rawImages[len(rawImages)-1]
	d2 := rawImages[0]
	dc := Delta(d1, d2, isSprite, ex.Size, mode, uint16(x0), uint16(y0), lineOctetsWidth)
	deltaData = append(deltaData, dc)
	fmt.Printf("%d bytes differ from the both images\n", len(dc.Items))

	return exportDeltaAnimate(rawImages[0], deltaData, palette, ex, initialAddress, mode, ex.OutputPath+string(filepath.Separator)+"delta.asm")
}

func convertToImage(g gif.GIF) []image.Image {
	c := make([]image.Image, 0)
	width := g.Image[0].Bounds().Max.X
	height := g.Image[0].Bounds().Max.Y
	for i := 1; i < len(g.Image)-1; i++ {
		img := image.NewNRGBA(image.Rect(0, 0, width, height))
		draw.Draw(img, img.Bounds(), g.Image[i], image.Point{0, 0}, draw.Over)
		c = append(c, img)
	}
	return c
}

func filloutGif(g gif.GIF, ex *export.ExportType) []image.Image {
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
		/*fw, _ := os.Create(ex.OutputPath + fmt.Sprintf("/%.2d.png", i))
		png.Encode(fw, reference)
		fw.Close()*/
		c = append(c, img)
	}
	return c
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
	dataCode += "table_delta:\n"
	file.ByteToken = "dw"
	dataCode += file.FormatAssemblyString(deltaIndex, "\n")

	file.ByteToken = "db"
	dataCode += "palette:\n" + file.ByteToken + " "
	dataCode += file.FormatAssemblyBasicPalette(palette, "\n")

	// replace the initial address
	address := fmt.Sprintf("#%.4x", initialAddress)
	header := strings.Replace(DeltaCodeDelta, "$INITIALADDRESS$", address, 1)

	// replace number of colors
	nbColors := fmt.Sprintf("%d", len(palette))
	header = strings.Replace(header, "$NBCOLORS$", nbColors, 1)

	// replace the number of delta
	nbDelta := fmt.Sprintf("%d", len(delta)+1)
	header = strings.Replace(header, "$NBDELTA$", nbDelta, 1)

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

	fw, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fw.Close()
	fw.WriteString(code)
	return nil
}

var DeltaCodeDelta string = ";--- dimensions du sprite ----\n" +
	"large equ $LARGE$\n" +
	"haut equ $HAUT$\n" +
	"loadingaddress equ #200\n" +
	"linewidth equ $LIGNELARGE$\n" +
	"nbdelta equ $NBDELTA$\n" +
	"nbcolors equ $NBCOLORS$\n" +
	";-----------------------------\n" +
	"org loadingaddress\n" +
	"run loadingaddress\n" +
	";-----------------------------\n" +
	"start\n" +
	";--- selection du mode ---------\n" +
	"    ld a,$SETMODE$\n" +
	"    call #BC0E\n" +
	";-------------------------------\n" +
	"\n" +
	";--- gestion de la palette ---- \n" +
	"    call palettefirmware\n" +
	";------------------------------\n" +
	"\n" +
	"call xvbl\n" +
	"\n" +
	";--- affichage du sprite initiale --  \n" +
	"    ; affichage du premier sprite\n" +
	"    ld de,$INITIALADDRESS$ ; adresse de l'ecran \n" +
	"    ld hl,sprite ; pointeur sur l'image en memoire \n" +
	"    ld b, haut ; hauteur de l'image \n" +
	"    loop \n" +
	"    push bc ; sauve le compteur hauteur dans la pile \n" +
	"    push de ; sauvegarde de l'adresse ecran dans la pile\n" +
	"    ld bc, large ; largeur de l'image a afficher\n" +
	"    ldir ; remplissage de n * largeur octets a l'adresse dans de \n" +
	"    pop de ; recuperation de l'adresse d'origine \n" +
	"    ex de,hl ; echange des valeurs des adresses\n" +
	"    call bc26 ; calcul de l'adresse de la ligne suivante\n" +
	"    ex de,hl ; echange des valeurs des adresses\n" +
	"    pop bc ; retabli le compteur \n" +
	"    djnz loop\n" +
	";------------------------------------\n" +
	"\n" +
	"mainloop    ; routine pour afficher les deltas provenant de martine \n" +
	"\n" +
	";call #bb06\n" +
	"\n" +
	"call xvbl\n" +
	"call next_delta\n" +
	"\n" +
	"jp mainloop\n" +
	"\n" +
	"\n" +
	";--- routine de deltapacking --------------------------\n" +
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
	" ld a,(hl) ; nombre de byte a poker\n" +
	" push af   ; stockage en mémoire\n" +
	" inc hl\n" +
	"init\n" +
	" ld a,(hl) ; octet a poker\n" +
	" ld (pixel),a\n" +
	" inc hl\n" +
	" ld c,(hl) ; nbfois\n" +
	" inc hl \n" +
	" ld b,(hl)\n" +
	" inc hl\n" +
	";\n" +
	"poke_octet\n" +
	" ld e,(hl)\n" +
	" inc hl\n" +
	" ld d,(hl) ; de=adresse\n" +
	" inc hl\n" +
	" ld a,(pixel)\n" +
	" ld (de),a ; poke a l'adresse dans de\n" +
	" dec bc\n" +
	" ld a,b ; test a t'on poke toutes les adresses compteur bc\n" +
	" or a \n" +
	" jr nz, poke_octet\n" +
	" ld a,c \n" +
	" or a\n" +
	" jr nz, poke_octet\n" +
	" pop af \n" +
	"; reste t'il d'autres bytes a poker ? \n" +
	" dec a \n" +
	" push af\n" +
	" jr nz,init\n" +
	" pop af\n" +
	" ret\n" +
	"\n" +
	"\n" +
	"\n" +
	";---------------------------------------------------------------\n" +
	";\n" +
	"; attente de plusieurs vbl\n" +
	";\n" +
	"xvbl ld e,50\n" +
	"	call waitvbl\n" +
	"	dec e\n" +
	"	jr nz,xvbl+2\n" +
	"	ret\n" +
	";-----------------------------------\n" +
	"\n" +
	";---- attente vbl ----------\n" +
	"waitvbl\n" +
	"    ld b,#f5 ; attente vbl\n" +
	"vbl     \n" +
	"    in a,(c)\n" +
	"    rra\n" +
	"    jp nc,vbl\n" +
	"    ret\n" +
	";---------------------------\n" +
	"\n" +
	";--- application palette firmware -------------\n" +
	"palettefirmware ; hl pointe sur les valeurs de la palette\n" +
	"ld e,nbcolors\n" +
	"ld a,0\n" +
	"ld hl,palette\n" +
	"\n" +
	"paletteloop\n" +
	"ld b,(hl)\n" +
	"ld c,b\n" +
	"push af\n" +
	"push de\n" +
	"push hl\n" +
	"call #bc32 ; af, de, hl corrupted\n" +
	"pop hl\n" +
	"pop de\n" +
	"pop af\n" +
	"inc a\n" +
	"inc hl\n" +
	"dec e\n" +
	"jr nz,paletteloop\n" +
	"ret\n" +
	";---------------------------------------------\n" +
	"\n" +
	";---------------------------------------------\n" +
	"\n" +
	";---- recuperation de l'adresse de la ligne en dessous ------------\n" +
	"bc26 \n" +
	"ld a,h\n" +
	"add a,8 \n" +
	"ld h,a ; <---- le fameux que tu as oublié !\n" +
	"ret nc \n" +
	"ld bc,linewidth ; on passe en 96 colonnes\n" +
	"add hl,bc\n" +
	"res 3,h\n" +
	"ret\n" +
	";-----------------------------------------------------------------\n" +
	"\n" +
	"\n" +
	";--- variables memoires -----\n" +
	"pixel db 0 \n" +
	";----------------------------\n"
