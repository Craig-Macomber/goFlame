package bitTable

type BitTable struct {
	stride uint // number of int64s per row
	b      []uint64
}

func (t *BitTable) bitToWord(x, y uint) (word, bit uint) {
	return uint(x/64) + uint(y*t.stride), x % 64
}

func (t *BitTable) Bit(x, y uint) bool {
	w, b := t.bitToWord(x, y)
	return t.b[w]&(1<<b) != 0
}

func (t *BitTable) SetBit(x, y uint) {
	w, b := t.bitToWord(x, y)
	t.b[w] |= 1 << b
}

func (t *BitTable) ClearBit(x, y uint) {
	w, b := t.bitToWord(x, y)
	t.b[w] = t.b[w] &^ (1 << b)
}

func NewBitTable(size uint) *BitTable {
	if size%64 != 0 {
		panic("Invalid Size")
	}
	s2 := size / 64 * size
	return &BitTable{uint(size / 64), make([]uint64, s2)}
}
