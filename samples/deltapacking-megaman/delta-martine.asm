;BUILDSNA V2
;BANK 0
;SETCPC 2
;RUN start
;--- dimensions du sprite ----
large equ 100 / 2 
haut equ 100
;-----------------------------

org #1000
run $

start

    ; gestion du mode 
    ;ld bc,#7f8c
    ;out (c),c 
;--- selection du mode ---------
    ld a,0
    call #BC0E
;-------------------------------
  

;--- gestion de la palette ---- 
    call palettefirmware
;------------------------------

call xvbl

;--- affichage du sprite initiale --  
    ; affichage du premier sprite
    ld de,#C000 ; adresse de l'ecran 
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
ld hl,delta00
call delta

;call #bb06

call xvbl
ld hl,delta01
call delta

;call #bb06
call xvbl
ld hl,delta02
call delta

;call #bb06
call xvbl
ld hl,delta03
call delta

;call #bb06
call xvbl
ld hl,delta04
call delta

;call #bb06
call xvbl
ld hl,delta05
call delta

;call #bb06
call xvbl
ld hl,delta06
call delta

;call #bb06
call xvbl
ld hl,delta07
call delta

;call #bb06
call xvbl
ld hl,delta08
call delta

;call #bb06
call xvbl
ld hl,delta09
call delta



jp mainloop 


;--- routine de deltapacking --------------------------

delta
 ld a,(hl) ; nombre de byte a poker
 ld (nbbytepoked),a ; stockage en mémoire
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

 ld a,(nbbytepoked) ; reste t'il d'autres bytes a poker ? 
 dec a 
 ld (nbbytepoked),a
 jr nz,init
 ret



;---------------------------------------------------------------
;
; attente de plusieurs vbl
;
xvbl ld e,30
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
ld a,0
ld b,0
ld c,0
call #bc32

ld a,1
ld b,15
ld c,15
call #bc32

ld a,2
ld b,1
ld c,1
call #bc32

ld a,3
ld b,10
ld c,10
call #bc32

ld a,4
ld b,11
ld c,11
call #bc32

ld a,5
ld b,12
ld c,12
call #bc32

ld a,6
ld b,13
ld c,13
call #bc32

ld a,7
ld b,14
ld c,14
call #bc32

ld a,8
ld b,23
ld c,23
call #bc32

ld a,9
ld b,16
ld c,16
call #bc32

ld a,10
ld b,26
ld c,26
call #bc32


ret
;---------------------------------------------

;---- recuperation de l'adresse de la ligne en dessous ------------
bc26 
ld a,h
add a,8 
ld h,a ; <---- le fameux que tu as oublié !
ret nc 
ld bc,#c050 ; on passe en 96 colonnes
add hl,bc
res 3,h
ret
;-----------------------------------------------------------------


;--- variables memoires -----
pixel db 0 
nbbytepoked db 0
;----------------------------


;------- data ---------------------------

include 'data.asm'



end

save'delta.bin',#1000,end-start,DSK,'bomberman.dsk'
;save'delta.bin',#1000,end-start
