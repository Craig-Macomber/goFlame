package affine

import (
	. "gomatrix.googlecode.com/hg/matrix"
)


type Affine struct {
	mat   *DenseMatrix
	shift *DenseMatrix
}

func NewAffine(size int) *Affine {
	return &Affine{Eye(size), Zeros(1, size)}
}

func MakeAffine(mat, shift *DenseMatrix) *Affine {
	return &Affine{mat, shift}
}

func (a *Affine) SetShift(mat *DenseMatrix) {
	x, _ := a.mat.GetSize()
	rows, cols := mat.GetSize()
	if x == rows && 1 == cols {
		a.shift = mat
	} else {
		panic("Invalid matrix size")
	}
}

func (a *Affine) GetShift() *DenseMatrix {
	return a.shift
}

func (a *Affine) SetMat(mat *DenseMatrix) {
	x, _ := a.mat.GetSize()
	rows, cols := mat.GetSize()
	if x == rows && x == cols {
		a.mat = mat
	} else {
		panic("Invalid matrix size")
	}
}

func (a *Affine) GetMat() (mat *DenseMatrix) {
	return a.mat
}

func (a *Affine) Transform(vec *DenseMatrix) *DenseMatrix {
	out, _ := a.mat.TimesDense(vec)
	out.AddDense(a.shift)
	return out
}

func (a *Affine) GetOrigin() *DenseMatrix {
	size, _ := a.mat.GetSize()
	i := Eye(size)
	i.SubtractDense(a.mat)
	i, _ = i.Inverse()
	out, _ := i.TimesDense(a.shift)
	return out
}

func FromOrigin(mat, origin *DenseMatrix) *Affine {
	size, _ := mat.GetSize()
	i := Eye(size)
	i.SubtractDense(mat)
	shift, _ := i.TimesDense(origin)
	return &Affine{mat, shift}
}

func (a *Affine) Inverse() *Affine {
	origin := a.GetOrigin()
	mat, _ := a.mat.Inverse()
	return FromOrigin(mat, origin)
}


func (a *Affine) Trans(x, y float64) (nx, ny float64) {
	loc := Zeros(2, 1)
	loc.Set(0, 0, x)
	loc.Set(1, 0, y)
	loc = a.Transform(loc)
	return loc.Get(0, 0), loc.Get(1, 0)
}

func FromOrigin2(mat *DenseMatrix, x, y float) *Affine {
	loc := Zeros(2, 1)
	loc.Set(0, 0, float64(x))
	loc.Set(1, 0, float64(y))
	return FromOrigin(mat, loc)
}
