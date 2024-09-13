package transformation

import (
	"bytes"
	"encoding/binary"
	"sort"

	"github.com/jeromelesaux/martine/export/amsdos"
)

type DeltaCollectionV2 struct {
	*DeltaCollection
}

type DeltaV2 struct {
	HighByte uint8
	LowBytes []uint8
	Byte     uint8
}

func NewDeltaV2() *DeltaV2 {
	return &DeltaV2{LowBytes: make([]uint8, 0)}
}

func (d *DeltaV2) AddLowByte(v uint8) {
	d.LowBytes = append(d.LowBytes, v)
}

/*
	--- details de la structure ---
	nombre d'occurence par frame occ
	0	byte à poker uint8 |  valeur du HB uint8  | nombre de LB uint16 |  LB[0] ..... LB[nombre de LB] uint8
	1
	2
	.
	.
	.
	occ
*/

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
	b, err := dc.Marshall()
	if err != nil {
		return err
	}
	return amsdos.SaveOSFile(filename, b)
}

// nolint: funlen
func (dc *DeltaCollectionV2) Marshall() ([]byte, error) {
	var b bytes.Buffer
	/*
		if err := binary.Write(&b, binary.LittleEndian, dc.OccurencePerFrame); err != nil {
			return b.Bytes(), err
		}
	*/
	if dc.OccurencePerFrame == 0 { // no difference between transitions
		return b.Bytes(), nil
	}

	deltas := make([]*DeltaV2, 0)
	// occurencesPerframe doit correspondre au nombre offsets modulo 255 et non au nombre d'items
	for _, item := range dc.Items {
		occ := len(item.Offsets)
		sort.Sort(offset(item.Offsets))
		var previousHB uint8 = 0
		delta := NewDeltaV2()
		delta.Byte = item.Byte
		for i := 0; i < occ; i++ {
			value := item.Offsets[i]
			currentHB := uint8(value >> 8)
			currentLB := uint8(value)
			if currentHB == previousHB || i == 0 {
				delta.AddLowByte(currentLB)
				delta.HighByte = currentHB
			} else {
				deltas = append(deltas, delta)
				// export all the value HB
				/*		if err := binary.Write(&b, binary.LittleEndian, previousHB); err != nil {
							return b.Bytes(), err
						}
						// export the number of LB
						if err := binary.Write(&b, binary.LittleEndian, uint16(len(lowBytes))); err != nil {
							return b.Bytes(), err
						}
						// export the LB values
						for j := 0; j < len(lowBytes); j++ {
							if err := binary.Write(&b, binary.LittleEndian, lowBytes[j]); err != nil {
								return b.Bytes(), err
							}
						}
				*/
				delta = NewDeltaV2()
				delta.HighByte = currentHB
				delta.Byte = item.Byte
				delta.AddLowByte(currentLB)
			}
			//			log.GetLogger().Info( "Value[%d]:%.4x\n", j, value)

		}
		deltas = append(deltas, delta)

	}
	if err := binary.Write(&b, binary.LittleEndian, uint16(len(deltas))); err != nil {
		return b.Bytes(), err
	}
	for _, v := range deltas {

		if err := binary.Write(&b, binary.LittleEndian, v.Byte); err != nil {
			return b.Bytes(), err
		}

		if err := binary.Write(&b, binary.LittleEndian, v.HighByte); err != nil {
			return b.Bytes(), err
		}
		// export the number of LB
		if err := binary.Write(&b, binary.LittleEndian, uint16(len(v.LowBytes))); err != nil {
			return b.Bytes(), err
		}
		for _, v1 := range v.LowBytes {
			if err := binary.Write(&b, binary.LittleEndian, v1); err != nil {
				return b.Bytes(), err
			}
		}
	}
	return b.Bytes(), nil
}
