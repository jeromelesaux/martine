package assembly

const (
	WaitVbl = `;---- attente vbl ----------
waitvbl
	ld b,#f5 ; attente vbl
vbl
	in a,(c)
	rra
	jp nc,vbl
	ret
;---------------------------
`
)
