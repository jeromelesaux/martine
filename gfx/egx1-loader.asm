
	org #3000
	nolist
	run start

mode0	equ #7f8c
mode1	equ #7f8d
mode2	equ #7f8e
depart	equ 71		; 71 part du border haut... ligne impaire...
lignes	equ 200/2	; 200 lignes dont 100 en mode x et 100 en mode y

start
	;ld hl,#c000
	;ld de,#c001
	;ld bc,#4000
	;ld a,r		; pour voir de visu si chgt ok...
	;ld (hl),a
	;ldir
    

	di
	ld hl,#c9fb	; pas forcement utile puisque tout est sous di...
	ld (#38),hl
main
	ld b,#f5
sync2	in a,(c)
	rra
	jr nc,sync2

	ds depart*64,0 	; 71*64=4544 nops, modes lances apres le border haut.
	ld e,lignes

loop	ld bc,mode1	; ligne impaire (1,3,5...195,197,199)
	out (c),c
	ds 57,0		; 64-7

	ld bc,mode0	; ligne paire (2,4,6...196,198,200)
	out (c),c
	ds 53,0		; 64-7-4

	dec e
	jp p,loop

	ld bc,#7f10	; border...
	ld a,#54	; noir...
	out (c),c	; si besoin...
	out (c),a

ink			; et la on balance les 16 encres...	
	ld hl,egx_palette
	ld bc,#7f0f	; 15+1 encres
cont
	out (c),c
	inc b
	outi
	dec c
	jp p,cont		

	jp main

egx_palette			
db #36,#4c,#58,#4e,#4b,#43,#5a,#59
db #4a,#46,#56,#5e,#47,#40,#5c,#54
; attention, les 15+1 encres sont stockees de la derniere a la premiere soit de ink 15 a ink 0.	
end 
save 'egx.bin',#3000,end-start,DSK,'D1.dsk'
;save 'egx.bin',#3000,end-start,AMSDOS
