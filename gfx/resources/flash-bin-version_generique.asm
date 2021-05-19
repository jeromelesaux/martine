; debut du code generique... On touche a rien en dessous...
; ---------------------------------------------------------

	org #3000
	nolist
	run start
start
	di
	ld hl,#c9fb
	ld (#38),hl
	ei
main
	ld b,#f5
sync	in a,(c)
	rra
	jr nc,sync

	ld hl,ecran1	
flip	ld a,0
	xor 1
	ld (flip+1),a
	jr nz,saut
	ld hl,ecran2
saut
	call inkmodecrtc
	halt
	halt
	halt

	jp main

inkmodecrtc
	ld bc,#7f0f	; 15+1 encres
cont
	out (c),c
	inc b
	outi
	dec c
	jp p,cont

	inc c	
	out (c),c
	ld a,(hl)
	out (c),a	; mode

	inc hl		
	ld bc,#bc0c	; crtc12
	out (c),c
	inc b
	ld a,(hl)
	out (c),a

	inc hl
	ld bc,#bc0d
	out (c),c
	inc b
	ld a,(hl)
	out (c),a	; crtc13
	
	ret

; fin du code generique... On touche a rien au dessus...
; ------------------------------------------------------

; C'est maintenant que tu "pokes" tes donnees (palette, mode, crtc)

; ici tu stock tes valeurs pour l'ecran1... c'est tout.
ecran1
ecr1_palette
db #53,#57,#4b,#54,#4b,#54,#4b,#54,#4b,#54,#4b,#54,#4a,#52,#4b,#54
; attention, les 15+1 encres sont stockees de la derniere a la premiere soit de ink 15 a ink 0.
; et pas grave si on envoi des encres inutiles (genre 15+1 encres pour le mode 1 ou 2...)
ecr1_mode
db #8e
; #8c=mode0 / #8d=mode1 / #8e=mode2 
ecr1_crtc
db #30,00
; crtc12=#30=#c000 / crtc13=#00=#c000+0 (utile si overscan plus tard...)
; on pourra ajouter d'autres valeurs crtc plus tard... (reg 2,1,8,6...) et modifier la routine du haut...

; ici tu stock tes valeurs pour l'ecran2... c'est tout.
ecran2
ecr2_palette
db #53,#55,#4b,#54,#4b,#54,#4b,#54,#4b,#54,#4b,#54,#4a,#52,#4b,#54
; attention, les 15+1 encres sont stockees de la derniere a la premiere soit de ink 15 a ink 0.
; et pas grave si on envoi des encres inutiles (genre 15+1 encres pour le mode 1 ou 2...)
ecr2_mode
db #8d
; #8c=mode0 / #8d=mode1 / #8e=mode2 
ecr2_crtc
db #10,00
; crtc12=#10=#4000 / crtc13=#00=#4000+0 (utile si overscan plus tard...)
; on pourra ajouter d'autres valeurs crtc plus tard... (reg 2,1,8,6...) et modifier la routine du haut...
end
;save 'flash.bin',#3000,end-start,DSK,'D1.dsk'
save 'flash.bin',#3000,end-start,AMSDOS
