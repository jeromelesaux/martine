org #3000
run $ 

start:
; delock asic 
di
ld bc,#bc00
ld hl,sequence
ld e,17
seq:
ld a,(hl)
out (c),a
inc hl
dec e
jr nz,seq
ei
;; page-in asic registers to #4000-#7fff
ld bc,#7fb8
out (c),c

ld hl,sprite_colours
ld de,#6400
ld bc,2
ldir

ld hl,#b7f9
jp #bcdd

;; page-out asic registers
ld bc,#7fa0
out (c),c

ret

sequence:
db #ff,#00,#ff,#77,#b3,#51,#a8,#d4,#62,#39,#9c,#46,#2b,#15,#8a,#cd,#ee

sprite_colours:
end:
db #66, #06, #63, #06, #00, #00, #96, #06, #33, #03, #63, #03, #93, #06, #96, #06
db #96, #09, #c9, #0c, #63, #06, #96, #06, #c6, #09, #c9, #09, #63, #03, #99, #09


save 'pal2plus.bin',#3000,end-start,AMSDOS
