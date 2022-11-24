package address

import "math"

// CpcScreenAddress returns the screen address according the screen mode, the initialAddress (always #C000)
// x the column number and y the line number on the screen
func CpcScreenAddress(intialeAddresse int, x, y int, mode uint8, isOverscan, doubleScreen bool) int {
	var addr int
	var adjustMode int
	switch mode {
	case 0:
		adjustMode = 2
	case 1:
		adjustMode = 4
	case 2:
		adjustMode = 8
	}
	if isOverscan {
		addr = (0x800 * (y % 8)) + (0x60 * (y / 8)) + ((x + 1) / adjustMode)
		if doubleScreen && y > 167 {
			addr += 0x3800
		} else {
			if !doubleScreen && y > 127 {
				addr += 0x3800
			}
		}
	} else {
		//	addr = (y>>4)*0x40 + (y&14)*0x400
		addr = (0x800 * (y % 8)) + (0x50 * (y / 8)) + ((x + 1) / adjustMode)
	}
	if intialeAddresse == 0 {
		return addr
	}
	return intialeAddresse + addr
}

func CpcScreenAddressOffset(line int) int {
	return int(math.Floor(float64(line)/8)*80) + (line%8)*2048
}
