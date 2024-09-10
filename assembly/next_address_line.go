package assembly

const (
	BC26 = `;---- recuperation de l'adresse de la ligne en dessous ------------
bc26
ld a,h
add a,8
ld h,a ; <---- le fameux que tu as oublié !
ret nc
ld bc,linewidth ; on passe en 96 colonnes
add hl,bc
res 3,h
ret
;-----------------------------------------------------------------`
)
