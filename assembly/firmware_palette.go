package assembly

const FirmwarePalette = `;--- application palette firmware -------------
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
;---------------------------------------------`
