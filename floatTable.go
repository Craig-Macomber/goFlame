package floatTable

type FloatTable struct {
    CellLength uint // number of float32 per cell
    width uint // number of cell per row
    Data []float32
}

func (t *FloatTable) GetCellStart(x,y uint) uint {
    return (x+y*t.width)*t.CellLength
}

func (t *FloatTable) GetCell(x,y uint) ([]float32) {
    start:=t.GetCellStart(x,y)
    return t.Data[start:start+t.CellLength]
}

func NewFloatTable(width, height, cellLength uint) *FloatTable{
    return &FloatTable{cellLength,width,make([]float32, cellLength*width*height)}
}

func (t *FloatTable) Fill(value float32) {
    lim:=len(t.Data)
    for i:=0; i<lim; i++ {
        t.Data[i]=value
    }
}