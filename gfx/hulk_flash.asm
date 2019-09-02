 org #3000 
 run $ 
 not list
start
	DI
	
	LD	HL,(#38) ; sauvegarde de RST 7 interruption clavier
	LD	(INTER+1),HL
	LD	HL,#C9FB
	LD	(#38),HL
	EXX	
	EX	AF,AF'
	PUSH	AF
	PUSH	BC
	PUSH	DE
	PUSH	HL
	EI	

debut
	ld b,#f5 ; synchronisation vbl 
sync	in a,(c)
	rra
	jr nc,sync

	ld b,#7f ; changement de mode 
mode	ld a,#8d
	xor 1 ; #8d/8c
	ld (mode+1),a
	out (c),a
;
	ld bc,#bc0c ; switch de bank 
	out (c),c
flip	ld a,#10
 	xor #20; #30/#10
	ld (flip+1),a
	inc b
	out (c),a
	bit 5,a
	jr z,pal2

	ld hl,palette2 ; pointe la palette
sendpal
	ld a, 16
        call setpalette

	jr haltsync


pal2
	ld a,4
        ld hl,palette1 ; pointe la palette
	jr sendpal

haltsync
	halt

key     LD   bc,#f40e                  ; Teste la barre Espace
        OUT  (c),c
        LD   bc,#f6c0
        OUT  (c),c
        XOR  a
        OUT  (c),a
        LD   bc,#f792
        OUT  (c),c
        LD   bc,#f645
        OUT  (c),c
        LD   b,#f4
        IN   a,(c)
        LD   bc,#f782
        OUT  (c),c
        LD   bc,#f600
        OUT  (c),c
        RLA 
	JP C,DEBUT

FIN
	DI	
INTER	LD	HL,0
	LD	(#38),HL
	POP	HL
	POP	DE
	POP	BC
	POP	AF
	EX	AF,AF'
	EXX	
	EI	
	JP	#BCA7
	ret

setpalette
	ld d,a
        ld bc,#7f00 ; set de la palette
setcolor
	out (c),c
        ld e,(hl) ; copie de la valeur de la premiere couleur
	out (c),e
        inc c ; incremente le numero de pen 
        inc l ; incremente la valeur de la couleur 
        dec a
        jr nz, setcolor
	ld a,d
	ret

palette1
db #54, #5E, #40, #43, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
palette2
db #4e, #4a, #43, #5c, #54, #40, #5e, #40, #4b, #56, #44, #59, #58, #46,0 ,0
end 

;save 'flash.bin',#3000,end-start,DSK,'HULK_reloaded.dsk'
save 'flash2.bin',#3000,end-start,AMSDOS
