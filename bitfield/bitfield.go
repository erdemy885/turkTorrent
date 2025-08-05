package bitfield

type Bitfield []byte

func (bits Bitfield) Get(index int) bool {
	byteIndex := index / 8
	bitIndex := index % 8
	return bits[byteIndex]>>(7-bitIndex)&1 != 0
}

func (bits Bitfield) Set(index int) {
	byteIndex := index / 8
	bitIndex := index % 8
	bits[byteIndex] |= 1 << (7 - bitIndex)
}
