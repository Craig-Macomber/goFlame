package floatTable

type FloatTable struct {
	CellLength int // number of float32 per cell
	width      int // number of cell per row
	Data       []float32
}

func (t *FloatTable) GetCellStart(x, y int) int {
	return (x + y*t.width) * t.CellLength
}

func (t *FloatTable) GetCell(x, y int) []float32 {
	start := t.GetCellStart(x, y)
	return t.Data[start : start+t.CellLength]
}

func NewFloatTable(width, height, cellLength int) *FloatTable {
	return &FloatTable{cellLength, width, make([]float32, cellLength*width*height)}
}

func (t *FloatTable) Fill(value float32) {
	lim := len(t.Data)
	for i := 0; i < lim; i++ {
		t.Data[i] = value
	}
}

func (t *FloatTable) Rez() (int, int) {
	return t.width, len(t.Data) / t.CellLength / t.width
}
