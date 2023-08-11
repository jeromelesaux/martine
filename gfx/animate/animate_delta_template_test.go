package animate_test

import (
	"bytes"
	"fmt"
	"image/color"
	"testing"
	"text/template"

	"github.com/jeromelesaux/martine/export/ascii"
	"github.com/jeromelesaux/martine/log"
	zx0 "github.com/jeromelesaux/zx0/encode"
)

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

	code += "\nsprite:\n"
	code += ascii.FormatAssemblyDatabyte(a.Image, "\n")

	ascii.ByteToken = "db"
	if a.Type.Compress {
		for i, v := range a.Delta {
			log.GetLogger().Info("Using Zx0 cruncher")
			d := zx0.Encode(v)
			code += fmt.Sprintf("delta%.2d:", i)
			code += ascii.FormatAssemblyDatabyte(d, "\n")
		}
	} else {
		for i, v := range a.Delta {
			code += fmt.Sprintf("delta%.2d:", i)
			code += ascii.FormatAssemblyDatabyte(v, "\n")
		}
	}

	return code
}

func (a AnimateValues) TableDelta() string {
	code := "table_delta:\n"
	deltaIndexes := make([]string, 0)
	for i := range a.Delta {
		deltaIndexes = append(deltaIndexes, fmt.Sprintf("delta%.2d", i))
	}
	ascii.ByteToken = "dw"
	code += ascii.FormatAssemblyString(deltaIndexes, "\n")
	return code
}

func (a AnimateValues) DisplayPalette() string {
	code := "palette:\n"
	if a.Type.CPCPlus {
		code += ascii.FormatAssemblyCPCPlusPalette(a.Palette, "\n")
	} else {
		code += ascii.FormatAssemblyCPCPalette(a.Palette, "\n")
	}
	return code
}

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
{{ .DisplayCode }}
{{ .TableDelta }}
{{ .DisplayPalette }}
buffer: 
`

func TestTemplate(t *testing.T) {
	var buf bytes.Buffer
	a := AnimateExportType{Compress: false, CPCPlus: false, IsSprite: false}
	vals := AnimateValues{
		Type:           a,
		InitialAddress: "#C000",
		Large:          "#800",
		Haut:           "#4000",
		LigneLarge:     "#C080",
		Image:          []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Delta: [][]byte{
			{1, 2},
			{2, 5},
		},
		Palette: color.Palette{color.Black, color.White},
		Mode:    0,
	}

	temp := template.Must(template.New("code").Parse(string(AnimateTemplates[a])))
	err := temp.Execute(&buf, vals)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}
