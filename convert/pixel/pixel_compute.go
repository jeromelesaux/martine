package pixel

// PixelMode0 converts palette position into byte in  screen mode 0
func PixelMode0(pp1, pp2 int) byte {
	var pixel byte
	if uint8(pp1)&1 == 1 {
		pixel += 128
	}
	if uint8(pp1)&2 == 2 {
		pixel += 8
	}
	if uint8(pp1)&4 == 4 {
		pixel += 32
	}
	if uint8(pp1)&8 == 8 {
		pixel += 2
	}
	if uint8(pp2)&1 == 1 {
		pixel += 64
	}
	if uint8(pp2)&2 == 2 {
		pixel += 4
	}
	if uint8(pp2)&4 == 4 {
		pixel += 16
	}
	if uint8(pp2)&8 == 8 {
		pixel++
	}
	return pixel
}

// PixelMode1 converts palette position into byte in  screen mode 1
func PixelMode1(pp1, pp2, pp3, pp4 int) byte {
	var pixel byte
	if uint8(pp1)&1 == 1 {
		pixel += 128
	}
	if uint8(pp1)&2 == 2 {
		pixel += 8
	}
	if uint8(pp2)&1 == 1 {
		pixel += 64
	}
	if uint8(pp2)&2 == 2 {
		pixel += 4
	}
	if uint8(pp3)&1 == 1 {
		pixel += 32
	}
	if uint8(pp3)&2 == 2 {
		pixel += 2
	}
	if uint8(pp4)&1 == 1 {
		pixel += 16
	}
	if uint8(pp4)&2 == 2 {
		pixel++
	}
	return pixel
}

// PixelMode converts palette position into byte in screen mode 2
func PixelMode2(pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 int) byte {
	var pixel byte
	if uint8(pp1)&1 == 1 {
		pixel += 128
	}
	if uint8(pp2)&1 == 1 {
		pixel += 64
	}
	if uint8(pp3)&1 == 1 {
		pixel += 32
	}
	if uint8(pp4)&1 == 1 {
		pixel += 16
	}
	if uint8(pp5)&1 == 1 {
		pixel += 8
	}
	if uint8(pp6)&1 == 1 {
		pixel += 4
	}
	if uint8(pp7)&1 == 1 {
		pixel += 2
	}
	if uint8(pp8)&1 == 1 {
		pixel++
	}
	return pixel
}

// RawPixelMode2 converts color  byte in palette position in  screen mode 2
func RawPixelMode2(b byte) (pp1, pp2, pp3, pp4, pp5, pp6, pp7, pp8 int) {
	val := int(b)
	if val-128 >= 0 {
		pp1 = 1
		val -= 128
	}
	if val-64 >= 0 {
		pp2 = 1
		val -= 64
	}
	if val-32 >= 0 {
		pp3 = 1
		val -= 32
	}
	if val-16 >= 0 {
		pp4 = 1
		val -= 16
	}
	if val-8 >= 0 {
		pp5 = 1
		val -= 8
	}
	if val-4 >= 0 {
		pp6 = 1
		val -= 4
	}
	if val-2 >= 0 {
		pp7 = 1
		val -= 2
	}
	if val-1 >= 0 {
		pp8 = 1
	}
	return
}

// RawPixelMode1 converts color  byte in palette position in screen mode 1
func RawPixelMode1(b byte) (pp1, pp2, pp3, pp4 int) {
	val := int(b)
	if val-128 >= 0 {
		pp1 |= 1
		val -= 128
	}
	if val-64 >= 0 {
		pp2 |= 1
		val -= 64
	}
	if val-32 >= 0 {
		pp3 |= 1
		val -= 32
	}
	if val-16 >= 0 {
		pp4 |= 1
		val -= 16
	}
	if val-8 >= 0 {
		pp1 |= 2
		val -= 8
	}
	if val-4 >= 0 {
		pp2 |= 2
		val -= 4
	}
	if val-2 >= 0 {
		pp3 |= 2
		val -= 2
	}
	if val-1 >= 0 {
		pp4 |= 2
	}

	return
}

// RawPixelMode0 converts color  byte in palette position in screen mode 0
func RawPixelMode0(b byte) (pp1, pp2 int) {
	val := int(b)
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-128 >= 0 {
		pp1 |= 1
		val -= 128
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-64 >= 0 {
		pp2 |= 1
		val -= 64
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-32 >= 0 {
		pp1 |= 4
		val -= 32
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-16 >= 0 {
		pp2 |= 4
		val -= 16
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-8 >= 0 {
		pp1 |= 2
		val -= 8
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-4 >= 0 {
		pp2 |= 2
		val -= 4
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-2 >= 0 {
		pp1 |= 8
		val -= 2
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	if val-1 >= 0 {
		pp2 |= 8
	}
	//fmt.Fprintf(os.Stderr,"v:%.8b\n",val)
	return
}
