
;--- dimensions du sprite ----
large equ 25
haut equ 100
loadingaddress equ #200
linewidth equ #c050
nbdelta equ 18
nbcolors equ 3
;-----------------------------
org loadingaddress
run loadingaddress
;-----------------------------
start
;--- selection du mode ---------
	ld a,1
	call #BC0E
;-------------------------------

;--- gestion de la palette ----
	call palettefirmware
;------------------------------

call xvbl

;--- affichage du sprite initiale --
	; affichage du premier sprite
	ld de,#c010 ; adresse de l'ecran 
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
	ld de,buffer
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
buffer dw 0
;----------------------------sprite:
db #92, #00, #20, #2d, #0f, #22, #e0, #00
db #46, #f0, #1e, #ce, #0a, #07, #02, #95
db #1e, #a0, #43, #29, #2c, #58, #96, #03
db #07, #95, #3c, #89, #01, #0f, #62, #2d
db #a4, #15, #a5, #0f, #6b, #00, #a0, #6b
db #0b, #5f, #ff, #ef, #d4, #2d, #ce, #a5
db #a0, #7c, #06, #e9, #cf, #48, #5a, #03
db #02, #90, #c0, #29, #40, #5a, #01, #02
db #90, #c0, #1a, #04, #56, #80, #00, #a1
db #c0, #0a, #83, #22, #78, #80, #05, #f8
db #61, #ce, #5a, #07, #68, #00, #a0, #69
db #5a, #0e, #8a, #40, #07, #24, #68, #08
db #a0, #40, #68, #03, #98, #70, #0e, #03
db #f5, #9c, #f9, #ce, #fc, #2d, #c3, #01
db #6b, #e2, #e0, #c0, #fb, #0e, #0c, #af
db #c4, #f0, #a0, #4f, #ce, #8a, #0d, #08
db #a8, #43, #5a, #06, #a2, #00, #4b, #24
db #aa, #07, #e0, #05, #8a, #0f, #42, #a2
db #87, #4a, #40, #04, #bc, #96, #6b, #7e
db #20, #ce, #02, #2a, #16, #2d, #10, #28
db #c0, #12, #08, #03, #0f, #1c, #e1, #a0
db #84, #59, #e5, #0b, #1e, #c3, #9c, #3e
db #ff, #cf, #81, #80, #28, #07, #68, #08
db #12, #86, #0e, #81, #00, #69, #07, #a0
db #1e, #48, #82, #01, #0c, #81, #0e, #26
db #21, #08, #03, #14, #aa, #03, #01, #14
db #88, #07, #00, #28, #1c, #13, #ba, #a8
db #ce, #84, #0c, #6c, #87, #85, #72, #ce
db #d7, #cf, #84, #07, #3e, #09, #f8, #d8
db #dc, #4f, #d7, #ac, #da, #21, #84, #40
db #ff, #73, #62, #21, #d6, #6e, #06, #ce
db #aa, #00, #28, #04, #a6, #a0, #60, #10
db #07, #d0, #4a, #40, #0b, #a0, #b0, #2e
db #84, #ac, #e8, #cf, #f0, #88, #10, #80
db #16, #9a, #00, #30, #40, #88, #70, #c0
db #46, #a8, #c0, #f0, #52, #24, #d0, #21
db #0a, #c3, #a1, #80, #18, #85, #a0, #1e
db #22, #a0, #a1, #28, #08, #42, #22, #68
db #03, #84, #00, #26, #1c, #87, #0f, #14
db #c3, #e1, #94, #38, #fa, #cf, #9d, #80
db #40, #4c, #0e, #f9, #ce, #16, #91, #87
db #69, #e1, #02, #84, #f0, #68, #07, #a1
db #d2, #0a, #20, #22, #1e, #e0, #40, #68
db #84, #46, #c9, #03, #a4, #48, #ce, #24
db #21, #13, #ce, #42, #60, #20, #02, #00
db #40, #88, #0c, #20, #53, #1a, #ed, #10
db #11, #a2, #00, #84, #08, #63, #10, #49
db #03, #34, #b0, #30, #61, #98, #77, #7b
db #c2, #d1, #18, #12, #a2, #0f, #0f, #1e
db #aa, #16, #e0, #00, #82, #03, #0e, #2c
db #0f, #2a, #1c, #95, #f0, #eb, #a2, #0e
db #88, #cf, #1e, #d2, #a5, #80, #2a, #03
db #62, #0d, #0f, #87, #8b, #70, #c0, #53
db #8a, #0a, #10, #5b, #b0, #a4, #cf, #f0
db #22, #21, #43, #48, #bc, #f0, #92, #68
db #2f, #01, #d0, #c2, #b4, #bb, #1e, #d2
db #bf, #68, #cd, #6e, #69, #0d, #78, #d0
db #cd, #9f, #5c, #94, #39, #a3, #e0, #96
db #0e, #b4, #ce, #6e, #0e, #a1, #61, #96
db #86, #e1, #e4, #04, #e5, #3f, #bc, #4b
db #1e, #78, #b4, #b0, #80, #70, #9c, #3e
db #6e, #5f, #0b, #4b, #f0, #b4, #69, #02
db #10, #ce, #29, #c0, #6a, #03, #68, #2d
db #69, #69, #25, #c3, #42, #40, #30, #a5
db #e0, #a9, #07, #65, #96, #3c, #d2, #c2
db #e1, #48, #30, #88, #c0, #20, #09, #a1
db #5a, #38, #f0, #ef, #a0, #f0, #40, #70
db #92, #a6, #a5, #f6, #c3, #6a, #ae, #2d
db #f0, #13, #bf, #50, #08, #b2, #cf, #90
db #d6, #4c, #b1, #23, #3c, #39, #e9, #97
db #fd, #c0, #30, #60, #d9, #f7, #86, #da
db #39, #01, #78, #e1, #ae, #70, #e1, #b4
db #c2, #67, #e3, #c0, #ce, #e4, #65, #ce
db #f8, #12, #b8, #d2, #00, #9d, #24, #60
db #00, #70, #d0, #c0, #00, #20
delta00:
db #00, #22, #40, #87, #09, #00, #17, #c0
db #18, #c8, #0e, #d9, #5b, #c9, #ad, #f1
db #4d, #da, #4d, #e2, #3a, #d3, #3b, #f3
db #78, #11, #00, #21, #c0, #71, #d0, #70
db #f8, #c0, #44, #ea, #d8, #c0, #e0, #c0
db #e8, #c0, #f0, #fd, #f9, #4d, #c2, #40
db #c3, #3e, #f3, #8c, #cb, #88, #eb, #85
db #fb, #d6, #db, #dd, #db, #08, #8d, #22
db #1a, #bc, #d0, #c1, #e8, #60, #42, #60
db #d1, #a9, #d9, #ee, #c2, #87, #f3, #8a
db #fb, #07, #16, #97, #17, #c8, #68, #c0
db #68, #c8, #b9, #92, #fd, #d0, #be, #d8
db #0a, #d1, #0e, #96, #0a, #d9, #5a, #c9
db #59, #f9, #5c, #82, #ac, #c1, #4d, #ea
db #a9, #f2, #a1, #fa, #9d, #c2, #aa, #ca
db #3e, #cb, #3b, #db, #e3, #4a, #eb, #2d
db #0a, #00, #20, #c8, #69, #d0, #0b, #c9
db #28, #e9, #3f, #cb, #39, #fb, #3c, #a0
db #86, #cb, #29, #d3, #84, #db, #80, #0e
db #00, #22, #c8, #a1, #d0, #c1, #f0, #a6
db #f8, #5a, #f1, #01, #da, #e2, #51, #ca
db #84, #68, #fa, #a1, #c2, #ed, #f2, #ea
db #fa, #8a, #d3, #e0, #db, #25, #01, #00
db #17, #d0, #21, #09, #00, #17, #d8, #4a
db #e0, #68, #d0, #5b, #e1, #59, #e9, #99
db #fa, #9d, #06, #3f, #c3, #83, #d3, #0f
db #25, #00, #20, #d8, #86, #e0, #6c, #c8
db #bb, #e0, #a9, #e8, #aa, #f0, #0b, #c1
db #d9, #69, #f1, #5f, #d1, #a6, #d9, #4f
db #e2, #f2, #9f, #e2, #9a, #ea, #9a, #f2
db #6a, #fa, #ea, #c2, #ca, #86, #d2, #eb
db #fa, #3b, #c3, #a9, #cb, #28, #d3, #3e
db #db, #3c, #e3, #3f, #a4, #37, #eb, #a6
db #f3, #89, #cb, #84, #d3, #89, #83, #db
db #85, #2a, #e3, #02, #eb, #d2, #c3, #f0
db #20, #00, #21, #d8, #aa, #e0, #e8, #a9
db #f0, #a4, #f8, #71, #c0, #88, #d8, #10
db #c9, #0d, #e9, #fe, #6a, #f1, #50, #fa
db #a0, #c2, #ca, #aa, #da, #e2, #92, #ea
db #8b, #d3, #89, #db, #8b, #9a, #88, #e3
db #89, #a6, #8b, #89, #eb, #93, #9a, #8f
db #fb, #90, #28, #d5, #cb, #aa, #d3, #d9
db #e2, #01, #da, #d5, #db, #4b, #05, #00
db #1b, #e0, #18, #e8, #bc, #c8, #bd, #f0
db #0c, #c9, #01, #14, #c9, #e8, #a9, #d8
db #aa, #e0, #be, #e8, #f0, #a0, #f8, #39
db #f1, #0a, #f9, #59, #e1, #a8, #d1, #ac
db #fd, #aa, #d9, #99, #e2, #ea, #32, #f2
db #ed, #3d, #4a, #ea, #39, #eb, #36, #f3
db #82, #32, #1e, #11, #71, #a2, #e8, #f8
db #6b, #9a, #70, #d0, #bd, #72, #0c, #e9
db #0f, #79, #a6, #c9, #e1, #5f, #e9, #82
db #ea, #4f, #fa, #9a, #b2, #9f, #7d, #5d
db #da, #83, #fb, #d2, #d3, #00, #23, #9a
db #29, #f0, #17, #f8, #0a, #f0, #68, #f8
db #bd, #6b, #09, #c1, #11, #0d, #70, #ed
db #83, #73, #f1, #a8, #c9, #af, #af, #59
db #e1, #f1, #65, #1d, #13, #f9, #fc, #c1
db #00, #c2, #51, #d2, #9c, #e2, #ed, #da
db #ed, #e2, #35, #fb, #44, #fb, #8a, #e3
db #8a, #eb, #8a, #f3, #87, #fb, #88, #fb
db #d7, #c3, #d8, #c3, #d9, #c3, #d1, #d3
db #da, #d3, #e3, #db, #0d, #01, #00, #6b
db #c0, #d0, #07, #00, #71, #c8, #3a, #db
db #87, #e3, #d8, #cb, #d7, #d3, #d8, #80
db #fd, #aa, #db, #70, #08, #00, #71, #e0
db #fd, #e1, #e9, #0a, #f1, #8b, #eb, #8f
db #54, #a3, #e0, #c3, #d7, #db, #10, #0c
db #00, #68, #e8, #be, #f0, #9c, #da, #ec
db #c2, #f0, #fa, #38, #e3, #93, #f3, #91
db #fb, #db, #c3, #df, #cb, #da, #db, #e1
db #60, #02, #83, #a5, #e8, #ea, #ea, #e0
db #85, #bd, #aa, #71, #f0, #c1, #c8, #10
db #e1, #10, #e9, #ea, #f2, #42, #43, #a5
db #43, #a9, #93, #e3, #93, #fb, #d9, #cb
db #e2, #a0, #86, #f8, #d6, #cb, #a4, #04
db #81, #e4, #c1, #c0, #a0, #f2, #8a, #c3
db #d6, #d3, #81, #dd, #a0, #bd, #c8, #85
db #c3, #3c, #09, #4a, #c0, #c8, #5f, #f1
db #5f, #f9, #ae, #e9, #ec, #e2, #3f, #d3
db #86, #c3, #86, #db, #89, #f3, #e5, #10
db #b3, #4a, #d0, #10, #d9, #10, #f9, #60
db #c1, #ae, #f9, #fe, #c1, #01, #c2, #fe
db #c9, #01, #ca, #fe, #d1, #fe, #e1, #01
db #ea, #01, #f2, #01, #fa, #51, #c2, #df
db #d9, #43, #47, #0f, #b9, #d8, #a8, #e9
db #ff, #e9, #4e, #c2, #82, #eb, #03, #13
db #e7, #61, #4f, #a5, #e0, #0e, #c1, #0e
db #c9, #a4, #e1, #5a, #c1, #5b, #d9, #59
db #f1, #6a, #c9, #a8, #e1, #e9, #c2, #e9
db #ca, #e9, #d2, #e9, #da, #e9, #e2, #f3
db #a8, #fb, #78, #e3, #88, #f3, #04, #01
db #17, #c2, #e0, #40, #d1, #a1, #c1, #2a
db #a8, #f1, #b1, #f9, #e3, #c3, #db, #db
db #14, #01, #bc, #ba, #24, #f6, #62, #f0
db #84, #05, #13, #f8, #0c, #c1, #9a, #d2
db #87, #eb, #da, #c3, #69, #b7, #dd, #c0
db #47, #21, #bb, #e6, #10, #c1, #3a, #fb
db #28, #c5, #0c, #d1, #e1, #f7, #e9, #39
db #60, #e7, #a3, #cb, #85, #f3, #d5, #48
db #03, #83, #dd, #d9, #38, #eb, #86, #4e
db #9d, #bd, #48, #0b, #e1, #ea, #e2, #0c
db #04, #e4, #11, #e0, #e1, #af, #c1, #a9
db #d1, #0b, #8f, #aa, #0b, #f9, #06, #03
db #a9, #1e, #e9, #fa, #37, #e3, #0e, #59
db #0a, #a9, #c9, #9a, #da, #4e, #e2, #ec
db #ea, #3d, #cb, #2c, #b5, #a7, #ae, #4d
db #ca, #89, #23, #3b, #eb, #18, #81, #9f
db #8a, #ae, #f1, #3f, #db, #3f, #eb, #3f
db #f3, #83, #ad, #ba, #41, #71, #ff, #82
db #3a, #c3, #30, #08, #85, #da, #9c, #fa
db #ec, #ca, #ec, #d2, #ec, #da, #89, #fb
db #8b, #43, #cb, #8e, #db, #20, #6b, #5e
db #ed, #ca, #ed, #d2, #da, #cb, #02, #ef
db #1a, #e9, #ea, #39, #c3, #8a, #22, #34
db #01, #39, #f2, #58, #f7, #89, #ed, #fa
db #52, #8a, #39, #cb, #94, #3a, #b8, #16
db #81, #7a, #39, #d3, #39, #e3, #96, #de
db #62, #db, #0a, #03, #4e, #e3, #3a, #eb
db #3a, #f3, #5a, #c3, #62, #3e, #eb, #c3
db #78, #3d, #fb, #b4, #89, #6a, #8d, #c3
db #87, #db, #d4, #49, #27, #86, #e3, #50
db #09, #e9, #ae, #86, #f3, #d7, #92, #cd
db #ae, #84, #b0, #f7, #62, #d6, #c3, #12
db #30, #db, #d3, #00, #08
delta01:
db #10, #0e, #3e, #07, #12, #00, #17, #c0
db #17, #d0, #b9, #c0, #0e, #d9, #5b, #c9
db #5a, #d1, #a9, #c1, #ad, #e9, #4d, #d2
db #4d, #da, #49, #e2, #4d, #e2, #e9, #da
db #eb, #f2, #36, #f3, #36, #fb, #3b, #fb
db #86, #c3, #87, #06, #00, #18, #c0, #bd
db #f0, #0b, #e1, #49, #ea, #8b, #c3, #86
db #cb, #2d, #07, #00, #1b, #c0, #1b, #d0
db #1b, #f8, #0b, #f1, #e9, #fa, #8c, #cb
db #84, #d3, #3c, #05, #00, #20, #c0, #20
db #d0, #70, #fd, #02, #f8, #c0, #c0, #80
db #16, #00, #22, #c0, #aa, #d8, #e0, #aa
db #e8, #f0, #81, #a5, #f8, #72, #c0, #bc
db #d8, #11, #c1, #60, #d1, #68, #e9, #ae
db #f9, #fe, #d9, #ea, #f2, #42, #d3, #87
db #cb, #86, #e3, #94, #eb, #42, #f3, #e3
db #c3, #dc, #d3, #db, #db, #0f, #24, #00
db #18, #c8, #1b, #e0, #18, #e8, #20, #f8
db #6b, #c0, #85, #a6, #c8, #0b, #e9, #0c
db #f1, #0b, #f9, #5f, #e1, #49, #f2, #fa
db #99, #c2, #aa, #ca, #d2, #aa, #da, #e2
db #a8, #ea, #a1, #f2, #9f, #a6, #39, #c3
db #37, #cb, #38, #37, #d3, #38, #9a, #37
db #db, #38, #a6, #3b, #37, #e3, #38, #ba
db #3b, #f9, #eb, #e0, #fd, #5a, #f3, #3c
db #fb, #83, #d3, #0d, #03, #00, #1b, #c8
db #0b, #c9, #39, #d3, #78, #12, #00, #21
db #c8, #aa, #d0, #d8, #aa, #e0, #e8, #68
db #f8, #71, #c0, #4a, #c8, #70, #f0, #c0
db #f8, #10, #c1, #ae, #f1, #fd, #e5, #50
db #55, #5e, #fa, #8f, #e3, #d3, #cb, #d3
db #d3, #01, #0a, #00, #17, #d8, #17, #e0
db #5b, #e1, #59, #e9, #a8, #d9, #ff, #e1
db #36, #c3, #3f, #c3, #36, #cb, #3a, #e3
db #1e, #d3, #aa, #20, #20, #76, #4f, #ca
db #e9, #f7, #12, #ea, #3d, #cb, #3e, #eb
db #88, #c3, #85, #eb, #d3, #88, #4b, #08
db #80, #a5, #e0, #0b, #d9, #3f, #e3, #3d
db #fb, #87, #e0, #8c, #c3, #83, #eb, #83
db #fb, #68, #1f, #29, #21, #f0, #0c, #d1
db #e9, #c2, #25, #01, #e8, #68, #c0, #43
db #f6, #e0, #c8, #03, #59, #a2, #68, #d0
db #b9, #d8, #d1, #a8, #82, #9d, #fa, #36
db #d3, #9a, #db, #36, #e3, #1a, #eb, #82
db #eb, #f0, #14, #46, #71, #d0, #c0, #f0
db #10, #d1, #fd, #e1, #fd, #e9, #88, #db
db #88, #b8, #8b, #fd, #a6, #f3, #92, #85
db #fb, #89, #a8, #8b, #69, #db, #c3, #da
db #cb, #db, #a6, #da, #d3, #db, #da, #db
db #dd, #85, #a4, #52, #03, #00, #71, #d8
db #ff, #e9, #89, #cb, #00, #10, #aa, #68
db #e8, #be, #f0, #ed, #ca, #d2, #1a, #f2
db #f0, #fa, #3a, #eb, #a0, #f3, #22, #fb
db #8a, #c3, #88, #cb, #87, #d3, #89, #5e
db #db, #93, #f3, #d8, #db, #50, #01, #a9
db #32, #e8, #40, #f5, #88, #72, #79, #d9
db #51, #ca, #e0, #0a, #37, #4a, #f8, #10
db #f1, #ad, #f9, #50, #fa, #88, #d3, #93
db #d3, #86, #db, #93, #aa, #86, #89, #03
db #58, #02, #00, #bc, #c8, #3f, #eb, #c1
db #2a, #c3, #bd, #ea, #a0, #e9, #c1, #0e
db #0d, #e9, #c0, #0b, #f3, #14, #e1, #d8
db #c1, #e8, #60, #c9, #01, #e2, #41, #c3
db #93, #e3, #91, #fb, #93, #fb, #d8, #c3
db #d9, #c3, #d6, #db, #10, #19, #68, #bc
db #e0, #59, #e1, #a8, #d1, #00, #ca, #ff
db #d9, #9c, #28, #ed, #c2, #ed, #ea, #e4
db #db, #35, #15, #aa, #db, #8a, #e3, #d6
db #c3, #d7, #d1, #7a, #d9, #db, #a4, #70
db #7a, #e0, #08, #13, #7f, #e0, #a9, #f0
db #07, #d9, #5a, #f1, #a9, #e1, #4a, #f2
db #4a, #fa, #9a, #c2, #9a, #ca, #9a, #d2
db #9a, #da, #9a, #e2, #9a, #ea, #ea, #da
db #ea, #e2, #b2, #f7, #ed, #9d, #0e, #d3
db #8a, #cb, #16, #06, #82, #e8, #bc, #f0
db #61, #63, #e8, #86, #a0, #e8, #e4, #f8
db #89, #c3, #48, #f3, #a5, #0c, #c1, #9a
db #f2, #e1, #03, #e8, #10, #d9, #a0, #ea
db #89, #eb, #2c, #e0, #69, #e1, #85, #f3
db #21, #04, #13, #0a, #e9, #a8, #f1, #eb
db #ea, #84, #cb, #0e, #07, #00, #5a, #d9
db #5f, #f1, #5f, #f9, #ae, #d9, #9f, #fa
db #39, #db, #85, #e3, #0c, #66, #09, #5a
db #e9, #a9, #85, #ac, #ae, #e1, #9a, #fa
db #39, #e3, #39, #eb, #39, #f3, #39, #20
db #98, #0b, #a8, #f9, #30, #80, #a0, #00
db #d2, #d7, #c3, #e0, #c3, #d0, #05, #6b
db #fe, #f1, #8a, #f3, #86, #fb, #8a, #fb
db #d9, #96, #83, #cb, #ff, #f1, #90, #04
db #86, #dd, #f9, #ea, #fa, #3f, #f3, #8e
db #d3, #94, #07, #62, #4d, #c2, #c3, #1b
db #4e, #c2, #86, #d3, #24, #83, #e9, #ca
db #d7, #cb, #02, #98, #db, #49, #da, #84
db #86, #a0, #f2, #ea, #d2, #04, #21, #ea
db #ea, #c2, #ea, #ca, #12, #b5, #ec, #29
db #b0, #0a, #a0, #e9, #ca, #8f, #4e, #87
db #eb, #87, #f3, #87, #fb, #d4, #c3, #d6
db #cb, #d6, #d3, #d8, #d3, #d7, #db, #70
db #09, #d3, #02, #d2, #88, #f3, #8f, #f3
db #88, #fb, #8f, #b6, #92, #a9, #a2, #cb
db #cb, #df, #ba, #1c, #9e, #62, #e2, #5a
db #01, #32, #f2, #81, #6b, #1a, #3a, #c3
db #3a, #cb, #86, #8b, #69, #39, #58, #e1
db #ea, #eb, #b4, #f3, #3f, #79, #d4, #d3
db #49, #c9, #e6, #85, #c3, #41, #f7, #87
db #db, #92, #63, #e3, #c2, #9a, #0b, #89
db #e3, #de, #89, #18, #eb, #8a, #eb, #d2
db #cb, #da, #00, #00, #80
delta02:
db #5e, #41, #0f, #32, #00, #18, #c0, #1b
db #fd, #a0, #c8, #6a, #d0, #18, #e0, #1b
db #f8, #bd, #f0, #0b, #c9, #d9, #a8, #e1
db #1a, #f1, #5b, #c9, #5f, #f1, #ae, #d9
db #48, #da, #aa, #e2, #ea, #a2, #f2, #fa
db #4f, #8a, #98, #c2, #a0, #ca, #29, #d2
db #9f, #fa, #e7, #ea, #e6, #f2, #e7, #a6
db #e6, #fa, #e7, #36, #c3, #37, #a9, #3a
db #a6, #36, #cb, #3a, #35, #d3, #36, #a9
db #3a, #a6, #35, #db, #36, #35, #e3, #3f
db #82, #35, #eb, #3a, #f3, #8a, #fb, #3b
db #69, #8a, #c3, #8b, #a5, #83, #eb, #84
db #a4, #83, #fb, #2d, #05, #00, #20, #c0
db #2a, #d0, #0b, #d1, #3a, #e3, #d3, #c3
db #c0, #0e, #00, #22, #c0, #c8, #86, #e0
db #bd, #c8, #c1, #e0, #9a, #f0, #51, #ca
db #0a, #da, #43, #f3, #85, #a2, #93, #e3
db #c3, #80, #a5, #cb, #d5, #db, #03, #12
db #00, #17, #c8, #6a, #d0, #68, #c8, #b9
db #d0, #be, #d8, #0a, #d9, #5a, #c9, #59
db #f9, #47, #c2, #ca, #aa, #d2, #da, #aa
db #e2, #ea, #a8, #f2, #4a, #fa, #e9, #ca
db #e5, #f2, #3c, #0b, #00, #20, #c8, #a0
db #f0, #69, #f8, #70, #f0, #bc, #c8, #c0
db #e0, #3e, #eb, #0a, #f3, #8d, #c3, #84
db #e3, #d4, #c3, #40, #05, #00, #22, #d0
db #51, #aa, #d8, #51, #d2, #86, #f3, #df
db #c3, #21, #07, #00, #17, #d8, #0e, #c1
db #a8, #e9, #97, #c2, #e6, #ea, #e9, #f2
db #3f, #c3, #f0, #1f, #00, #21, #e0, #e8
db #a6, #f0, #f8, #71, #c0, #aa, #d8, #e0
db #9a, #f8, #10, #d9, #00, #a6, #e9, #fe
db #d1, #50, #fa, #98, #e2, #a0, #ea, #87
db #c3, #8b, #cb, #8f, #d3, #93, #86, #db
db #8a, #22, #e3, #8f, #9a, #85, #eb, #8f
db #1a, #8e, #fb, #d3, #d3, #d8, #6a, #d7
db #db, #d8, #db, #a8, #e2, #a0, #00, #2e
db #a0, #17, #e8, #68, #d8, #1a, #e0, #72
db #f0, #bc, #e0, #b9, #f8, #59, #e1, #5b
db #e9, #a8, #d1, #ff, #d9, #49, #da, #aa
db #e2, #ea, #22, #f2, #4a, #6a, #fa, #9a
db #c2, #ca, #aa, #d2, #da, #aa, #e2, #ea
db #a9, #f2, #aa, #fa, #ea, #c2, #ca, #a1
db #d2, #29, #da, #38, #cb, #39, #e3, #38
db #f3, #87, #eb, #94, #a6, #87, #f3, #94
db #87, #fb, #88, #a9, #92, #aa, #d7, #c3
db #d8, #e0, #28, #d6, #cb, #2a, #d3, #d1
db #db, #d6, #e1, #84, #d8, #87, #06, #00
db #18, #f0, #0e, #e1, #ff, #f1, #35, #53
db #a0, #db, #83, #02, #07, #0a, #00, #18
db #f8, #69, #f8, #be, #d0, #0a, #c9, #a9
db #c9, #4e, #ca, #9d, #d2, #e5, #fa, #35
db #c3, #85, #c3, #61, #02, #00, #68, #c0
db #e9, #ea, #1e, #08, #a8, #6b, #12, #5f
db #d9, #ae, #e1, #e8, #ca, #ec, #ea, #3e
db #d3, #8c, #c3, #84, #cb, #e0, #18, #00
db #71, #20, #4e, #c8, #0d, #e9, #10, #f9
db #fd, #c1, #01, #c2, #fe, #e1, #fd, #f1
db #a0, #c2, #38, #eb, #85, #4f, #38, #db
db #89, #fd, #a2, #e3, #93, #89, #eb, #8a
db #fb, #91, #a2, #93, #d9, #c3, #8a, #cb
db #de, #68, #d4, #db, #d9, #6a, #80, #0e
db #00, #72, #c8, #d0, #86, #d8, #11, #c9
db #60, #d9, #96, #e1, #af, #d9, #01, #e2
db #a1, #ca, #92, #d2, #ed, #f2, #39, #fb
db #86, #81, #68, #d6, #c3, #01, #13, #00
db #68, #d0, #be, #e0, #b9, #e8, #0a, #e9
db #5a, #c1, #59, #f1, #a8, #e1, #89, #f1
db #f7, #aa, #f9, #97, #ca, #d2, #2b, #da
db #e7, #eb, #62, #53, #fa, #39, #b2, #39
db #4b, #78, #eb, #78, #06, #0b, #69, #e8
db #70, #f8, #c0, #c0, #32, #e1, #8a, #f3
db #d2, #d3, #d0, #c7, #16, #71, #f0, #fe
db #f9, #88, #c3, #86, #cb, #87, #cb, #87
db #a9, #d5, #ca, #dc, #d3, #68, #93, #c1
db #83, #ae, #f1, #a4, #01, #e5, #f3, #c5
db #4e, #0d, #00, #bc, #d8, #0c, #d1, #11
db #d1, #5a, #f9, #af, #d1, #a9, #e9, #f8
db #e9, #49, #fa, #99, #c2, #99, #ca, #ea
db #f2, #38, #d3, #39, #db, #a0, #03, #bd
db #4e, #d8, #fe, #d9, #86, #d3, #42, #9d
db #1e, #bd, #e0, #e7, #e2, #12, #f3, #4a
db #bc, #e8, #e9, #c2, #d2, #04, #aa, #bd
db #e8, #d8, #e8, #6d, #e1, #e3, #09, #cf
db #98, #be, #e8, #3f, #eb, #06, #69, #f8
db #35, #fb, #4a, #01, #e1, #0c, #c1, #43
db #c5, #a0, #0e, #d1, #e9, #d2, #e9, #e0
db #e9, #e2, #0c, #0a, #e1, #e6, #d9, #af
db #c9, #ae, #53, #f1, #f8, #f9, #aa, #d2
db #da, #93, #fa, #89, #cb, #84, #d3, #0e
db #80, #19, #28, #0c, #e1, #5a, #e1, #af
db #c1, #a9, #d1, #48, #c2, #48, #ca, #48
db #d2, #99, #e2, #0a, #ea, #36, #e3, #35
db #f3, #37, #f3, #84, #e1, #86, #2b, #88
db #0b, #e9, #5a, #d9, #4b, #62, #5b, #c1
db #84, #f3, #c2, #06, #60, #c1, #a0, #f2
db #89, #d3, #d9, #d3, #84, #21, #89, #60
db #d1, #89, #c3, #48, #29, #5a, #f1, #ee
db #c2, #10, #0c, #75, #ac, #d9, #97, #e2
db #97, #f2, #97, #9a, #ed, #c2, #e7, #ca
db #4e, #d2, #37, #e3, #37, #eb, #39, #81
db #0d, #c3, #87, #e3, #30, #b6, #ff, #a8
db #15, #50, #a7, #f9, #9c, #fa, #d7, #cb
db #da, #d3, #da, #db, #60, #03, #00, #b1
db #f9, #86, #eb, #da, #c3, #70, #05, #00
db #fd, #e9, #98, #ea, #86, #e3, #8a, #d7
db #d3, #90, #21, #4b, #ac, #fe, #f1, #37
db #fb, #85, #34, #9e, #99, #4d, #c2, #83
db #f7, #ae, #4e, #16, #1b, #4a, #4d, #ca
db #ec, #e2, #e1, #07, #5a, #50, #f2, #86
db #c3, #88, #cb, #8a, #5e, #8a, #eb, #d4
db #cb, #d4, #d3, #c1, #c7, #72, #98, #da
db #02, #2f, #18, #99, #f2, #99, #fa, #c3
db #9a, #98, #fa, #3d, #e1, #20, #85, #ac
db #ed, #ca, #ed, #d2, #39, #96, #9a, #0f
db #e8, #d2, #38, #89, #5a, #b8, #e8, #e2
db #58, #d7, #d7, #fa, #3f, #07, #0f, #35
db #53, #9e, #e1, #3f, #cb, #18, #f6, #88
db #d3, #0d, #9e, #3a, #eb, #92, #b5, #1b
db #36, #f3, #88, #db, #25, #8e, #f3, #fb
db #1c, #f7, #62, #85, #cb, #a5, #ac, #8a
db #b4, #86, #e1, #8c, #cb, #8f, #db, #85
db #9a, #d8, #cb, #d5, #e2, #b0, #bb, #88
db #d3, #8b, #eb, #52, #8e, #f3, #f3, #29
db #b1, #a8, #84, #38, #aa, #d5, #24, #2b
db #da, #00, #00, #80
delta03:
db #01, #79, #3d, #80, #0f, #00, #22, #c0
db #22, #c8, #22, #d0, #22, #d8, #22, #e0
db #11, #d1, #51, #da, #a1, #da, #ea, #f2
db #41, #c3, #39, #e3, #38, #f3, #86, #d3
db #94, #eb, #e3, #c3, #07, #bf, #aa, #17
db #c8, #18, #17, #05, #e0, #68, #c0, #b9
db #d0, #b9, #d8, #be, #d8, #0a, #d9, #5b
db #d1, #5a, #d9, #ad, #f1, #47, #d2, #97
db #fd, #2a, #da, #35, #cb, #2d, #0c, #00
db #1b, #c8, #d0, #a9, #e0, #0a, #e8, #6b
db #c0, #9f, #d2, #3f, #cb, #3a, #fb, #8a
db #c3, #8c, #06, #84, #d3, #d2, #c3, #e0
db #0e, #00, #21, #c8, #95, #a0, #d0, #71
db #f8, #fd, #d1, #50, #fa, #87, #c3, #93
db #d3, #92, #f3, #e2, #c3, #88, #cb, #d9
db #d3, #de, #a0, #db, #e2, #6a, #3c, #08
db #00, #20, #d0, #70, #f8, #c0, #c0, #d0
db #a8, #d8, #4a, #e8, #bd, #f0, #4d, #ca
db #03, #0c, #00, #17, #d8, #01, #a0, #e0
db #68, #d0, #be, #e0, #b9, #e8, #0e, #d1
db #0a, #e9, #5a, #d1, #a9, #c1, #e7, #c2
db #5a, #ca, #e6, #f2, #0d, #03, #00, #1b
db #d8, #0b, #f9, #3a, #f3, #f0, #13, #00
db #21, #d8, #10, #c1, #68, #e1, #0d, #e9
db #69, #f1, #fd, #c1, #50, #c2, #6a, #f2
db #3e, #eb, #93, #c3, #87, #cb, #d3, #82
db #e3, #88, #eb, #8a, #96, #88, #f3, #92
db #fb, #d3, #cb, #d5, #84, #28, #01, #0e
db #00, #17, #e8, #68, #e0, #be, #e8, #b9
db #f8, #0e, #c1, #5b, #e9, #ac, #d9, #a8
db #e9, #ff, #2a, #ed, #ea, #39, #db, #f3
db #94, #a1, #fb, #84, #c3, #4b, #06, #00
db #18, #e8, #0b, #e1, #0d, #f9, #3a, #c3
db #a2, #cb, #84, #e3, #00, #2b, #22, #e8
db #a9, #f0, #aa, #f8, #72, #c0, #c8, #a1
db #d0, #29, #d8, #0a, #f9, #59, #e9, #a8
db #d9, #ae, #f9, #b1, #69, #ff, #e1, #f8
db #e9, #f7, #f1, #f8, #a6, #f7, #f9, #f8
db #47, #c2, #48, #2a, #ca, #a1, #d2, #aa
db #da, #49, #fa, #99, #c2, #ca, #69, #d2
db #ed, #ca, #28, #d2, #ea, #e2, #e6, #ea
db #ea, #69, #e5, #f2, #38, #d3, #42, #1a
db #38, #fb, #89, #c3, #93, #f3, #85, #fb
db #d6, #c3, #d7, #cb, #d8, #40, #a5, #d7
db #d3, #10, #0b, #00, #17, #f0, #be, #f8
db #4e, #c2, #9c, #fa, #ec, #c2, #e5, #fa
db #3f, #c3, #37, #f3, #d8, #d3, #df, #a2
db #db, #db, #0f, #1f, #00, #18, #f0, #f8
db #69, #9a, #bc, #c8, #bd, #0a, #be, #d0
db #0a, #c9, #29, #d1, #0b, #29, #0e, #e1
db #5b, #c1, #5f, #c9, #aa, #d9, #47, #da
db #e2, #28, #ea, #4f, #a6, #47, #f2, #fa
db #97, #c2, #80, #aa, #ca, #98, #da, #3a
db #db, #36, #e3, #37, #38, #a7, #3e, #35
db #f3, #3d, #20, #81, #da, #db, #d2, #cb
db #1e, #7e, #29, #f0, #20, #77, #6b, #c8
db #0b, #e9, #0f, #f9, #4f, #f5, #67, #5d
db #fb, #d2, #04, #6a, #1b, #f0, #11, #70
db #f8, #a0, #f2, #88, #e3, #58, #02, #00
db #71, #c0, #4d, #c2, #05, #01, #00, #68
db #c8, #87, #0b, #00, #6d, #c8, #6c, #d0
db #b9, #c0, #bd, #e8, #ad, #e9, #9d, #d2
db #e8, #e2, #e8, #ea, #35, #d3, #84, #cb
db #8a, #cb, #78, #ab, #85, #70, #71, #c2
db #fd, #fd, #a0, #c8, #1e, #f0, #0c, #d1
db #00, #e2, #fd, #f9, #3e, #f3, #8a, #e3
db #8f, #eb, #d5, #d3, #21, #07, #8f, #4e
db #d8, #ac, #d1, #a8, #f9, #e7, #a1, #4e
db #f2, #39, #d3, #88, #fb, #34, #5e, #87
db #e8, #85, #cb, #30, #06, #95, #f3, #5d
db #f0, #bc, #e8, #97, #ea, #e7, #ea, #e0
db #c3, #df, #db, #68, #03, #00, #c1, #d0
db #86, #cb, #89, #e3, #40, #06, #00, #c1
db #d8, #bd, #e0, #fe, #c1, #fe, #c9, #44
db #fb, #e3, #cb, #08, #0d, #00, #bc, #e0
db #c8, #f0, #0c, #d9, #60, #e9, #aa, #c1
db #a9, #f1, #48, #e2, #48, #ea, #99, #da
db #99, #e2, #ed, #f2, #38, #db, #89, #ca
db #23, #b3, #be, #96, #0c, #e1, #99, #fa
db #ee, #c2, #2c, #24, #af, #bd, #f8, #38
db #eb, #0e, #0c, #0c, #b3, #64, #23, #e4
db #e1, #98, #c2, #98, #ca, #e8, #fd, #39
db #d2, #ec, #ea, #3d, #cb, #36, #eb, #35
db #fb, #85, #c3, #c0, #6d, #42, #60, #d1
db #60, #d9, #01, #c2, #fe, #d1, #01, #da
db #fe, #e1, #fe, #f1, #42, #db, #93, #cb
db #86, #fb, #93, #fb, #dc, #d3, #d9, #db
db #0c, #08, #9a, #5a, #f1, #a9, #70, #48
db #f2, #48, #f2, #e3, #3f, #13, #0c, #c3
db #89, #d3, #86, #9e, #1d, #a9, #d1, #48
db #53, #1e, #af, #d1, #89, #db, #1c, #f3
db #48, #ae, #e9, #85, #d3, #28, #01, #8e
db #f1, #c3, #17, #07, #ad, #f9, #98, #e2
db #98, #ea, #8a, #db, #20, #9a, #eb, #f8
db #c1, #9c, #4a, #97, #f2, #39, #eb, #70
db #07, #42, #00, #da, #ff, #f1, #86, #eb
db #db, #c3, #da, #d3, #d7, #db, #8c, #db
db #d0, #b8, #5f, #fd, #37, #3a, #e9, #fd
db #f1, #60, #ee, #48, #e9, #86, #f3, #da
db #cb, #16, #96, #ff, #f9, #97, #e2, #e9
db #da, #06, #20, #a6, #47, #ca, #e9, #ca
db #d2, #83, #01, #ab, #4e, #96, #98, #31
db #4d, #d2, #c1, #b9, #98, #15, #e3, #cb
db #18, #05, #f3, #fa, #ea, #08, #47, #a7
db #db, #3f, #d5, #c3, #5a, #a8, #e7, #9f
db #76, #e8, #c2, #e8, #da, #8a, #39, #39
db #f3, #24, #b5, #e5, #e9, #c2, #81, #7d
db #88, #e7, #e2, #37, #eb, #36, #f3, #43
db #aa, #e9, #ea, #f2, #22, #fa, #09, #78
db #ed, #fa, #90, #af, #58, #3f, #d3, #85
db #db, #88, #db, #87, #eb, #91, #fb, #02
db #9e, #36, #fb, #c2, #f7, #ae, #37, #a4
db #c9, #0a, #86, #c3, #e1, #03, #a6, #88
db #d3, #c3, #d4, #b8, #92, #3d, #62, #86
db #db, #84, #fb, #b0, #02, #8f, #db, #85
db #e3, #85, #eb, #dc, #cb, #8e, #db, #a0
db #01, #62, #93, #eb, #b4, #62, #85, #f3
db #69, #c2, #8a, #9f, #38, #fb, #a5, #e9
db #c0, #d4, #cb, #00, #20
delta04:
db #00, #ea, #3e, #16, #02, #00, #17, #c0
db #e7, #d2, #0f, #3b, #00, #20, #c0, #18
db #c8, #1b, #fd, #d0, #a0, #d8, #a6, #e0
db #18, #e8, #1b, #68, #c0, #6b, #92, #6d
db #c8, #6c, #d0, #b9, #c0, #86, #c8, #bc
db #d0, #0b, #e9, #81, #a6, #f9, #5f, #e9
db #5c, #f9, #ad, #e9, #4f, #c2, #ca, #4d
db #da, #a8, #e2, #6a, #ea, #9f, #d2, #e8
db #c2, #ca, #86, #d2, #3a, #cb, #3e, #d3
db #92, #db, #3a, #e3, #36, #eb, #37, #9a
db #36, #f3, #37, #aa, #3a, #3e, #6a, #35
db #fb, #36, #37, #aa, #3a, #3c, #9a, #8a
db #c3, #8c, #a6, #8d, #85, #cb, #8b, #a8
db #8c, #a0, #84, #d3, #a6, #e3, #82, #fb
db #84, #d2, #c3, #d3, #a9, #d4, #a1, #d2
db #d3, #d3, #a6, #f0, #1c, #00, #21, #c0
db #f8, #70, #f0, #9a, #f8, #c0, #c0, #6a
db #f8, #fd, #e1, #e9, #96, #f1, #a0, #c2
db #43, #fb, #93, #cb, #82, #d3, #8f, #db
db #93, #86, #87, #eb, #86, #fb, #88, #9a
db #d5, #c3, #d7, #aa, #d8, #da, #a6, #e2
db #d8, #cb, #e2, #86, #d8, #d3, #da, #db
db #e2, #8a, #00, #40, #2a, #22, #c0, #c8
db #aa, #d0, #d8, #8a, #e0, #bd, #00, #a5
db #bc, #e8, #be, #f8, #11, #d1, #0d, #e9
db #af, #d9, #fe, #c9, #00, #ca, #fe, #d9
db #aa, #e1, #ff, #e9, #51, #c2, #47, #ca
db #d2, #a2, #da, #e2, #48, #98, #47, #ea
db #48, #a8, #f2, #29, #fa, #98, #c2, #a1
db #a6, #98, #ca, #a1, #98, #d2, #a1, #9a
db #98, #da, #9c, #a2, #a1, #98, #e2, #a8
db #ea, #6a, #f2, #97, #fa, #e9, #ca, #d2
db #a8, #e2, #29, #ea, #ea, #f2, #ed, #a2
db #e5, #fa, #ea, #35, #c3, #aa, #cb, #d3
db #88, #db, #38, #aa, #e3, #eb, #92, #f3
db #89, #cb, #87, #d3, #88, #20, #0a, #db
db #94, #eb, #92, #f3, #90, #fb, #e3, #cb
db #d7, #db, #69, #03, #00, #20, #c8, #a4
db #d0, #a0, #e0, #d0, #07, #00, #21, #c8
db #28, #d0, #37, #c3, #3f, #cb, #86, #f3
db #88, #1a, #dc, #cb, #07, #10, #00, #18
db #d8, #68, #c8, #16, #d0, #6c, #d8, #bd
db #f0, #0e, #e1, #9d, #d2, #e6, #ea, #a8
db #f2, #e1, #fa, #39, #63, #78, #e3, #35
db #eb, #87, #c3, #83, #d3, #82, #eb, #2d
db #05, #87, #36, #d8, #0c, #c9, #85, #e9
db #6b, #cb, #d2, #db, #04, #81, #7f, #ea
db #d8, #c0, #f0, #88, #eb, #8a, #f3, #3c
db #0d, #57, #e8, #41, #ea, #f0, #5f, #d9
db #5f, #f1, #e7, #e2, #e7, #ea, #e7, #fa
db #3f, #e3, #86, #c3, #8a, #cb, #8c, #d3
db #85, #eb, #d3, #db, #e1, #0e, #b1, #e8
db #39, #c8, #c0, #55, #0a, #f8, #0d, #f9
db #4e, #ca, #50, #e2, #50, #fa, #3e, #fb
db #8b, #1e, #8a, #fb, #d4, #cb, #da, #fd
db #02, #d3, #01, #0f, #00, #17, #f0, #68
db #d8, #a9, #e8, #2a, #f0, #09, #c1, #0a
db #f9, #59, #a8, #82, #99, #ea, #e7, #c2
db #80, #68, #ca, #e6, #da, #3f, #c3, #36
db #d3, #d1, #db, #b4, #06, #00, #21, #f0
db #00, #e2, #ec, #52, #3f, #fb, #8e, #c3
db #8a, #e3, #4b, #03, #00, #18, #f8, #5b
db #c9, #84, #eb, #78, #0b, #00, #20, #f8
db #70, #c0, #aa, #d0, #d8, #ae, #e0, #3d
db #a0, #e0, #38, #e8, #4d, #c2, #8d, #d3
db #8d, #db, #e0, #cf, #57, #71, #c0, #71
db #d0, #c1, #c0, #10, #e9, #4e, #c2, #50
db #ea, #50, #f2, #99, #fa, #87, #e3, #93
db #eb, #93, #f3, #0d, #2a, #09, #6b, #4a
db #e8, #e2, #38, #cb, #a4, #01, #a3, #71
db #0e, #04, #81, #e5, #eb, #d0, #0c, #e9
db #5a, #e9, #ae, #e9, #68, #03, #9d, #d8
db #90, #7f, #3a, #d0, #03, #09, #00, #68
db #e0, #b9, #f0, #5b, #e1, #ac, #d1, #a8
db #f1, #e6, #e2, #3a, #c3, #39, #e3, #39
db #eb, #60, #06, #c5, #e8, #84, #e0, #f0
db #d9, #c3, #d6, #cb, #d6, #d3, #dc, #db
db #49, #87, #ae, #bd, #d0, #c0, #0f, #c1
db #3d, #42, #d8, #60, #c1, #01, #e2, #41
db #cb, #42, #e3, #88, #cb, #93, #97, #87
db #fb, #db, #c3, #d5, #cb, #d9, #b9, #b9
db #db, #ff, #ea, #e3, #db, #40, #97, #bd
db #17, #0c, #e1, #ae, #f9, #01, #c2, #01
db #fa, #e9, #c2, #21, #96, #3d, #be, #e8
db #f8, #c1, #39, #f3, #30, #21, #0a, #bc
db #f0, #00, #da, #ff, #f1, #36, #cb, #39
db #fb, #87, #cb, #86, #db, #91, #fb, #df
db #d3, #80, #15, #69, #be, #f0, #c1, #f1
db #60, #e9, #fe, #61, #83, #f1, #d1, #01
db #d2, #fe, #92, #ed, #f1, #fe, #f9, #51
db #ca, #51, #91, #3a, #99, #ca, #e9, #da
db #42, #db, #44, #fb, #94, #d3, #94, #db
db #94, #e3, #89, #f3, #d5, #db, #1e, #19
db #00, #bf, #f0, #0b, #c9, #0f, #fd, #d1
db #aa, #d9, #e1, #86, #e9, #0c, #f1, #5f
db #c1, #aa, #d1, #e1, #86, #f9, #4d, #d2
db #4f, #da, #a9, #ea, #aa, #fa, #9f, #c2
db #da, #a4, #e2, #c2, #ea, #3d, #cb, #3e
db #eb, #8d, #f3, #28, #db, #8a, #e3, #52
db #01, #13, #f8, #86, #82, #f7, #0c, #c1
db #5a, #03, #25, #e6, #d1, #88, #c3, #85
db #d3, #48, #04, #ef, #d9, #ee, #c2, #83
db #db, #85, #e3, #43, #87, #eb, #0a, #e1
db #a9, #c9, #36, #23, #23, #db, #0b, #21
db #af, #a2, #0b, #e1, #3a, #eb, #84, #c2
db #02, #95, #76, #10, #f9, #a0, #f2, #87
db #07, #00, #5b, #d1, #ad, #f1, #4e, #d2
db #4d, #f2, #4d, #fa, #e8, #da, #36, #db
db #08, #05, #00, #60, #d9, #af, #d1, #51
db #f2, #9a, #6f, #3a, #fb, #06, #8f, #5a
db #22, #d1, #39, #26, #cb, #10, #09, #ab
db #a8, #00, #25, #a1, #3b, #da, #99, #e2
db #ed, #ea, #f0, #fa, #89, #db, #d7, #cb
db #d8, #db, #0c, #77, #a9, #62, #41, #f1
db #b0, #25, #e9, #fd, #c1, #86, #eb, #85
db #fb, #a0, #ee, #39, #c9, #9a, #d2, #e3
db #d3, #70, #53, #1a, #ff, #f9, #89, #c3
db #89, #d3, #d6, #c3, #e0, #c3, #dd, #db
db #df, #e0, #02, #6d, #79, #47, #f2, #47
db #fa, #97, #c2, #97, #f2, #12, #0d, #ae
db #97, #ca, #99, #81, #f2, #83, #d2, #97
db #da, #84, #65, #05, #ac, #9a, #da, #e8
db #f2, #e8, #fa, #86, #88, #8c, #01, #fb
db #04, #25, #79, #89, #97, #e2, #e9, #f2
db #e9, #fa, #24, #88, #97, #ea, #20, #1e
db #98, #fa, #ed, #ca, #ec, #d2, #ed, #d2
db #82, #e1, #79, #e7, #da, #1c, #c5, #a7
db #e8, #ea, #e7, #8a, #db, #18, #b4, #87
db #ed, #c7, #f5, #01, #82, #e1, #38, #c3
db #38, #d3, #88, #d3, #90, #9a, #37, #cb
db #86, #28, #1a, #01, #b8, #3f, #db, #e9
db #62, #3f, #eb, #3f, #f3, #41, #79, #84
db #cb, #2c, #75, #8b, #87, #db, #50, #89
db #d3, #4b, #39, #63, #47, #d3, #34, #8c
db #e1, #85, #f3, #00, #02
delta05:
db #44, #ea, #40, #0b, #06, #00, #17, #c0
db #18, #c0, #17, #c8, #17, #d0, #6b, #d0
db #5b, #c9, #3c, #1e, #00, #19, #c0, #1a
db #c8, #1f, #d0, #6f, #fd, #e8, #6a, #f8
db #bf, #c0, #c8, #aa, #d8, #e0, #aa, #e8
db #f0, #9a, #f8, #0f, #c9, #aa, #d1, #d9
db #aa, #e1, #e9, #6a, #f9, #5f, #c1, #c9
db #a8, #d1, #22, #e1, #4d, #d2, #4f, #56
db #e2, #47, #fa, #9f, #f2, #88, #d3, #8b
db #eb, #d2, #04, #00, #1a, #c0, #c0, #91
db #aa, #8f, #c3, #8b, #f3, #c3, #09, #00
db #1b, #c0, #5d, #c1, #e9, #ea, #f2, #69
db #fa, #37, #d3, #a0, #db, #89, #d3, #68
db #e3, #2d, #08, #00, #1f, #c0, #18, #e0
db #1f, #06, #0b, #d9, #4f, #da, #e8, #e2
db #3d, #c3, #86, #cb, #78, #0b, #00, #20
db #c0, #0f, #c1, #9a, #f1, #5f, #d9, #53
db #e9, #47, #f2, #40, #e3, #36, #eb, #3f
db #fb, #8e, #f3, #89, #fb, #48, #a9, #69
db #21, #62, #0c, #e9, #ae, #f9, #a0, #f2
db #4b, #16, #19, #c8, #68, #c8, #69, #f8
db #3a, #cb, #87, #c3, #8a, #8a, #8d, #e3
db #1a, #eb, #8a, #eb, #87, #10, #1e, #1b
db #c8, #19, #d0, #1b, #fd, #aa, #d8, #e0
db #46, #e8, #18, #f0, #0c, #d1, #0e, #e1
db #4e, #da, #e8, #f2, #35, #d3, #8a, #db
db #8b, #28, #d4, #cb, #0a, #d3, #1e, #1b
db #00, #1f, #c8, #1a, #d0, #22, #d8, #1f
db #a0, #e8, #8a, #f8, #68, #c0, #6f, #c8
db #a8, #d8, #06, #e0, #bf, #d0, #0b, #d1
db #5c, #c1, #4f, #ca, #4d, #da, #97, #c2
db #9f, #ca, #e8, #c2, #a8, #ca, #69, #da
db #3a, #fb, #8e, #c3, #a0, #cb, #8c, #e3
db #2a, #eb, #84, #f3, #e0, #18, #00, #20
db #c8, #d0, #a9, #d8, #a6, #e0, #70, #d8
db #e0, #10, #e1, #9a, #f9, #60, #c1, #05
db #a6, #c9, #fd, #e9, #48, #c2, #50, #da
db #a0, #ea, #41, #d3, #37, #e3, #42, #eb
db #85, #c3, #94, #d3, #db, #86, #eb, #84
db #aa, #fb, #e2, #c3, #e4, #db, #c0, #12
db #00, #21, #c8, #e0, #aa, #e8, #f0, #6a
db #f8, #71, #c0, #c8, #aa, #d0, #e0, #81
db #a4, #f8, #c1, #f0, #60, #e1, #01, #d2
db #fd, #d9, #28, #e1, #87, #db, #86, #e3
db #92, #eb, #40, #0a, #00, #21, #d0, #a2
db #d8, #71, #b1, #d1, #80, #a5, #f1, #f8
db #c9, #94, #eb, #d6, #c3, #d8, #a1, #d5
db #d3, #21, #09, #00, #17, #d8, #68, #f0
db #68, #e8, #0e, #d1, #0d, #f9, #96, #d2
db #39, #cb, #0a, #d3, #82, #db, #0f, #42
db #00, #18, #d8, #12, #f8, #6b, #c8, #68
db #d0, #6c, #d8, #b9, #d0, #bd, #98, #bc
db #d8, #be, #85, #a4, #e0, #0a, #d9, #0b
db #e1, #5b, #d1, #ac, #c1, #4d, #f2, #a6
db #fa, #9d, #c2, #97, #ca, #9d, #97, #d2
db #9d, #9a, #97, #da, #98, #a6, #9d, #97
db #e2, #98, #aa, #9d, #9f, #9a, #97, #ea
db #98, #aa, #9d, #9f, #6a, #97, #f2, #98
db #9d, #a8, #eb, #6a, #e7, #fa, #36, #c3
db #37, #3a, #9a, #35, #cb, #36, #aa, #37
db #3d, #a1, #3e, #89, #36, #d3, #3a, #eb
db #3e, #a2, #fb, #88, #c3, #cb, #89, #aa
db #8a, #8d, #9a, #8b, #d3, #8c, #a6, #8d
db #8c, #db, #8d, #9a, #8a, #e3, #8b, #68
db #82, #f3, #89, #68, #d3, #cb, #d2, #db
db #d3, #0a, #03, #12, #00, #17, #e8, #68
db #f0, #be, #0a, #b9, #f8, #0a, #f1, #12
db #f9, #5a, #d9, #5b, #e9, #a9, #c9, #a8
db #f9, #ad, #85, #a5, #f8, #c1, #96, #ca
db #ed, #c2, #eb, #ea, #39, #c3, #3f, #a6
db #39, #db, #f0, #3a, #00, #20, #e8, #f8
db #70, #c0, #a9, #d0, #aa, #e8, #c0, #c8
db #d8, #aa, #e0, #e8, #5a, #f0, #0c, #d9
db #10, #e9, #00, #da, #4a, #e2, #4f, #c2
db #48, #ca, #4e, #2a, #48, #d2, #da, #8a
db #e2, #50, #69, #48, #ea, #50, #88, #48
db #f2, #50, #0a, #fa, #40, #c3, #36, #db
db #35, #e3, #36, #28, #35, #eb, #a0, #f3
db #3f, #6a, #85, #cb, #89, #db, #88, #e3
db #86, #f3, #8f, #92, #a9, #93, #aa, #8a
db #fb, #8f, #91, #6a, #db, #c3, #de, #df
db #a9, #e0, #aa, #da, #cb, #dc, #de, #69
db #dc, #d3, #de, #aa, #db, #db, #dd, #de
db #aa, #df, #e0, #a4, #e3, #5a, #e1, #06
db #00, #20, #f0, #70, #c8, #4e, #d2, #35
db #fb, #8a, #f3, #db, #d3, #01, #1b, #00
db #17, #f8, #67, #c0, #68, #f8, #09, #c9
db #5a, #12, #5b, #f1, #a9, #c1, #ac, #e1
db #a8, #f1, #fc, #d1, #aa, #d9, #e1, #aa
db #e9, #f1, #86, #f9, #46, #fa, #96, #da
db #a9, #e2, #aa, #ea, #e7, #da, #e2, #06
db #ea, #e6, #f2, #ea, #fa, #39, #fb, #d1
db #cb, #92, #d3, #07, #10, #00, #68, #d8
db #8a, #e0, #69, #a0, #b9, #a5, #be, #e8
db #0a, #e1, #1a, #e9, #0b, #f1, #5b, #d9
db #5a, #e1, #ac, #c9, #a9, #d1, #ad, #f1
db #98, #fa, #3f, #cb, #82, #e3, #0d, #01
db #00, #6b, #d8, #84, #08, #00, #71, #e8
db #6a, #f0, #c1, #c0, #c8, #81, #28, #d0
db #0c, #c9, #9f, #fa, #88, #f3, #d0, #07
db #00, #c0, #d0, #37, #f3, #86, #db, #87
db #e3, #a5, #f3, #91, #0a, #e3, #d3, #29
db #01, #00, #bd, #d8, #80, #14, #00, #c1
db #e0, #c8, #f0, #0c, #e1, #60, #d9, #aa
db #c9, #af, #d1, #b1, #f9, #fd, #c9, #01
db #da, #a0, #e2, #a1, #ea, #51, #c2, #ee
db #68, #41, #cb, #42, #e3, #37, #eb, #43
db #f3, #94, #c3, #d5, #52, #e4, #d3, #16
db #02, #00, #b9, #e8, #3f, #d3, #08, #07
db #00, #bc, #e8, #be, #f8, #a9, #f9, #98
db #d2, #99, #e2, #a9, #ea, #a1, #f2, #00
db #2e, #a6, #bd, #e8, #bc, #f0, #bd, #11
db #c1, #5a, #96, #59, #f1, #a8, #e1, #fd
db #c1, #fe, #82, #01, #c2, #fe, #d1, #a9
db #e9, #a6, #f1, #01, #f2, #fa, #51, #ea
db #a8, #f2, #08, #fa, #99, #ca, #9a, #d2
db #99, #da, #9a, #8a, #e2, #9c, #2a, #e6
db #da, #e2, #92, #ea, #e8, #fa, #38, #c3
db #41, #8a, #38, #cb, #12, #d3, #42, #db
db #44, #fb, #84, #c3, #85, #fb, #87, #a9
db #93, #a6, #d7, #c3, #e3, #d6, #cb, #d8
db #9a, #d6, #d3, #d8, #6b, #d6, #db, #d8
db #60, #58, #3d, #4a, #bc, #f8, #00, #ca
db #94, #e3, #d5, #cb, #0e, #08, #16, #0b
db #c1, #0b, #c9, #0c, #f9, #e7, #f2, #38
db #fb, #84, #aa, #84, #84, #84, #c8, #a5
db #03, #00, #0c, #c1, #e8, #ea, #87, #d3
db #a4, #65, #2b, #0d, #c1, #10, #11, #09
db #48, #cb, #cd, #e9, #0d, #f7, #a8, #75
db #c1, #46, #c9, #ff, #e1, #ff, #e9, #4c
db #c2, #47, #ca, #47, #d2, #46, #f2, #96
db #f2, #e7, #d2, #3f, #db, #38, #eb, #9e
db #c3, #0c, #05, #9b, #13, #f1, #5a, #f9
db #99, #fa, #ec, #e2, #37, #fb, #86, #80
db #81, #a7, #5a, #e9, #e9, #da, #e9, #e2
db #2c, #02, #5f, #f9, #a9, #0c, #b1, #86
db #37, #a9, #d9, #d4, #db, #30, #27, #b0
db #f9, #00, #ac, #8b, #da, #70, #f7, #83
db #eb, #c2, #9d, #fa, #20, #9a, #df, #f8
db #d1, #98, #e0, #9c, #97, #32, #c2, #ec
db #ca, #f0, #fa, #d7, #d3, #a0, #d3, #0a
db #fd, #d1, #ed, #f2, #28, #cb, #e2, #a4
db #70, #0c, #1c, #00, #d2, #00, #ea, #39
db #eb, #36, #f3, #93, #eb, #d9, #c3, #d9
db #cb, #df, #cb, #d9, #d3, #df, #d3, #d9
db #db, #e1, #db, #b0, #96, #c3, #47, #ea
db #8f, #d3, #d5, #db, #1c, #27, #97, #fa
db #34, #9a, #47, #e7, #c2, #88, #c9, #04
db #5d, #aa, #e9, #c2, #e9, #38, #79, #88
db #fb, #05, #d3, #b8, #e7, #ca, #24, #e1
db #38, #d2, #ec, #da, #02, #5d, #a0, #ed
db #da, #ac, #e2, #ed, #ea, #84, #68, #86
db #e7, #ed, #fa, #87, #cb, #43, #21, #e7
db #35, #c3, #86, #c3, #61, #b7, #39, #e3
db #50, #9c, #c5, #3f, #e3, #39, #70, #0b
db #23, #5f, #f3, #90, #86, #d3, #3f, #eb
db #90, #fb, #14, #27, #38, #f3, #82, #98
db #51, #36, #fb, #69, #86, #89, #c3, #86
db #d3, #b4, #20, #62, #8f, #cb, #8e, #d3
db #8e, #db, #8f, #db, #49, #79, #85, #d3
db #0a, #bf, #89, #84, #db, #25, #89, #85
db #db, #96, #ae, #8a, #db, #da, #a1, #e8
db #8c, #e3, #83, #98, #ed, #88, #eb, #92
db #9b, #8e, #fb, #dc, #00, #00, #80
delta06:
db #5e, #44, #0f, #30, #00, #17, #c0, #18
db #fd, #6a, #f0, #68, #c0, #d8, #88, #e0
db #69, #82, #f8, #b9, #d8, #bd, #8a, #b9
db #e0, #a6, #e8, #f0, #0a, #e1, #a8, #e9
db #a2, #f1, #0b, #5b, #c9, #a8, #d9, #68
db #e1, #ac, #c9, #ae, #e1, #6a, #e9, #ad
db #f1, #fa, #c1, #c9, #a8, #d1, #2a, #d9
db #48, #d2, #49, #4a, #9a, #49, #da, #4a
db #a6, #4b, #4a, #e2, #4b, #a9, #9c, #a1
db #e9, #d2, #ea, #a4, #37, #d3, #36, #db
db #37, #aa, #89, #c3, #8a, #d3, #8b, #db
db #eb, #06, #f3, #8a, #fb, #f0, #9b, #00
db #19, #c0, #1a, #a9, #20, #aa, #19, #c8
db #1a, #1f, #a6, #20, #1f, #d0, #20, #9a
db #1f, #d8, #20, #2a, #1f, #e0, #e8, #8a
db #f0, #20, #0a, #1f, #f8, #6f, #c0, #aa
db #c8, #d0, #29, #d8, #70, #a2, #6f, #e0
db #70, #6f, #e8, #a8, #f0, #28, #f8, #bf
db #c0, #c0, #aa, #bf, #c8, #d0, #aa, #d8
db #e0, #aa, #e8, #f0, #6a, #f8, #0f, #c1
db #c9, #aa, #d1, #d9, #aa, #e1, #e9, #a9
db #f1, #aa, #f9, #5f, #c1, #c9, #aa, #d1
db #d9, #aa, #e9, #f1, #68, #f9, #f7, #e1
db #88, #e9, #fa, #a6, #f1, #fb, #fa, #f9
db #fb, #a9, #fc, #aa, #4a, #c2, #4b, #4c
db #6a, #4a, #ca, #4b, #4c, #2a, #d2, #2a
db #da, #4d, #4f, #8a, #4d, #e2, #aa, #ea
db #f2, #4a, #fa, #99, #da, #98, #e2, #99
db #a6, #9a, #98, #ea, #99, #aa, #9a, #9b
db #9a, #98, #f2, #99, #aa, #9a, #9b, #6a
db #98, #fa, #99, #9a, #a9, #9b, #aa, #e8
db #c2, #e9, #ea, #a6, #eb, #e8, #ca, #e9
db #aa, #ea, #eb, #ba, #ed, #f9, #d2, #86
db #da, #ed, #e2, #ec, #ea, #ea, #fd, #f2
db #0a, #fa, #39, #d3, #3a, #a6, #3b, #39
db #db, #3a, #aa, #3b, #3c, #a9, #3d, #aa
db #39, #e3, #3a, #3b, #aa, #3c, #3d, #6a
db #36, #eb, #39, #3a, #aa, #3b, #3c, #a9
db #3d, #aa, #36, #f3, #39, #3a, #aa, #3b
db #3c, #a6, #3d, #35, #fb, #36, #aa, #39
db #3b, #aa, #3c, #3d, #a9, #3e, #aa, #85
db #c3, #8c, #8d, #a6, #8e, #8d, #cb, #94
db #9a, #85, #d3, #8d, #aa, #8e, #94, #68
db #8d, #db, #8e, #1a, #93, #eb, #91, #f3
db #89, #fb, #da, #d3, #db, #69, #dc, #db
db #e4, #08, #87, #0f, #00, #1b, #c0, #1d
db #c8, #19, #d8, #68, #c8, #6b, #91, #28
db #d0, #0e, #f1, #f8, #c9, #fa, #e1, #49
db #e2, #4b, #ea, #e8, #d2, #ea, #da, #e8
db #e2, #35, #e3, #1e, #20, #00, #1c, #c0
db #18, #c8, #a2, #d0, #1e, #18, #d8, #8a
db #e0, #1a, #a4, #1e, #a6, #1a, #e8, #1e
db #f8, #6e, #c0, #d8, #be, #c0, #aa, #c8
db #d8, #a9, #e8, #a1, #f0, #0b, #c1, #a5
db #d9, #5c, #c9, #5e, #d1, #a1, #e9, #ae
db #d1, #aa, #f9, #f8, #e1, #a0, #e9, #4b
db #f2, #9c, #ca, #5a, #d2, #8b, #c3, #84
db #db, #8a, #f3, #e1, #09, #00, #1d, #c0
db #19, #d0, #0d, #c1, #f7, #f1, #fc, #0a
db #f7, #f9, #47, #c2, #a0, #ca, #1a, #d2
db #78, #0e, #00, #1f, #c0, #1a, #d0, #5f
db #e1, #00, #fa, #4b, #d2, #4c, #e2, #4f
db #6a, #98, #da, #9a, #9c, #aa, #e8, #eb
db #81, #aa, #36, #e3, #8c, #cb, #80, #34
db #00, #21, #c0, #c8, #08, #d0, #16, #d8
db #21, #aa, #e0, #e8, #a9, #f0, #aa, #f8
db #71, #c0, #d0, #aa, #d8, #e0, #aa, #e8
db #f0, #6a, #f8, #c1, #c0, #c8, #aa, #d0
db #d8, #a9, #e8, #a5, #f0, #07, #e9, #29
db #f9, #60, #e1, #aa, #d1, #af, #d9, #fb
db #c1, #01, #c2, #f9, #c9, #fb, #28, #fd
db #d1, #01, #d2, #f9, #d9, #a4, #e1, #fd
db #a6, #01, #fa, #49, #ea, #9e, #c2, #ca
db #9d, #da, #8a, #e2, #9e, #16, #ec, #ca
db #41, #d3, #42, #eb, #8f, #db, #86, #eb
db #94, #81, #62, #e2, #c3, #d8, #cb, #e5
db #db, #1a, #01, #00, #17, #c8, #3c, #0e
db #00, #1e, #c8, #1a, #d8, #1e, #0a, #e8
db #48, #ea, #4c, #a2, #4f, #48, #f2, #86
db #fa, #98, #c2, #e7, #e2, #95, #aa, #ea
db #3a, #fb, #8c, #d3, #d2, #01, #00, #17
db #d0, #43, #09, #00, #17, #d8, #e8, #29
db #f0, #68, #29, #5a, #e1, #ac, #d9, #e6
db #e2, #5a, #ea, #88, #eb, #61, #02, #00
db #17, #e0, #47, #e2, #07, #1a, #00, #1b
db #e0, #aa, #e8, #f0, #56, #f8, #6b, #c0
db #68, #d0, #6b, #d8, #68, #e8, #b9, #f8
db #09, #c1, #0e, #e1, #86, #e9, #0a, #f9
db #5a, #c1, #8a, #e9, #5b, #10, #28, #ac
db #d1, #a9, #d9, #ad, #f9, #4b, #fa, #9b
db #c2, #e6, #f2, #3a, #c3, #35, #d3, #88
db #c3, #8b, #cb, #20, #03, #00, #17, #f8
db #fb, #e9, #86, #c3, #03, #10, #00, #67
db #c0, #68, #f8, #b8, #c0, #09, #c9, #52
db #d1, #0e, #d9, #5a, #c9, #5b, #f1, #a9
db #d1, #ac, #e1, #fd, #c1, #9b, #ca, #e6
db #fa, #86, #cb, #89, #db, #d1, #86, #e0
db #20, #00, #70, #c0, #96, #c8, #10, #e9
db #4d, #d2, #50, #e2, #aa, #ea, #f2, #9a
db #fa, #a0, #c2, #02, #ca, #9b, #e2, #9e
db #f2, #ec, #c2, #ee, #8a, #3e, #cb, #aa
db #d3, #db, #28, #e3, #41, #a2, #3e, #eb
db #f3, #42, #82, #43, #fb, #8e, #cb, #a8
db #e3, #a6, #eb, #91, #8e, #f3, #93, #82
db #8e, #fb, #de, #c3, #96, #d3, #01, #11
db #00, #67, #c8, #b8, #2a, #d0, #56, #d8
db #09, #e1, #5a, #d1, #5b, #f9, #ab, #c1
db #ac, #e9, #9b, #d2, #9f, #e2, #91, #a1
db #ea, #36, #c3, #90, #d3, #82, #db, #88
db #f3, #d1, #c3, #00, #61, #a0, #71, #c8
db #bc, #f8, #c1, #aa, #11, #c9, #0d, #e9
db #f1, #11, #a0, #f9, #b1, #d1, #a8, #e9
db #b1, #f1, #b0, #f9, #fc, #c1, #00, #c2
db #fd, #c9, #00, #d2, #fd, #d9, #ff, #e1
db #29, #e9, #01, #ea, #ff, #f1, #fe, #f9
db #ff, #a6, #4e, #c2, #51, #4e, #ca, #51
db #9a, #4e, #d2, #51, #62, #4e, #da, #51
db #69, #e2, #46, #f2, #aa, #fa, #96, #c2
db #ca, #a2, #d2, #da, #97, #9a, #96, #e2
db #97, #69, #96, #ea, #97, #88, #96, #f2
db #97, #a5, #fa, #9f, #a6, #e7, #c2, #ec
db #d2, #35, #c3, #3f, #35, #cb, #3f, #a8
db #41, #aa, #3f, #d3, #db, #28, #e3, #42
db #6a, #3f, #eb, #37, #f3, #3f, #43, #9a
db #37, #fb, #3f, #69, #87, #c3, #8f, #89
db #8a, #cb, #8f, #29, #d3, #90, #db, #8f
db #e3, #90, #a6, #8f, #eb, #90, #86, #f3
db #87, #aa, #8f, #90, #9a, #8f, #fb, #90
db #6a, #d6, #c3, #d8, #df, #a9, #e0, #aa
db #d5, #cb, #d7, #df, #aa, #e0, #e2, #6a
db #d7, #d3, #df, #e0, #aa, #e2, #e3, #a9
db #e4, #aa, #d5, #db, #df, #e0, #33, #10
db #18, #2d, #4d, #65, #82, #01, #c9, #5d
db #c1, #59, #96, #ac, #f9, #f7, #c1, #00
db #ca, #fc, #9a, #00, #da, #f9, #c0, #4f
db #c8, #a3, #e7, #28, #ca, #97, #d2, #9f
db #7c, #e7, #ca, #40, #f3, #40, #63, #8a
db #11, #cb, #94, #4a, #87, #eb, #81, #fb
db #4b, #04, #50, #e0, #6b, #e8, #5b, #d1
db #8c, #db, #88, #e3, #40, #08, #00, #66
db #f0, #bd, #f8, #0c, #e9, #57, #d9, #9e
db #da, #f0, #fa, #87, #fb, #8c, #fb, #d0
db #c5, #62, #c0, #d8, #4d, #ca, #9d, #ca
db #3b, #cb, #2d, #04, #a0, #bc, #e0, #0b
db #e1, #48, #da, #8b, #d3, #8a, #db, #8c
db #e3, #8a, #eb, #89, #f3, #c3, #07, #da
db #bd, #e0, #5d, #c9, #96, #07, #d1, #e6
db #da, #87, #d3, #88, #db, #68, #03, #97
db #b9, #e8, #af, #c1, #fb, #e1, #21, #06
db #84, #cd, #ac, #e8, #09, #d9, #ea, #fa
db #84, #c3, #83, #d3, #89, #52, #60, #6f
db #ad, #be, #f8, #70, #13, #0c, #61, #e7
db #aa, #d1, #f7, #d9, #00, #f2, #da, #a1
db #ea, #0e, #f2, #ed, #c2, #ed, #d2, #40
db #cb, #8d, #e3, #93, #e3, #8d, #eb, #92
db #eb, #88, #fb, #8d, #fb, #91, #8d, #0e
db #c3, #e1, #d3, #34, #75, #02, #0c, #c9
db #4c, #f2, #4f, #f2, #48, #05, #9a, #0d
db #c9, #0c, #1a, #fd, #f9, #89, #cb, #88
db #c8, #b0, #eb, #4a, #10, #c9, #9c, #fa
db #3d, #c3, #85, #eb, #0e, #0a, #57, #0b
db #d1, #5a, #f1, #a9, #e1, #a9, #e9, #aa
db #f1, #ae, #f1, #48, #e2, #37, #58, #1d
db #2f, #f3, #8b, #fb, #0c, #91, #02, #ef
db #38, #c1, #a9, #cf, #29, #f9, #ae, #f9
db #fd, #a0, #98, #d2, #e9, #e0, #3c, #cb
db #a4, #02, #d7, #e1, #e1, #a0, #ea, #0b
db #f3, #e5, #0b, #e9, #e7, #da, #08, #8d
db #aa, #0c, #f1, #aa, #c9, #f9, #c1, #fd
db #f8, #28, #f8, #f9, #1e, #c2, #48, #ca
db #a0, #f2, #e9, #da, #16, #03, #c5, #5e
db #f9, #e8, #fa, #d4, #c3, #2c, #04, #7b
db #29, #f9, #5c, #43, #98, #ca, #8c, #f3
db #c0, #1c, #00, #60, #c1, #60, #c9, #af
db #c9, #fe, #c1, #fb, #d1, #fb, #d9, #f9
db #e9, #4e, #f2, #4e, #fa, #9e, #d2, #a0
db #a8, #fd, #da, #6a, #e2, #9e, #ea, #ee
db #ca, #d2, #aa, #da, #e2, #aa, #ea, #f2
db #84, #a6, #fa, #3e, #c3, #38, #cb, #3c
db #d3, #41, #db, #85, #92, #e3, #94, #b9
db #84, #19, #a2, #af, #d1, #e9, #04, #01
db #ab, #aa, #41, #82, #f7, #a8, #f9, #30
db #1f, #85, #e4, #f7, #c9, #00, #e2, #f9
db #f1, #49, #c2, #49, #ca, #4f, #fd, #a2
db #d2, #4e, #e2, #47, #fa, #4f, #9f, #c2
db #a8, #ca, #89, #da, #ec, #29, #e2, #3a
db #cb, #3d, #d3, #40, #a6, #38, #db, #40
db #38, #e3, #40, #8a, #38, #eb, #b1, #f3
db #92, #df, #db, #8d, #f3, #86, #fb, #d9
db #aa, #d9, #d9, #b6, #96, #ed, #42, #f8
db #d9, #eb, #e2, #e8, #ea, #e8, #f2, #38
db #c3, #c2, #02, #91, #bb, #4d, #c2, #85
db #e3, #50, #03, #00, #4e, #ea, #9d, #d2
db #d5, #c3, #e1, #4c, #26, #89, #db, #90
db #06, #97, #97, #c2, #9d, #ea, #9c, #f2
db #9d, #d6, #fd, #87, #32, #e3, #1c, #6b
db #4b, #9c, #c2, #e7, #f2, #e7, #fa, #03
db #19, #27, #f9, #9d, #c2, #12, #9a, #a3
db #9f, #d2, #e7, #72, #40, #eb, #18, #d3
db #6e, #9b, #da, #83, #c5, #39, #ea, #a0
db #ed, #8a, #9e, #fa, #25, #eb, #bb, #86
db #c7, #e9, #70, #f1, #8e, #a5, #fb, #b4
db #d1, #62, #ec, #f2, #92, #a8, #ec, #06
db #9e, #39, #c3, #0d, #4b, #d3, #3c, #d5
db #4e, #ad, #e9, #a2, #40, #60, #05, #9a
db #39, #cb, #87, #5e, #85, #f3, #85, #fb
db #e1, #cb, #38, #dd, #7a, #3d, #cb, #58
db #f7, #87, #a2, #4a, #62, #84, #d3, #14
db #a9, #86, #49, #9a, #db, #29, #02, #0c
db #89, #e3, #89, #eb, #00, #02
delta07:
db #15, #7a, #34, #80, #09, #00, #16, #c0
db #16, #c8, #16, #e0, #fb, #d1, #4d, #e2
db #4d, #f2, #3d, #c3, #41, #db, #86, #e3
db #87, #06, #00, #17, #c0, #bd, #e0, #4c
db #f2, #4b, #fa, #38, #c3, #87, #d3, #78
db #10, #00, #1e, #c0, #1e, #c8, #6e, #c0
db #6e, #f9, #d0, #aa, #d8, #e0, #aa, #e8
db #f0, #86, #f8, #be, #c8, #be, #d8, #a9
db #e0, #5e, #e8, #e9, #ca, #ec, #ea, #f0
db #0c, #00, #1f, #c0, #19, #d0, #1a, #d0
db #1e, #fd, #aa, #d8, #e0, #aa, #e8, #f0
db #05, #a6, #f8, #5f, #e1, #00, #fa, #9a
db #da, #92, #01, #00, #17, #c8, #1e, #07
db #00, #1c, #c8, #5e, #f1, #f9, #ae, #c1
db #a8, #c9, #4a, #e1, #3a, #fb, #07, #03
db #00, #1b, #d0, #e8, #f2, #56, #fa, #a0
db #02, #00, #16, #d8, #85, #eb, #3c, #0d
db #00, #1a, #e0, #be, #c0, #86, #d0, #0e
db #f9, #5e, #c1, #aa, #c9, #d1, #aa, #d9
db #e1, #84, #68, #e9, #ae, #d9, #f8, #e9
db #e8, #da, #20, #03, #00, #16, #e8, #0c
db #c1, #4d, #ca, #e0, #0d, #00, #20, #e8
db #1a, #f8, #af, #c1, #fa, #e9, #fc, #f9
db #4c, #ca, #a5, #d2, #b9, #da, #86, #c3
db #94, #cb, #dd, #c3, #f9, #e0, #db, #61
db #02, #2f, #5b, #f0, #40, #cb, #21, #01
db #00, #67, #c0, #01, #04, #00, #67, #d0
db #fc, #d1, #9c, #f2, #9f, #f2, #86, #d3
db #bd, #e8, #ea, #fa, #16, #20, #a4, #be
db #f8, #0e, #03, #4e, #0b, #c1, #aa, #f9
db #84, #db, #30, #08, #00, #0e, #d1, #0e
db #d9, #00, #ea, #47, #f2, #8a, #cb, #94
db #db, #88, #f3, #d9, #c3, #b0, #b1, #6b
db #10, #d1, #47, #34, #80, #bb, #89, #0e
db #e1, #0e, #e9, #86, #d3, #b4, #01, #a1
db #f1, #0c, #07, #43, #0c, #f9, #fd, #f9
db #9d, #c2, #9d, #ca, #e7, #f2, #e7, #fa
db #8c, #f3, #00, #40, #00, #fb, #c1, #fe
db #c1, #01, #fa, #4f, #c2, #4d, #d2, #4d
db #da, #49, #ea, #4d, #ea, #4e, #f2, #4e
db #fa, #9e, #c2, #9e, #ca, #97, #d2, #9e
db #aa, #fd, #da, #e2, #aa, #ea, #f2, #82
db #fa, #ec, #c2, #ed, #a9, #ee, #a8, #ec
db #ca, #ed, #ee, #82, #d2, #ec, #da, #ee
db #98, #ed, #e2, #ee, #a8, #ea, #28, #f2
db #ed, #fa, #ee, #a0, #3e, #c3, #a6, #cb
db #3d, #d3, #3e, #3d, #db, #3e, #2a, #e3
db #aa, #eb, #f3, #6a, #fb, #8e, #c3, #cb
db #a8, #d3, #a6, #db, #8f, #87, #e3, #8e
db #9a, #86, #eb, #87, #88, #8e, #2a, #f3
db #81, #fb, #86, #87, #a8, #8e, #a0, #de
db #c3, #89, #cb, #d5, #d3, #de, #0a, #db
db #10, #09, #00, #00, #d2, #49, #c2, #4e
db #ea, #9f, #e2, #02, #ea, #9c, #fa, #ea
db #f2, #38, #db, #d6, #b2, #70, #35, #18
db #f7, #e9, #4f, #d2, #c1, #9e, #fc, #f1
db #42, #e9, #18, #4d, #c2, #85, #fb, #28
db #98, #4c, #e2, #2c, #92, #4c, #ea, #8c
db #eb, #4b, #03, #24, #e7, #fa, #89, #e3
db #d4, #db, #48, #d7, #4d, #fa, #2d, #82
db #f7, #9c, #c2, #0f, #04, #20, #39, #d2
db #9c, #da, #89, #f3, #d4, #d3, #94, #ad
db #8b, #9d, #d2, #58, #98, #8c, #e5, #cb
db #90, #82, #13, #9b, #da, #08, #05, #82
db #9d, #da, #9d, #e2, #ac, #ea, #63, #dd
db #cb, #84, #86, #cf, #a0, #ea, #e9, #e2
db #e1, #27, #e8, #ca, #40, #a9, #09, #ed
db #a0, #ed, #da, #ed, #6a, #8d, #e3, #8d
db #eb, #8d, #f3, #8d, #fb, #d8, #e5, #ba
db #1c, #bf, #e7, #98, #3c, #cb, #43, #82
db #e6, #f2, #c0, #0c, #20, #ea, #f2, #37
db #e3, #3d, #fd, #eb, #a6, #f3, #fb, #8d
db #c3, #aa, #cb, #d3, #93, #db, #93, #f3
db #dd, #d3, #1a, #aa, #0b, #ec, #86, #27
db #39, #c3, #d0, #98, #ed, #36, #fb, #83
db #98, #86, #cb, #02, #98, #87, #db, #0b
db #98, #89, #eb, #06, #99, #85, #f3, #69
db #8e, #fb, #04, #bb, #b0, #d5, #00, #08
delta08:
db #57, #3e, #e0, #13, #00, #16, #c0, #16
db #c8, #19, #d0, #20, #d8, #70, #d0, #70
db #83, #f9, #e0, #0d, #c9, #10, #90, #fd
db #a6, #d9, #af, #c9, #f8, #e1, #f7, #f1
db #4b, #f2, #4a, #fa, #4b, #87, #db, #93
db #80, #8a, #92, #fb, #0f, #6b, #00, #17
db #c0, #1b, #c8, #aa, #d0, #d8, #aa, #e0
db #e8, #a9, #f0, #29, #f8, #6b, #c0, #68
db #c8, #6b, #88, #68, #d0, #6b, #02, #d8
db #bd, #e0, #bc, #e8, #b9, #f8, #09, #c1
db #0a, #f9, #5b, #e9, #aa, #d1, #ab, #22
db #d9, #ac, #98, #ab, #e1, #ac, #86, #e9
db #fd, #d9, #f9, #e1, #a8, #e9, #a6, #f1
db #fa, #f9, #f9, #fa, #a9, #fb, #a8, #49
db #c2, #4a, #4b, #80, #a6, #ca, #4c, #da
db #9d, #ca, #9a, #d2, #9b, #99, #da, #9a
db #a9, #9b, #aa, #99, #e2, #9a, #9b, #6a
db #99, #ea, #9b, #9c, #20, #a1, #f2, #98
db #fa, #9c, #a0, #e8, #c2, #e7, #ca, #e8
db #a1, #ed, #da, #ea, #f2, #a2, #fa, #3a
db #c3, #3b, #cb, #d3, #3c, #9a, #3a, #db
db #3b, #a6, #3c, #3b, #e3, #3c, #9a, #3a
db #eb, #3b, #a6, #3c, #3a, #f3, #3b, #9a
db #3a, #fb, #3b, #a6, #3c, #8b, #c3, #8c
db #9a, #8a, #cb, #8b, #a6, #8d, #8b, #d3
db #8c, #a8, #8d, #a2, #8c, #db, #e3, #8d
db #9a, #8a, #eb, #8c, #a1, #8d, #aa, #8c
db #f3, #89, #fb, #8b, #8d, #6a, #d4, #c3
db #da, #db, #86, #d4, #cb, #db, #d3, #dd
db #9a, #d4, #db, #da, #aa, #db, #dc, #06
db #1e, #12, #00, #18, #c0, #1c, #d0, #1d
db #d8, #aa, #e0, #e8, #a9, #f0, #a4, #f8
db #6d, #c0, #0a, #c8, #6b, #f8, #ad, #f9
db #fd, #c9, #fa, #e9, #48, #da, #ed, #d2
db #e8, #da, #3d, #fb, #dd, #db, #0e, #15
db #00, #1c, #c0, #1a, #c8, #1a, #f0, #0b
db #e9, #6a, #f1, #5c, #d9, #e1, #90, #22
db #e9, #a9, #f1, #fd, #e1, #4c, #d2, #4d
db #ea, #9c, #e2, #98, #f2, #ed, #ca, #3c
db #cb, #3d, #f3, #8d, #c3, #8c, #cb, #da
db #12, #d3, #10, #14, #00, #1d, #c0, #67
db #c8, #06, #c1, #86, #c9, #0e, #d1, #06
db #d9, #81, #a6, #e1, #a9, #c1, #ae, #e1
db #fc, #e9, #00, #f2, #fa, #50, #c2, #81
db #28, #ca, #49, #fa, #e6, #ca, #38, #d3
db #85, #cb, #88, #f3, #e1, #db, #f0, #1a
db #00, #1e, #c0, #68, #c8, #19, #d8, #20
db #e8, #28, #f8, #6e, #c0, #70, #aa, #6e
db #c8, #d0, #aa, #d8, #e0, #aa, #e8, #f0
db #69, #f8, #be, #c0, #0a, #c8, #c0, #d0
db #be, #d8, #0e, #f1, #fa, #d9, #9a, #c2
db #9b, #29, #37, #eb, #4a, #f3, #38, #fb
db #93, #e3, #07, #18, #00, #17, #c8, #6b
db #e8, #68, #f0, #69, #f8, #09, #c9, #69
db #d1, #5a, #c9, #5b, #e1, #5d, #e9, #4a
db #f1, #a9, #e1, #4c, #e2, #4d, #fa, #9d
db #d2, #ec, #c2, #ea, #e2, #ed, #4a, #38
db #c3, #39, #cb, #35, #db, #1a, #e3, #3c
db #f3, #dd, #c3, #06, #cb, #80, #1f, #00
db #1a, #c8, #16, #d8, #1a, #80, #0a, #16
db #e8, #71, #c8, #0d, #d1, #0c, #e9, #07
db #f1, #57, #c1, #60, #d1, #5a, #f9, #ac
db #c9, #ad, #d9, #fe, #c1, #fa, #c9, #fc
db #d9, #f8, #e9, #fc, #f1, #49, #d2, #aa
db #da, #e2, #5a, #ea, #4e, #f2, #97, #c2
db #9e, #ca, #05, #a0, #da, #35, #fb, #87
db #c3, #8e, #db, #93, #f3, #e5, #db, #12
db #04, #00, #1d, #c8, #0e, #d9, #ae, #28
db #91, #f3, #c0, #10, #00, #16, #d0, #1a
db #a0, #16, #f0, #69, #f8, #66, #c0, #aa
db #c1, #ab, #c9, #af, #d1, #ab, #d9, #fe
db #c9, #d1, #20, #b5, #5e, #ca, #42, #f3
db #88, #c3, #94, #d3, #87, #11, #00, #17
db #d0, #19, #e0, #bd, #e8, #0b, #f9, #5a
db #bf, #0c, #e9, #f8, #f9, #99, #30, #2b
db #80, #6f, #88, #ea, #e9, #fa, #39, #c3
db #3a, #cb, #39, #aa, #d3, #37, #8d, #39
db #b4, #02, #5b, #17, #d0, #5e, #c1, #3c
db #0a, #00, #18, #d8, #18, #e0, #5e, #f9
db #ae, #c1, #ae, #c9, #fb, #f1, #48, #d2
db #48, #e2, #e7, #da, #36, #db, #40, #9e
db #d3, #16, #e0, #b6, #29, #13, #f1, #10
db #f9, #5b, #f1, #9b, #ca, #e9, #da, #38
db #86, #4d, #fb, #d5, #c3, #96, #03, #25
db #72, #e0, #fd, #c1, #99, #ca, #08, #0d
db #00, #1a, #e0, #af, #e1, #4d, #e2, #9e
db #ad, #a1, #d2, #3a, #e2, #9d, #fa, #e9
db #ca, #88, #d3, #86, #e3, #8e, #fd, #eb
db #9e, #f3, #d2, #01, #b5, #78, #e8, #0d
db #05, #4f, #5f, #e8, #0c, #d1, #8a, #db
db #d9, #c3, #d9, #db, #0c, #16, #a3, #8c
db #e7, #d9, #0b, #20, #8b, #a0, #c1, #5c
db #c9, #ad, #b0, #ad, #e9, #4c, #c2, #c7
db #41, #37, #84, #29, #c6, #ea, #3c, #c3
db #3d, #e3, #87, #cb, #dc, #c3, #de, #c3
db #2a, #cb, #de, #d9, #aa, #de, #de, #ba
db #41, #32, #5e, #f0, #e6, #d2, #39, #fb
db #43, #05, #ef, #15, #d9, #f8, #ac, #f1
db #4a, #ca, #4b, #d2, #86, #cb, #03, #19
db #00, #67, #c0, #b8, #c8, #b8, #d0, #09
db #d9, #09, #e1, #5a, #0f, #1a, #e1, #5d
db #f9, #ad, #c1, #aa, #d9, #fc, #c1, #4e
db #ca, #4c, #ea, #a5, #da, #a0, #ca, #39
db #db, #39, #e3, #39, #eb, #ae, #f3, #89
db #d3, #82, #85, #fd, #0a, #eb, #d1, #d3
db #dc, #28, #00, #83, #aa, #71, #c0, #d8
db #aa, #e0, #e8, #a6, #f0, #f8, #c1, #c0
db #aa, #c8, #d0, #aa, #d8, #e0, #a8, #e8
db #06, #f0, #bd, #f8, #0c, #c1, #07, #e9
db #5d, #c1, #60, #e9, #5a, #f1, #ac, #c1
db #ad, #c9, #92, #d1, #aa, #e9, #a8, #f1
db #aa, #a9, #ae, #aa, #aa, #f9, #ae, #b1
db #4a, #fa, #c1, #fb, #c9, #00, #ca, #42
db #d2, #f8, #d9, #01, #e2, #fd, #f1, #fc
db #f9, #4d, #ca, #4f, #9a, #4a, #d2, #4f
db #a6, #50, #4a, #da, #4b, #aa, #4f, #50
db #9a, #4a, #e2, #4b, #aa, #4e, #4f, #a6
db #50, #4a, #ea, #4b, #aa, #4e, #4f, #a9
db #50, #a6, #4f, #f2, #50, #48, #fa, #4f
db #a9, #50, #aa, #98, #c2, #9f, #a0, #6a
db #97, #ca, #98, #9f, #a9, #a0, #aa, #98
db #d2, #9f, #a0, #6a, #98, #da, #9f, #a0
db #9a, #9d, #e2, #9f, #a6, #a0, #9d, #ea
db #9f, #a9, #a0, #aa, #99, #f2, #9a, #9d
db #aa, #9f, #a0, #6a, #99, #fa, #9a, #9b
db #9a, #e9, #c2, #ea, #a6, #eb, #ea, #ca
db #eb, #9a, #e9, #d2, #ea, #89, #eb, #a4
db #da, #e9, #ea, #a6, #f2, #f0, #fa, #3d
db #c3, #40, #3d, #cb, #40, #22, #d3, #41
db #9a, #38, #db, #40, #a6, #41, #38, #e3
db #40, #a9, #41, #aa, #38, #eb, #40, #42
db #62, #35, #f3, #40, #0a, #fb, #84, #c3
db #90, #69, #88, #cb, #90, #a6, #83, #d3
db #90, #82, #db, #85, #92, #94, #eb, #e2
db #c3, #e1, #cb, #80, #aa, #d3, #d8, #db
db #60, #09, #00, #66, #c8, #e8, #1a, #f8
db #ab, #c1, #f7, #d9, #28, #e9, #fb, #1a
db #fe, #f1, #ea, #da, #20, #0d, #00, #66
db #d8, #6a, #e0, #b6, #d0, #e0, #a8, #e8
db #29, #f0, #5b, #f9, #5c, #a1, #aa, #c9
db #f9, #4a, #49, #f2, #85, #c3, #91, #eb
db #4b, #04, #00, #6b, #f0, #3a, #e3, #88
db #db, #d4, #d3, #78, #07, #00, #be, #d0
db #a1, #f0, #aa, #f8, #0e, #f9, #5e, #d1
db #d9, #4a, #f1, #01, #12, #00, #b8, #e0
db #5a, #e8, #bd, #f0, #0e, #c1, #09, #e9
db #0a, #f1, #5d, #c9, #a9, #41, #a1, #aa
db #e1, #ae, #e9, #fc, #c9, #4c, #f2, #97
db #fa, #ec, #d2, #ed, #f2, #89, #c3, #81
db #fb, #e1, #c3, #d0, #05, #00, #c0, #e0
db #18, #e8, #fa, #d1, #37, #fb, #94, #cb
db #30, #0f, #00, #b6, #f8, #0e, #e1, #4e
db #c2, #4a, #f2, #36, #eb, #41, #82, #f3
db #36, #fb, #41, #8a, #91, #c3, #aa, #cb
db #d3, #aa, #db, #e3, #38, #fb, #a4, #1f
db #71, #0d, #c1, #49, #ca, #ec, #cd, #27
db #7d, #24, #c5, #ad, #0c, #c9, #ec, #fa
db #a0, #03, #10, #71, #a3, #9f, #33, #da
db #2c, #09, #e1, #23, #c7, #d1, #ad, #49
db #bf, #5a, #ca, #4c, #ca, #48, #f2, #e7
db #d2, #ec, #e2, #3d, #eb, #02, #02, #00
db #5c, #ae, #fd, #48, #37, #17, #ac, #d1
db #4e, #fa, #9c, #d2, #3d, #db, #8e, #fb
db #16, #81, #79, #ed, #ae, #d1, #e8, #d2
db #e8, #e2, #ec, #ea, #21, #87, #a8, #63
db #3d, #f9, #89, #7a, #e7, #67, #f8, #82
db #ed, #ea, #84, #01, #ab, #f8, #b0, #87
db #e9, #f9, #d1, #35, #eb, #1c, #9a, #db
db #fd, #d1, #eb, #b8, #94, #e5, #c6, #d9
db #50, #25, #58, #fb, #d9, #e9, #e2, #85
db #d3, #70, #9e, #f7, #e1, #fb, #25, #4e
db #d2, #36, #e3, #94, #db, #68, #bd, #18
db #fa, #e1, #fe, #e9, #4a, #9e, #fe, #e1
db #06, #6d, #1a, #fd, #e9, #e8, #ea, #dc
db #98, #81, #39, #f1, #4c, #fa, #61, #f3
db #29, #fe, #f9, #ec, #f2, #0b, #04, #a1
db #48, #c2, #4d, #e6, #9d, #c2, #98, #e2
db #04, #49, #4d, #c2, #c1, #2e, #9c, #e9
db #38, #ea, #a1, #c5, #62, #9c, #ca, #8c
db #fb, #c3, #06, #9a, #ea, #37, #e3, #87
db #d3, #89, #db, #a5, #2a, #e7, #e6, #e1
db #f7, #eb, #e2, #09, #26, #ed, #fa, #c2
db #26, #84, #d3, #14, #2a, #86, #e1, #2d
db #4f, #ae, #8a, #e3, #8d, #f3, #d9, #86
db #da, #b0, #f3, #00, #08
delta09:
db #15, #e5, #3e, #80, #1c, #00, #16, #c0
db #71, #d8, #71, #e0, #60, #c1, #5b, #c9
db #60, #c9, #5d, #d1, #5e, #e9, #aa, #c1
db #af, #fd, #a0, #c9, #ac, #e1, #ad, #f9
db #f9, #d1, #a6, #d9, #4c, #ca, #4f, #4c
db #d2, #4f, #2a, #da, #41, #a4, #e2, #4d
db #ea, #49, #f2, #9a, #da, #42, #f3, #85
db #c3, #94, #eb, #e2, #c3, #0f, #35, #00
db #18, #c0, #17, #c8, #68, #e8, #6b, #a0
db #bc, #e0, #09, #c9, #0b, #e9, #22, #f1
db #5a, #c1, #5c, #d9, #5b, #e1, #5c, #69
db #e9, #ac, #c1, #a2, #c9, #aa, #e1, #e9
db #ab, #9a, #aa, #f1, #ab, #6a, #aa, #f9
db #ab, #ac, #8a, #fc, #c1, #a4, #c9, #88
db #d1, #fd, #e1, #49, #ca, #4a, #a6, #d2
db #4b, #4a, #da, #4b, #8a, #4c, #e2, #6a
db #ea, #9d, #d2, #ea, #81, #a0, #f2, #9b
db #fa, #eb, #c2, #e8, #da, #ec, #ea, #a2
db #f2, #3a, #cb, #3c, #3d, #e3, #81, #28
db #eb, #3c, #f3, #3d, #fb, #8c, #cb, #8a
db #db, #da, #cb, #d4, #d3, #2c, #09, #00
db #1c, #c0, #4a, #c8, #98, #ca, #ee, #c2
db #ec, #d2, #e7, #da, #ec, #29, #8e, #cb
db #a0, #d3, #00, #55, #69, #1d, #c0, #1a
db #d0, #66, #f0, #06, #c1, #0d, #69, #06
db #d9, #07, #f1, #0b, #f9, #0c, #88, #5b
db #c1, #5c, #a6, #c9, #5d, #57, #d9, #5e
db #22, #e1, #60, #9a, #59, #f9, #5f, #2a
db #a9, #c1, #c9, #a8, #d1, #a0, #d9, #ac
db #69, #ad, #f1, #f7, #c1, #01, #c2, #fa
db #c9, #fe, #aa, #fa, #d1, #fb, #fe, #69
db #fa, #d9, #fb, #a6, #00, #da, #01, #f9
db #e1, #fa, #aa, #fb, #fc, #80, #62, #00
db #e2, #f9, #e9, #00, #ea, #f9, #f1, #00
db #f2, #f9, #f9, #00, #fa, #4c, #c2, #50
db #68, #ca, #4d, #da, #29, #e2, #4a, #f2
db #4b, #aa, #49, #fa, #4a, #4b, #a6, #4e
db #9a, #c2, #9b, #a9, #9e, #aa, #99, #ca
db #9a, #9b, #a6, #9e, #99, #d2, #9a, #8a
db #99, #da, #a0, #e2, #22, #ea, #9c, #fa
db #ec, #c2, #e9, #da, #ea, #29, #e2, #ed
db #a1, #ea, #ea, #ed, #5a, #41, #eb, #43
db #fb, #94, #c3, #85, #d3, #91, #eb, #88
db #fb, #e5, #db, #d0, #07, #00, #20, #c0
db #0c, #d1, #ae, #e1, #a1, #e9, #2a, #f1
db #9b, #da, #87, #c3, #c0, #0a, #00, #16
db #c8, #d8, #a9, #e0, #68, #e8, #10, #e9
db #ab, #d1, #af, #e1, #46, #e9, #9b, #d2
db #93, #f3, #a0, #04, #00, #16, #d0, #10
db #d9, #4a, #c2, #38, #fb, #07, #10, #00
db #17, #d0, #a8, #e0, #44, #aa, #e8, #b8
db #c0, #09, #d9, #5b, #f1, #aa, #d9, #a9
db #e9, #fb, #c1, #fd, #e9, #49, #d2, #9a
db #f2, #e8, #d2, #eb, #ea, #83, #db, #88
db #eb, #f0, #15, #00, #19, #d0, #20, #d8
db #70, #c8, #c0, #d8, #5e, #f9, #ae, #c1
db #c9, #a6, #d1, #d9, #ff, #c1, #82, #c9
db #f7, #f1, #fb, #21, #a5, #f9, #49, #c2
db #47, #d2, #a0, #e2, #36, #e3, #37, #fb
db #92, #eb, #46, #fb, #0e, #1d, #00, #1c
db #d0, #1a, #f8, #0b, #d9, #5d, #e1, #5a
db #f9, #ac, #f1, #fd, #c9, #48, #c2, #4c
db #da, #9d, #c2, #ec, #e2, #e9, #fa, #39
db #c3, #3d, #cb, #8e, #db, #8d, #e3, #8e
db #9a, #8d, #eb, #8e, #6a, #84, #f3, #8a
db #8e, #26, #fb, #de, #c3, #82, #cb, #dd
db #d3, #de, #9a, #dd, #db, #de, #16, #87
db #0c, #00, #17, #d8, #6b, #f0, #5c, #f1
db #f8, #d1, #94, #a1, #d9, #4b, #e2, #9a
db #ea, #e8, #e2, #39, #cb, #3d, #db, #d4
db #c3, #db, #68, #1e, #0e, #00, #18, #d8
db #6d, #d0, #0b, #c1, #5d, #e9, #c8, #f1
db #ad, #bd, #2b, #d1, #48, #ca, #9b, #3c
db #88, #5d, #d3, #8d, #83, #f3, #8d, #fb
db #96, #26, #5d, #1d, #d8, #ad, #87, #38
db #cb, #87, #d3, #c3, #84, #eb, #aa, #19
db #e0, #17, #f0, #f8, #c9, #86, #cb, #40
db #0b, #1a, #5e, #66, #c8, #66, #d0, #66
db #d8, #66, #f1, #a3, #e8, #f8, #57, #66
db #cd, #f9, #fb, #e9, #8e, #e2, #3c, #cf
db #a1, #1d, #4f, #18, #e8, #1d, #e8, #0d
db #e1, #0d, #e9, #0d, #f1, #5c, #f9, #4e
db #da, #48, #f2, #48, #fa, #98, #c2, #68
db #03, #9d, #87, #d7, #d1, #9e, #ea, #0c
db #10, #90, #ef, #e0, #f0, #5d, #d9, #fd
db #c1, #f8, #e9, #f8, #f1, #f8, #f9, #4d
db #c9, #e1, #d2, #98, #da, #ee, #fd, #aa
db #e2, #3d, #c3, #3e, #eb, #f3, #05, #29
db #fb, #8e, #c3, #43, #06, #00, #67, #c0
db #5c, #d1, #5b, #f9, #ed, #c2, #39, #d3
db #82, #eb, #03, #0a, #00, #67, #c8, #b8
db #d8, #ac, #d1, #fd, #f1, #4e, #ea, #4c
db #f2, #9d, #fa, #eb, #ca, #ed, #02, #ea
db #fa, #10, #15, #00, #67, #d0, #06, #d1
db #5d, #c1, #5f, #f1, #a9, #e1, #ff, #e9
db #4e, #f2, #97, #ca, #eb, #da, #ed, #86
db #41, #f3, #35, #fb, #41, #8a, #91, #c3
db #aa, #cb, #d3, #aa, #db, #e3, #b1, #f3
db #8e, #c9, #c3, #01, #a6, #93, #d8, #b8
db #f0, #0c, #c9, #a8, #34, #b3, #83, #71
db #d9, #fc, #e9, #4d, #4e, #79, #69, #0a
db #fa, #9c, #c2, #9c, #1e, #9c, #d2, #9c
db #da, #e6, #8b, #b0, #d2, #91, #d5, #aa
db #e2, #ea, #f2, #85, #db, #82, #e3, #78
db #03, #00, #6e, #e0, #be, #be, #ca, #4b
db #97, #6c, #92, #a9, #f1, #89, #db, #2d
db #07, #9c, #6b, #f8, #bd, #d2, #81, #fb
db #0d, #d1, #ec, #fa, #3d, #8c, #3f, #eb
db #60, #a9, #bb, #b6, #0a, #b6, #c8, #fa
db #e9, #fe, #e9, #9e, #fa, #85, #fb, #69
db #01, #0e, #bc, #e8, #20, #05, #d9, #34
db #f8, #5b, #d6, #d1, #ad, #0d, #f9, #9c
db #f2, #30, #90, #ef, #2b, #0e, #c1, #5e
db #c9, #5f, #e1, #aa, #c9, #ad, #e1, #ad
db #e9, #f7, #d1, #ff, #e1, #97, #c2, #38
db #d3, #85, #cb, #24, #02, #0d, #48, #7d
db #e7, #f9, #1c, #9d, #0b, #e1, #08, #95
db #53, #89, #0c, #e1, #f9, #c9, #ee, #ca
db #e9, #d2, #ed, #f2, #ed, #fa, #3e, #e3
db #21, #4f, #09, #e9, #ab, #c1, #4d, #ca
db #9a, #fa, #89, #c3, #d1, #cb, #a4, #02
db #c1, #9e, #41, #f2, #90, #04, #f3, #03
db #f1, #4b, #c2, #9b, #ea, #36, #f3, #34
db #25, #cb, #b9, #0e, #f1, #0e, #f9, #5e
db #c1, #16, #73, #38, #f9, #e8, #ea, #dc
db #cb, #06, #b9, #62, #5a, #c9, #84, #fb
db #12, #06, #5e, #d1, #5f, #e9, #4d, #c2
db #87, #db, #0b, #2c, #5b, #83, #af, #c9
db #9d, #da, #86, #a8, #91, #5a, #0a, #f8
db #e1, #e9, #f2, #dd, #c3, #e0, #07, #a4
db #5e, #79, #ab, #c9, #af, #f1, #af, #f9
db #fe, #f1, #93, #c3, #94, #cb, #b4, #93
db #2a, #5d, #f9, #4e, #d2, #70, #09, #aa
db #9c, #ff, #d1, #f7, #60, #13, #ac, #d9
db #f7, #e9, #fa, #fa, #22, #a1, #ca, #47
db #39, #d7, #69, #ab, #d9, #ac, #e9, #f9
db #c1, #fe, #69, #4d, #f2, #9e, #e2, #98
db #ea, #ec, #c9, #ee, #d2, #28, #e9, #2a
db #fe, #e1, #4b, #ca, #81, #01, #fc, #b8
db #18, #e9, #62, #fd, #f9, #9a, #e2, #a5
db #62, #4e, #c2, #a1, #6a, #49, #da, #9c
db #14, #66, #e2, #84, #03, #aa, #9e, #98
db #9c, #e9, #ea, #0d, #65, #ab, #ae, #9d
db #e2, #e7, #d2, #8a, #e3, #d9, #c1, #b7
db #62, #9c, #ea, #e1, #4a, #e6, #d2, #37
db #e3, #04, #03, #a1, #e9, #e7, #ee, #ea
db #d5, #c3, #50, #d7, #36, #fb, #61, #ab
db #f7, #39, #02, #96, #ad, #84, #d3, #85
db #e3, #85, #eb, #dd, #a2, #58, #79, #86
db #db, #83, #cd, #8a, #88, #db, #29, #8c
db #c0, #00, #20
delta10:
db #11, #4e, #24, #00, #06, #00, #16, #c0
db #1a, #e0, #f9, #d9, #49, #f2, #9a, #da
db #94, #eb, #70, #03, #00, #19, #c0, #37
db #fb, #93, #e3, #78, #01, #00, #1e, #c0
db #e0, #05, #00, #20, #c8, #10, #c1, #9e
db #f2, #37, #eb, #92, #fb, #80, #08, #00
db #16, #d0, #16, #e0, #71, #e8, #9b, #e2
db #9b, #ea, #e9, #da, #38, #fb, #95, #d3
db #0e, #02, #00, #18, #d0, #9b, #f2, #d2
db #01, #00, #1d, #d8, #1e, #04, #00, #1d
db #e0, #1d, #e8, #5c, #f9, #e7, #ca, #40
db #d3, #18, #16, #e8, #66, #f0, #28, #9e
db #1a, #e8, #3c, #5f, #a1, #6d, #e5, #6d
db #e8, #48, #da, #c0, #ef, #e6, #66, #c8
db #60, #c1, #98, #ea, #12, #a5, #0e, #c1
db #2c, #20, #79, #0d, #d1, #0c, #e9, #98
db #c2, #98, #d2, #16, #e1, #ea, #5a, #c9
db #52, #f7, #5e, #e4, #30, #91, #a5, #5b
db #d1, #5e, #d1, #10, #07, #7a, #aa, #c9
db #a9, #d9, #f7, #c1, #4b, #c2, #97, #c2
db #85, #c3, #85, #db, #a0, #bd, #ab, #a2
db #21, #62, #f8, #c1, #c3, #a7, #f8, #89
db #db, #f0, #86, #3f, #f7, #e1, #f7, #e9
db #02, #27, #4d, #c2, #e1, #98, #f7, #47
db #d2, #0c, #98, #98, #ca, #20, #98, #99
db #d2, #4a, #9e, #9e, #ea, #43, #c1, #18
db #e6, #fa, #86, #cb, #b0, #9e, #36, #e3
db #4b, #87, #62, #3a, #eb, #60, #62, #38
db #f3, #50, #6b, #87, #c3, #94, #01, #9e
db #e9, #89, #c3, #87, #cb, #62, #87, #d3
db #0d, #a8, #8a, #2d, #8c, #d9, #db, #00
db #02
delta11:
db #05, #0e, #3a, #81, #01, #00, #16, #c0
db #f0, #0a, #00, #19, #c0, #20, #c8, #1d
db #d0, #be, #d0, #be, #e0, #c0, #e8, #36
db #e3, #94, #db, #93, #e3, #db, #c3, #e0
db #09, #00, #1a, #c0, #20, #e8, #70, #c8
db #10, #d9, #ab, #c9, #fb, #f1, #4a, #c2
db #9b, #e2, #42, #fb, #80, #08, #00, #16
db #c8, #1a, #d0, #66, #c8, #66, #d0, #57
db #c9, #94, #d3, #86, #e3, #92, #e3, #52
db #01, #00, #1d, #c8, #00, #12, #d1, #aa
db #d0, #d8, #29, #f8, #b6, #01, #3a, #b6
db #f0, #06, #c9, #06, #d1, #07, #f9, #57
db #c1, #57, #e1, #a9, #d9, #f9, #d1, #9c
db #f2, #38, #d3, #36, #f3, #42, #f3, #86
db #c3, #85, #db, #0f, #0f, #00, #17, #d0
db #17, #d8, #6c, #e8, #6d, #e8, #68, #f0
db #6b, #f0, #bc, #e8, #09, #d1, #0c, #d9
db #5b, #f1, #35, #e3, #3a, #e3, #3a, #eb
db #87, #d3, #d4, #c3, #1e, #08, #00, #18
db #d0, #1d, #d8, #6a, #c0, #0b, #d9, #0b
db #e9, #5a, #c9, #e8, #ca, #8a, #f3, #2c
db #03, #00, #1c, #d0, #e7, #e2, #db, #cb
db #a0, #02, #00, #16, #d8, #10, #c9, #3c
db #04, #bb, #d8, #8e, #e8, #9b, #53, #55
db #e7, #cb, #c0, #06, #00, #19, #d8, #1a
db #e0, #16, #e8, #1a, #e8, #10, #f9, #93
db #db, #87, #05, #00, #17, #e0, #f8, #d1
db #e8, #d2, #36, #cb, #83, #db, #a4, #01
db #c9, #e0, #a5, #01, #9e, #dd, #e8, #e1
db #02, #d3, #a3, #f0, #c1, #4b, #a9, #f3
db #18, #8a, #88, #db, #48, #1a, #23, #d3
db #25, #a3, #cc, #f8, #0c, #86, #5d, #1a
db #f8, #98, #e2, #21, #26, #67, #c0, #b8
db #93, #09, #f1, #89, #c3, #01, #08, #a9
db #eb, #c8, #4b, #f8, #bd, #f8, #09, #f9
db #eb, #da, #89, #cb, #85, #d3, #91, #e3
db #03, #a8, #db, #d0, #1b, #e0, #09, #e9
db #ab, #c1, #e6, #fa, #82, #e3, #83, #fd
db #eb, #d1, #cb, #10, #a4, #d6, #68, #e0
db #0c, #d1, #06, #d9, #06, #e9, #0c, #f1
db #fd, #f9, #97, #d2, #3a, #cb, #91, #eb
db #0b, #63, #6b, #82, #0b, #f1, #40, #05
db #a8, #b6, #1e, #b6, #c8, #85, #c3, #88
db #fb, #d5, #c3, #07, #af, #01, #ac, #b8
db #c8, #09, #e1, #5a, #d1, #5b, #f9, #e8
db #e2, #39, #cb, #39, #d3, #3a, #d3, #82
db #43, #a9, #5b, #b8, #2a, #39, #f3, #88
db #eb, #20, #01, #b6, #b2, #d0, #75, #a1
db #c0, #72, #ab, #d1, #35, #eb, #37, #eb
db #38, #f3, #37, #fb, #68, #bf, #a6, #bc
db #ee, #c2, #b4, #2a, #be, #c8, #0e, #5d
db #a2, #0b, #c1, #e1, #e7, #87, #e7, #da
db #84, #fb, #16, #86, #cf, #0e, #c1, #fd
db #e9, #61, #27, #0c, #c9, #60, #98, #8f
db #0d, #c9, #70, #9a, #0e, #f1, #ae, #06
db #f7, #e1, #f7, #e9, #49, #c2, #47, #e2
db #b0, #21, #ea, #5b, #d1, #36, #fb, #12
db #f3, #5e, #76, #4d, #c2, #06, #e7, #16
db #5a, #f1, #39, #c3, #84, #f3, #dd, #c3
db #dc, #cb, #96, #27, #5d, #f9, #08, #97
db #15, #f9, #c1, #e9, #da, #8c, #fb, #30
db #96, #ef, #47, #f2, #e6, #c2, #87, #c3
db #90, #25, #e1, #9b, #da, #35, #fb, #e5
db #db, #84, #4b, #89, #98, #ea, #e9, #e2
db #85, #89, #e7, #c2, #a1, #d9, #e6, #d2
db #86, #b5, #a8, #e9, #fa, #84, #50, #9e
db #35, #f3, #05, #f7, #66, #86, #d3, #49
db #22, #db, #c2, #ae, #87, #2d, #cb, #18
db #8a, #db, #8a, #eb, #c3, #ab, #8b, #42
db #98, #c1, #85, #fb, #83, #aa, #dc, #1c
db #2a, #d9, #8a, #78, #e2, #c0, #00, #20
delta12:
db #00, #7a, #1b, #c0, #06, #00, #1a, #c8
db #ab, #d1, #fb, #e9, #9b, #e2, #98, #ea
db #95, #d3, #80, #04, #00, #16, #d0, #1a
db #e0, #1a, #e8, #9a, #da, #d0, #01, #00
db #19, #d8, #5a, #f7, #1d, #8a, #1e, #02
db #1e, #6d, #c0, #9b, #f2, #f0, #df, #62
db #6e, #e0, #4b, #62, #6b, #f8, #0f, #4e
db #bd, #e8, #36, #cb, #10, #03, #f3, #4e
db #f0, #9b, #da, #35, #fb, #41, #d7, #7a
db #0c, #c9, #24, #f7, #0d, #ba, #50, #ec
db #22, #d1, #0d, #08, #0b, #e1, #90, #02
db #83, #f1, #4b, #c2, #07, #86, #89, #5c
db #d1, #fd, #e9, #21, #27, #ab, #c1, #1a
db #aa, #b3, #ad, #87, #26, #ad, #d9, #04
db #26, #ae, #f9, #60, #26, #fb, #f1, #a0
db #21, #e6, #fc, #f9, #4c, #d2, #02, #d5
db #4d, #c2, #30, #26, #4e, #ea, #2c, #26
db #98, #ca, #0c, #21, #89, #98, #d2, #8e
db #cb, #0e, #88, #8d, #fb, #01, #c0, #e1
db #db, #00, #20
delta13:
db #5e, #1b, #1e, #01, #00, #1d, #d8, #b4
db #f6, #98, #e0, #0f, #02, #83, #f8, #e8
db #d2, #c0, #92, #f3, #66, #c8, #66, #d0
db #0b, #01, #9e, #6b, #f8, #f0, #e9, #a6
db #be, #f7, #e9, #78, #e3, #f3, #d0, #2d
db #98, #f7, #bd, #e8, #c1, #98, #0c, #c9
db #34, #82, #0d, #f9, #70, #03, #90, #89
db #5e, #c9, #f7, #f1, #47, #f2, #96, #01
db #00, #ad, #d9, #0c, #8a, #ac, #e9, #40
db #fb, #ba, #16, #ed, #fd, #a6, #e0, #39
db #f1, #20, #7d, #88, #fc, #f9, #80, #6e
db #4b, #c2, #4c, #c2, #38, #f3, #3a, #ca
db #03, #f7, #4d, #a2, #50, #62, #4e, #e2
db #00, #62, #97, #d2, #84, #62, #98, #ea
db #c2, #ae, #9e, #48, #ad, #62, #ee, #e2
db #30, #18, #85, #cb, #88, #f3, #86, #8c
db #dc, #cb, #00, #02
delta14:
db #13, #13, #3c, #02, #00, #6d, #c0, #98
db #c2, #1e, #01, #8e, #f3, #e0, #0f, #f7
db #62, #bd, #e8, #a1, #0a, #f8, #c1, #70
db #03, #42, #f7, #d1, #f7, #e9, #37, #f3
db #0e, #01, #00, #fd, #d1, #07, #66, #e9
db #f0, #02, #87, #f7, #f1, #47, #e2, #12
db #98, #df, #4d, #c2, #28, #98, #4b, #ca
db #c0, #9a, #4c, #d2, #9e, #89, #81, #89
db #e6, #d2, #87, #e6, #35, #db, #80, #ed
db #36, #f3, #50, #62, #fb, #01, #26, #89
db #c3, #b0, #26, #85, #cb, #41, #63, #cb
db #10, #ac, #c5, #e1, #00, #02
delta15:
db #14, #3a, #2e, #f0, #08, #00, #1e, #c0
db #f7, #e1, #f7, #e9, #fb, #f1, #49, #c2
db #4a, #c2, #35, #eb, #37, #f3, #d0, #04
db #00, #19, #c8, #c0, #f8, #35, #f3, #85
db #cb, #0c, #02, #00, #1c, #c8, #98, #ca
db #a0, #01, #00, #1a, #d0, #b4, #f7, #1d
db #b8, #e0, #df, #78, #19, #d8, #9b, #e2
db #3c, #d1, #62, #1d, #d8, #5c, #f9, #96
db #08, #19, #e0, #1e, #06, #90, #e1, #e0
db #1d, #f8, #6d, #c0, #0d, #e9, #fd, #d1
db #fd, #d9, #c0, #cb, #e8, #1a, #e8, #ab
db #c9, #68, #f2, #22, #f0, #fe, #d9, #0e
db #4a, #6a, #c0, #db, #cb, #0f, #03, #55
db #e1, #6d, #c8, #ad, #c1, #8d, #f3, #80
db #05, #00, #66, #d0, #4c, #d2, #49, #f2
db #35, #fb, #e5, #db, #00, #04, #00, #67
db #e0, #b6, #c0, #aa, #c9, #f7, #c1, #78
db #b1, #89, #be, #e0, #0d, #f1, #01, #7a
db #bd, #f0, #e1, #c3, #e1, #db, #08, #83
db #c8, #9c, #4d, #ea, #10, #81, #eb, #88
db #0c, #d1, #97, #d2, #9c, #fa, #d8, #c3
db #1c, #62, #0b, #e1, #f8, #f9, #b0, #62
db #0c, #f1, #43, #79, #5c, #d1, #03, #a9
db #62, #ab, #c1, #4e, #ea, #ed, #c2, #0b
db #18, #ad, #c9, #dc, #d3, #21, #96, #ad
db #e9, #9c, #e2, #36, #c3, #89, #a2, #e1
db #e3, #f8, #cb, #c2, #60, #a9, #cf, #ff
db #e1, #86, #c3, #30, #e5, #89, #f7, #d1
db #47, #fa, #20, #ae, #fa, #e9, #85, #83
db #7f, #a8, #fd, #c2, #98, #fe, #f1, #70
db #b8, #47, #a3, #39, #fb, #37, #fb, #2c
db #35, #a1, #98, #c2, #98, #e7, #e7, #d2
db #38, #cb, #50, #cf, #9b, #da, #84, #98
db #f7, #98, #e2, #48, #aa, #9e, #41, #26
db #9c, #ea, #4a, #2a, #9e, #98, #24, #89
db #fa, #87, #e1, #36, #cb, #07, #71, #e6
db #35, #db, #39, #f3, #61, #df, #37, #e3
db #90, #63, #f3, #86, #98, #c1, #84, #fb
db #18, #98, #8c, #fb, #06, #ac, #dc, #00
db #02
delta16:
db #5e, #08, #d0, #01, #00, #19, #d0, #20
db #f7, #a8, #1a, #f0, #98, #19, #d8, #0f
db #98, #0d, #d9, #70, #98, #f7, #e1, #0c
db #99, #98, #e2, #c0, #88, #ea, #80, #8c
db #ec, #c2, #00, #02
delta17:
db #41, #17, #3d, #00, #04, #01, #16, #c0
db #2f, #c0, #48, #c0, #61, #c0, #16, #c8
db #2f, #c8, #48, #c8, #61, #c8, #16, #d0
db #17, #d0, #2f, #d0, #30, #d0, #48, #d0
db #49, #d0, #61, #d0, #62, #d0, #16, #d8
db #2f, #d8, #48, #d8, #18, #e0, #2f, #e0
db #31, #e0, #4a, #e0, #4e, #e0, #61, #e0
db #63, #e0, #67, #e0, #68, #e0, #16, #e8
db #18, #e8, #1c, #e8, #2f, #e8, #31, #e8
db #4a, #e8, #61, #e8, #63, #e8, #68, #e8
db #32, #f0, #36, #f0, #4b, #f0, #4f, #f0
db #61, #f0, #64, #f0, #68, #f0, #16, #f8
db #19, #f8, #2f, #f8, #32, #f8, #4b, #f8
db #64, #f8, #80, #c0, #89, #c0, #a2, #c0
db #bb, #c0, #9c, #c8, #b5, #c8, #b6, #c8
db #6a, #d0, #6b, #d0, #84, #d0, #9c, #d0
db #9d, #d0, #b5, #d0, #b6, #d0, #ba, #d0
db #69, #d8, #6a, #d8, #6b, #d8, #6f, #d8
db #82, #d8, #83, #d8, #84, #d8, #85, #d8
db #88, #d8, #9b, #d8, #9c, #d8, #9d, #d8
db #a1, #d8, #b4, #d8, #b5, #d8, #b6, #d8
db #ba, #d8, #68, #e0, #69, #e0, #6b, #e0
db #6f, #e0, #80, #e0, #81, #e0, #82, #e0
db #84, #e0, #85, #e0, #88, #e0, #99, #e0
db #9a, #e0, #9e, #e0, #a1, #e0, #b2, #e0
db #b3, #e0, #b7, #e0, #ba, #2a, #d9, #e8
db #2a, #e8, #80, #81, #aa, #83, #84, #aa
db #85, #99, #aa, #9a, #9c, #a9, #9d, #1c
db #9e, #e8, #b2, #e8, #b3, #e8, #b5, #e8
db #b6, #e8, #b7, #e8, #67, #a0, #d4, #51
db #e2, #69, #f0, #6a, #f0, #6b, #f0, #6c
db #f0, #80, #f0, #81, #f0, #82, #f0, #83
db #f0, #84, #f0, #85, #f0, #99, #f0, #9a
db #f0, #9b, #f0, #9c, #f0, #9d, #f0, #9e
db #f0, #b2, #f0, #b3, #f0, #b4, #f0, #b5
db #f0, #b6, #f0, #b7, #f0, #67, #f8, #68
db #f8, #69, #f8, #6b, #f8, #6c, #f8, #80
db #f8, #81, #f8, #82, #f8, #84, #f8, #85
db #f8, #99, #f8, #9a, #f8, #9b, #f8, #9e
db #f8, #b2, #f8, #b3, #f8, #b7, #fd, #c0
db #b8, #aa, #bc, #d0, #aa, #d1, #d5, #aa
db #e9, #ea, #aa, #eb, #ed, #a9, #ee, #aa
db #03, #c1, #05, #06, #a6, #07, #b8, #c8
db #ba, #aa, #bb, #d1, #aa, #d3, #d4, #aa
db #ea, #ec, #a9, #ed, #aa, #02, #c9, #03
db #05, #a6, #06, #b6, #d0, #b7, #aa, #b8
db #bb, #aa, #cf, #d0, #aa, #d1, #d2, #aa
db #d4, #d7, #aa, #e8, #e9, #aa, #ea, #eb
db #aa, #ed, #f0, #9a, #01, #d1, #02, #aa
db #03, #04, #aa, #06, #09, #6a, #b6, #d8
db #b7, #b8, #aa, #b9, #bb, #aa, #be, #cf
db #aa, #d0, #d1, #aa, #d2, #d7, #aa, #e8
db #e9, #aa, #ea, #eb, #a9, #ec, #aa, #01
db #d9, #02, #03, #a6, #04, #b6, #e0, #b7
db #aa, #b8, #b9, #aa, #cf, #d0, #aa, #d1
db #d2, #aa, #d4, #e7, #aa, #e8, #e9, #9a
db #00, #e1, #01, #a6, #02, #b5, #e8, #b6
db #aa, #b7, #ce, #aa, #cf, #e7, #9a, #dd
db #f0, #f6, #a1, #f7, #a6, #0f, #f1, #c4
db #f8, #fc, #01, #c1, #1a, #aa, #56, #5f
db #a9, #60, #4a, #2d, #07, #00, #17, #c0
db #52, #c8, #35, #d0, #4e, #e8, #1b, #f8
db #cd, #02, #02, #c1, #0f, #be, #00, #19
db #c0, #1a, #aa, #1c, #1d, #aa, #1e, #1f
db #aa, #20, #31, #aa, #32, #33, #aa, #35
db #36, #aa, #37, #38, #aa, #39, #4a, #aa
db #4b, #4c, #aa, #4e, #4f, #aa, #50, #51
db #aa, #52, #63, #aa, #64, #65, #aa, #68
db #69, #a9, #6a, #aa, #18, #c8, #19, #1a
db #aa, #1d, #1e, #aa, #1f, #33, #aa, #36
db #37, #aa, #38, #4a, #aa, #4c, #4f, #aa
db #50, #51, #aa, #65, #68, #aa, #69, #6a
db #a6, #6b, #1a, #d0, #1d, #aa, #1e, #1f
db #aa, #37, #38, #aa, #4f, #50, #aa, #51
db #69, #a9, #6a, #aa, #1d, #d8, #1e, #1f
db #aa, #37, #38, #aa, #50, #51, #aa, #66
db #69, #a6, #6a, #1e, #e0, #1f, #aa, #37
db #38, #aa, #50, #51, #aa, #69, #6a, #9a
db #1f, #e8, #38, #aa, #51, #6a, #6a, #1f
db #f0, #38, #4d, #aa, #51, #69, #a9, #6a
db #aa, #1e, #f8, #1f, #34, #aa, #37, #38
db #aa, #50, #51, #aa, #66, #67, #aa, #68
db #69, #a6, #6a, #6b, #c0, #6c, #aa, #6d
db #6e, #aa, #6f, #84, #aa, #85, #86, #aa
db #87, #88, #aa, #9e, #9f, #aa, #a0, #b8
db #a9, #b9, #aa, #6d, #c8, #6e, #86, #aa
db #87, #9e, #aa, #9f, #a0, #aa, #b7, #b8
db #a6, #b9, #6e, #d0, #86, #aa, #87, #9f
db #aa, #a0, #b8, #9a, #6d, #d8, #86, #6a
db #88, #f0, #9f, #a0, #aa, #b8, #b9, #a9
db #ba, #aa, #6d, #f8, #6e, #86, #aa, #87
db #88, #aa, #9f, #a0, #aa, #b8, #b9, #a6
db #ba, #be, #c0, #bf, #aa, #d7, #d8, #aa
db #f0, #f1, #9a, #09, #c1, #0a, #6a, #be
db #c8, #d7, #f0, #96, #09, #c9, #07, #d9
db #bc, #e0, #bd, #aa, #d7, #ef, #a9, #f0
db #a6, #04, #e1, #08, #b9, #e8, #d1, #aa
db #d2, #ea, #aa, #eb, #ef, #9a, #02, #e9
db #03, #a6, #04, #b6, #f0, #b7, #aa, #b8
db #b9, #aa, #ce, #cf, #aa, #d0, #d1, #aa
db #e7, #ea, #a8, #ff, #0a, #b4, #f8, #e0
db #15, #00, #21, #c0, #3a, #aa, #53, #6c
db #68, #21, #c8, #3a, #0a, #21, #d0, #39
db #f0, #bb, #f8, #c0, #c0, #29, #c8, #d9
db #1a, #0a, #c9, #da, #e0, #dc, #f0, #0e
db #f1, #dc, #f8, #09, #f9, #0a, #69, #0e
db #c1, #44, #2a, #1e, #19, #00, #30, #c0
db #49, #6b, #81, #a6, #4e, #f8, #a1, #c0
db #b9, #d0, #a1, #e8, #ba, #6f, #f0, #a1
db #9a, #6f, #f8, #a1, #6a, #bf, #c8, #d8
db #f1, #96, #d5, #d8, #d6, #e0, #05, #e1
db #09, #86, #be, #e8, #07, #e9, #08, #84
db #28, #d2, #f0, #01, #f1, #1b, #c1, #07
db #15, #00, #62, #c0, #31, #c8, #4a, #d0
db #64, #e0, #19, #e8, #1e, #6a, #37, #f0
db #6a, #c0, #9d, #b6, #96, #9b, #c8, #85
db #d0, #68, #d8, #9f, #86, #08, #c1, #bd
db #c8, #d6, #a8, #ef, #68, #05, #e9, #e4
db #f8, #fd, #00, #aa, #43, #05, #00, #17
db #c8, #1e, #f0, #e9, #e8, #01, #e9, #fe
db #f0, #2c, #02, #00, #20, #c8, #ed, #e0
db #03, #14, #00, #30, #c8, #63, #d0, #18
db #d8, #32, #e8, #37, #4b, #a6, #69, #33
db #f8, #4c, #a9, #65, #a6, #6b, #c8, #82
db #9a, #d0, #b3, #96, #6e, #f0, #d4, #d8
db #ba, #e0, #eb, #90, #aa, #b8, #e8, #cb
db #f8, #3c, #0e, #00, #39, #c8, #20, #d0
db #39, #52, #a6, #6b, #20, #d8, #39, #aa
db #52, #6b, #b6, #20, #cf, #32, #e8, #bc
db #1f, #0e, #f8, #1d, #c1, #01, #6d, #a6
db #49, #62, #c8, #31, #96, #4a, #d8, #63
db #d8, #4f, #e0, #50, #96, #64, #e8, #19
db #f0, #b4, #c0, #84, #aa, #9d, #81, #95
db #ea, #6c, #d8, #00, #e9, #cd, #f0, #e6
db #f0, #b2, #f8, #33, #c1, #4c, #c1, #a4
db #09, #53, #0a, #6c, #c8, #0b, #03, #6b
db #1b, #d0, #31, #b5, #08, #fd, #75, #06
db #00, #3a, #d0, #9a, #a3, #ba, #a0, #e8
db #9e, #2f, #02, #f9, #48, #c2, #86, #d0
db #eb, #f8, #c0, #17, #a9, #6c, #1d, #3a
db #d8, #53, #d8, #21, #e0, #20, #f8, #52
db #f8, #bc, #d8, #a3, #6a, #ca, #bc, #90
db #7d, #71, #e8, #c2, #e8, #0e, #e9, #c3
db #f8, #f4, #f8, #0e, #f9, #1f, #c1, #22
db #c1, #27, #c1, #3e, #c1, #40, #c1, #57
db #c1, #59, #f5, #20, #03, #69, #21, #d8
db #1d, #e8, #53, #f0, #b3, #d8, #d3, #a2
db #04, #f9, #3d, #04, #01, #ee, #0f, #83
db #f7, #0a, #36, #e0, #78, #05, #a1, #39
db #39, #52, #e0, #02, #f1, #34, #c1, #4d
db #c1, #80, #1d, #23, #2a, #e0, #53, #e0
db #6c, #e0, #3a, #5a, #b6, #6b, #25, #2c
db #d0, #b4, #0b, #c9, #6e, #80, #4f, #e2
db #e0, #87, #e0, #a0, #e0, #8a, #e8, #a3
db #fd, #f0, #bc, #9a, #71, #f8, #8a, #a5
db #a3, #6a, #be, #d0, #c1, #e0, #0d, #e1
db #f5, #e8, #0f, #e9, #05, #f1, #09, #c1
db #2c, #5e, #84, #aa, #69, #05, #00, #6b
db #e0, #52, #e8, #b9, #f8, #ce, #cf, #4a
db #0e, #0f, #00, #1b, #e8, #34, #aa, #36
db #4d, #6a, #83, #c0, #9c, #ba, #9a, #6f
db #c8, #88, #a0, #b4, #69, #69, #d0, #ec
db #e0, #06, #e1, #d4, #f0, #e8, #28, #68
db #03, #00, #20, #e8, #39, #4a, #6d, #f0
db #08, #0b, #00, #21, #e8, #1c, #f0, #67
db #1a, #70, #c0, #6a, #c8, #83, #6a, #6f
db #d0, #88, #9b, #80, #68, #bc, #e8, #bb
db #e0, #70, #0b, #00, #35, #e8, #a2, #e0
db #bf, #e8, #06, #f1, #05, #f9, #aa, #c1
db #0a, #1e, #aa, #37, #50, #a1, #5c, #a5
db #c3, #04, #00, #4f, #e8, #0a, #f8, #88
db #e8, #d1, #f8, #0c, #06, #00, #67, #e8
db #b5, #c0, #ba, #c8, #82, #d0, #bc, #d8
db #08, #d9, #f0, #76, #00, #6b, #e8, #6d
db #e0, #86, #aa, #9f, #b8, #a6, #bb, #6d
db #e8, #70, #aa, #86, #89, #aa, #9f, #a2
db #a9, #bb, #aa, #70, #f0, #89, #a2, #a6
db #bb, #70, #f8, #89, #a8, #a2, #6a, #f2
db #c0, #c0, #e0, #d9, #f2, #a9, #f3, #a6
db #0b, #e1, #0c, #c0, #e8, #c1, #aa, #d9
db #da, #aa, #db, #f0, #aa, #f2, #f3, #a9
db #f4, #aa, #09, #e9, #0a, #0b, #aa, #0c
db #0d, #6a, #bd, #f0, #be, #bf, #aa, #c0
db #c1, #aa, #d6, #d7, #aa, #d8, #d9, #aa
db #da, #ee, #aa, #ef, #f0, #aa, #f1, #f2
db #a9, #f3, #aa, #07, #f1, #08, #09, #aa
db #0a, #0b, #a6, #0c, #b7, #f8, #bc, #aa
db #bd, #be, #aa, #bf, #c0, #aa, #c1, #c2
db #aa, #d0, #d5, #aa, #d6, #d7, #aa, #d8
db #d9, #aa, #da, #ee, #aa, #ef, #f0, #aa
db #f1, #f2, #a9, #f3, #aa, #00, #f9, #03
db #07, #aa, #08, #0d, #6a, #04, #c1, #06
db #08, #aa, #0c, #0d, #aa, #10, #11, #aa
db #12, #1c, #aa, #21, #25, #aa, #26, #29
db #aa, #2a, #2b, #aa, #35, #36, #aa, #3a
db #3c, #aa, #3f, #42, #aa, #43, #4e, #aa
db #4f, #51, #aa, #52, #53, #ab, #54, #58
db #c6, #9d, #b2, #b5, #f0, #ab, #e7, #b5
db #5a, #4a, #9a, #35, #f0, #fe, #3a, #e5
db #e1, #34, #86, #00, #f1, #b6, #f8, #24
db #2a, #35, #7a, #bc, #d0, #42, #d3, #4e
db #a8, #d2, #ab, #87, #08, #50, #18, #9f
db #4b, #87, #f0, #bd, #c0, #d6, #c0, #ef
db #c0, #ed, #d8, #52, #55, #02, #01, #00
db #52, #f0, #96, #04, #00, #66, #aa, #d3
db #eb, #9e, #e6, #f8, #16, #ad, #4a, #1a
db #f8, #d8, #e0, #10, #09, #d0, #1d, #d0
db #01, #ad, #d8, #6c, #84, #05, #bc, #e0
db #d6, #d8, #05, #d9, #d0, #e8, #bb, #f8
db #1c, #67, #70, #eb, #6b, #97, #d8, #23
db #1f, #e0, #e1, #ab, #e7, #36, #d9, #c0
db #79, #c4, #71, #8e, #f1, #c1, #84, #95
db #ac, #39, #f2, #0d, #85, #6a, #85, #69
db #6b, #b7, #d0, #6f, #bd, #c3, #a3, #a7
db #1a, #4d, #81, #a2, #28, #0a, #87, #d8
db #30, #0a, #a0, #89, #5c, #0b, #d9, #bf
db #e0, #d4, #f8, #ec, #f8, #ed, #f8, #0c
db #f9, #23, #c1, #28, #c1, #41, #c1, #d0
db #2a, #0f, #6e, #a0, #87, #aa, #5d, #c1
db #a1, #01, #b9, #ec, #d2, #89, #0b, #23
db #4b, #e1, #e8, #42, #a9, #29, #f9, #55
db #4a, #20, #09, #00, #bc, #c8, #d5, #c8
db #ee, #c8, #07, #c9, #d5, #d0, #ee, #18
db #07, #d1, #bd, #d8, #f5, #c9, #b5, #6b
db #d6, #d0, #ef, #d0, #08, #d1, #ba, #49
db #9f, #9d, #06, #d9, #b0, #17, #46, #bf
db #72, #04, #f1, #92, #1b, #4a, #f1, #e8
db #06, #f9, #b4, #04, #aa, #d5, #03, #a7
db #b8, #39, #c1, #86, #2a, #59, #ec, #ca
db #c2, #db, #e9, #9e, #3b, #c1, #38, #f3
db #aa, #ff, #20, #f5, #50, #db, #2e, #ad
db #90, #f7, #ae, #0f, #60, #df, #6c, #24
db #c1, #5a, #00, #02
table_delta:
dw delta00, delta01, delta02, delta03, delta04, delta05, delta06, delta07
dw delta08, delta09, delta10, delta11, delta12, delta13, delta14, delta15
dw delta16, delta17
palette:
db 00, 13, 26
end

save'disc.bin',#200, end - start,DSK,'delta.dsk'