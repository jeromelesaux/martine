package transformation

import (
	"bytes"
	"encoding/binary"
	"os"
	"sort"
)

type DeltaCollectionV2 struct {
	*DeltaCollection
}

type offset []uint16

func (f offset) Len() int {
	return len(f)
}

func (f offset) Less(i, j int) bool {
	return f[i] < f[j]
}

func (f offset) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (dc *DeltaCollectionV2) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := dc.Marshall()
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}

func (dc *DeltaCollectionV2) Marshall() ([]byte, error) {
	var b bytes.Buffer

	if err := binary.Write(&b, binary.LittleEndian, dc.OccurencePerFrame); err != nil {
		return b.Bytes(), err
	}
	if dc.OccurencePerFrame == 0 { // no difference between transitions
		return b.Bytes(), nil
	}
	// occurencesPerframe doit correspondre au nombre offsets modulo 255 et non au nombre d'items
	for _, item := range dc.Items {
		occ := len(item.Offsets)
		if err := binary.Write(&b, binary.LittleEndian, item.Byte); err != nil {
			return b.Bytes(), err
		}
		sort.Sort(offset(item.Offsets))
		var previousHB uint8 = 0
		lbs := make([]uint8, 0)
		for i := 0; i < occ; i++ {
			value := item.Offsets[i]
			currentHB := uint8(value >> 8)
			currentLB := uint8(value)
			if currentHB == previousHB {
				lbs = append(lbs, currentHB)
			} else {
				// export all the value
				if err := binary.Write(&b, binary.LittleEndian, currentHB); err != nil {
					return b.Bytes(), err
				}
				if err := binary.Write(&b, binary.LittleEndian, uint16(len(lbs))); err != nil {
					return b.Bytes(), err
				}
				for j := 0; j < len(lbs); j++ {
					if err := binary.Write(&b, binary.LittleEndian, lbs[j]); err != nil {
						return b.Bytes(), err
					}
				}
				lbs = make([]uint8, 0)
				previousHB = currentHB
				lbs = append(lbs, currentLB)
			}
			//			fmt.Fprintf(os.Stdout, "Value[%d]:%.4x\n", j, value)

		}
	}
	return b.Bytes(), nil
}
