		org &1000
		nolist
		run start
start
		di
		ld (pile+1),sp
		ld hl,&c000
		ld sp,debutrle16
pool
		pop bc
		pop de
loop
		ld (hl),e
		inc hl
		dec bc
		ld a,b
		or c
		jr nz,loop
		dec sp
		ex de,hl
		ld hl,finrle16
		sbc hl,sp
		ex de,hl
 		jr nz,pool

pile		ld sp,0
		ei
		ret

		org &4000
		nolist
debutrle16
incbin		"fichier.scr"
finrle16
