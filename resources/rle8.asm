		org &1000
		nolist
		run start


rlesize		equ finrle-debutrle	; taille du fichier rle
rlesizediv2	equ rlesize/2		; taille du fichier rle div 2

start
		ld hl,&4000		; rle_source
		ld de,&c000		; destination

		ld bc,rlesizediv2
pool
		push bc	
		ld b,(hl)
		inc hl
		ld a,(hl)
loop
		ld (de),a
		inc de
		djnz loop
		inc hl
		pop bc
		dec bc
		ld a,b
		or c
		jr nz,pool

fin		ret

		org &4000
		nolist
debutrle
incbin		"fichier.scr"
finrle