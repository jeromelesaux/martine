	org #8000
	nolist
	run start

mode0	equ #7f8c
mode1	equ #7f8d
mode2	equ #7f8e
depart	equ 71		; 71 part du border haut... ligne impaire...
lignes	equ 200/2	; 200 lignes dont 100 en mode x et 100 en mode y

start
	di
	ld hl,#c9fb	; pas forcement utile puisque tout est sous di...
	ld (#38),hl

delock_asic
	ld bc,#bc00
	ld hl,sequence
	ld e,17
	seq:
	ld a,(hl)
	out (c),a
	inc hl
	dec e
	jr nz,seq

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
; Connecting I/O Page
	ld bc,#7FB8
	out (c),c

	ld hl,egx_palette
	ld de,#6400
	ld bc,32
	ldir

	; Deconnecting I/O Page
	ld bc,#7fA0
	out (c),c

	jp main

sequence:
db #ff,#00,#ff,#77,#b3,#51,#a8,#d4,#62,#39,#9c,#46,#2b,#15,#8a,#cd,#ee


egx_palette			
DB #00, #00, #03, #33, #0C, #C9, #03, #00, #0C, #CC, #06, #63, #09, #96, #03, #33, #09, #C9, #0C, #C9, #03, #33, #00, #30, #00, #30, #06, #66, #03, #33, #0F, #FC
; attention, les 15+1 encres sont stockees de la derniere a la premiere soit de ink 15 a ink 0.	
end 
;save 'egx.bin',#8000,end-start,DSK,'D1.dsk'
save 'egx.bin',#8000,end-start
