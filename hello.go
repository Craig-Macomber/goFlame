package main

import (
	//"runtime"
	"time"
	"os"
	"image"
	"image/png"
	"/affine"
	. "gomatrix.googlecode.com/hg/matrix"
	"/floatTable"
	"math"
	"io"
)

const fill = 1.0 / (1 << 126) // min float32 that does not have performance overheads
const channelCount = 3
const iterCount = 10


func fit(trans []*affine.Affine) (x1, y1, x2, y2 float64) {
	org := trans[0].GetOrigin()
	x2 = org.Get(0, 0)
	x1 = x2
	y2 = org.Get(1, 0)
	y1 = y2
	nx1, ny1, nx2, ny2 := x1, y1, x2, y2
	done := false
	for !done {
		//print (x1,x2,y1,y2,"\n")
		for i := 0; i < 4; i++ {
			xx := x1
			yy := y1
			if i%2 == 0 {
				xx = x2
			}
			if i > 1 {
				yy = y2
			}
			for _, t := range trans {
				x, y := t.Trans(xx, yy)
				nx1 = math.Fmin(nx1, x)
				nx2 = math.Fmax(nx2, x)
				ny1 = math.Fmin(ny1, y)
				ny2 = math.Fmax(ny2, y)
			}
		}
		if nx1 > x1 && nx2 < x2 && ny1 > y1 && ny2 < y2 {
			done = true
		}
		x1, y1, x2, y2 = nx1, ny1, nx2, ny2
		dx := nx2 - nx1
		const f = .001
		x1 -= dx * f
		x2 += dx * f
		dy := ny2 - ny1
		y1 -= dy * f
		y2 += dy * f

	}
	return
}

func col(v, maxv float64) uint8 {
	v = math.Log2(v / fill)
	v = v * 255 / maxv
	return uint8(math.Fmax(0, math.Fmin(255, v)))
}


func MakeTransformer(t *affine.Affine, rezx, rezy int) func(x, y int) (xout, yout int) {
	a11 := t.GetMat().Get(0, 0)
	a12 := t.GetMat().Get(0, 1)
	a21 := t.GetMat().Get(1, 0)
	a22 := t.GetMat().Get(1, 1)
	shiftx := t.GetShift().Get(0, 0)
	shifty := t.GetShift().Get(1, 0)

	nyx := make([]int, rezy)
	nyy := make([]int, rezy)
	for y := 0; y < rezy; y++ {
		nyx[y] = int(a12*float64(y)) << 4
		nyy[y] = int(a22*float64(y)) << 4
	}

	nxx := make([]int, rezx)
	nxy := make([]int, rezx)

	for x := 0; x < rezx; x++ {
		nxx[x] = int(a11*float64(x)+shiftx) << 4
		nxy[x] = int(a21*float64(x)+shifty) << 4
	}

	return func(x, y int) (xout, yout int) {
		return (nxx[x] + nyx[y]) >> 4, (nxy[x] + nyy[y]) >> 4
	}

}


func Render(xmin, xmax, ymin, ymax int, ft, ft2 *floatTable.FloatTable, trans []*affine.Affine) {
	cfactor := make([]float32, channelCount)

	xform := make([]func(x, y int) (xout, yout int), len(trans))
	for tn, t := range trans {
		xform[tn] = MakeTransformer(t, xmax+1, ymax+1)
	}

	for i := 0; i < iterCount; i++ {
		print("Iter-")
		ft2.Fill(0)
		for tn, t := range xform {
			for k := 0; k < channelCount; k++ {
				cfactor[k] = .5
			}
			cfactor[tn] = 4.0
			for x := xmin; x <= xmax; x++ {
				for y := ymin; y <= ymax; y++ {
					out := ft2.GetCellStart(t(x, y))
					src := ft.GetCellStart(x, y)
					for k := 0; k < ft.CellLength; k++ {
						ft2.Data[out+k] += ft.Data[src+k] * cfactor[k]
					}
				}
			}
		}
		ft2, ft = ft, ft2
	}
}

func MakeImage(out io.Writer, ft *floatTable.FloatTable, colorFunc func(cell []float32) image.NRGBAColor) {
	rezx, rezy := ft.Rez()
	m := image.NewNRGBA(rezx, rezy)
	for x := 0; x < rezx; x++ {
		for y := 0; y < rezy; y++ {
			m.Pix[y*m.Stride+x] = colorFunc(ft.GetCell(x, y))
		}
	}
	png.Encode(out, m)
}

func MakeColorizer(ft *floatTable.FloatTable) func(cell []float32) image.NRGBAColor {
	maxv := float64(0.0)
	for _, t := range ft.Data {
		maxv = math.Fmax(maxv, float64(t))
	}

	println(maxv)
	maxv = math.Log2(maxv / fill)
	return func(cell []float32) image.NRGBAColor {
		c := new(image.NRGBAColor)
		c.R = col(float64(cell[0]), maxv)
		c.G = col(float64(cell[1]), maxv)
		c.B = col(float64(cell[2]), maxv)
		c.A = 255
		return *c
	}
}

// func Render(xmin, xmax, ymin, ymax int, ft, ft2 *floatTable.FloatTable, trans []*affine.Affine) {
//     cfactor:=make([]float32, channelCount)
// 	nyx:=make([]int, ymax+1)
// 	nyy:=make([]int, ymax+1)
//     
//     for i := 0; i < iterCount; i++ {
// 		print("Iter-")
// 		tt := time.Nanoseconds()
// 		ft2.Fill(0)
// 		ttt := time.Nanoseconds()
// 		for tn, t := range trans {
// 			for k := 0; k < channelCount; k++ { 
// 			    cfactor[k]=.5
// 			}
// 			cfactor[tn]=4.0
// 			a11 := t.GetMat().Get(0, 0)
// 			a12 := t.GetMat().Get(0, 1)
// 			a21 := t.GetMat().Get(1, 0)
// 			a22 := t.GetMat().Get(1, 1)
// 			shiftx := t.GetShift().Get(0, 0)
// 			shifty := t.GetShift().Get(1, 0)
// 			for y := ymin; y <= ymax; y++ {
// 			    nyx[y]=int(a12*float64(y))<<4
// 			    nyy[y]=int(a22*float64(y))<<4
// 			}
// 			
// 			for x := xmin; x <= xmax; x++ {
// 				nnx := int(a11*float64(x) + shiftx)<<4
// 				nny := int(a21*float64(x) + shifty)<<4
// 				
// 				
// 				
// 				for y := ymin; y <= ymax; y++ {
// 					nx := (nnx + nyx[y])>>4
// 					ny := (nny + nyy[y])>>4
// 					out := ft2.GetCellStart(nx, ny)
// 					src := ft.GetCellStart(x, y)
// 					for k := 0; k < ft.CellLength; k++ {
// 						ft2.Data[out+k] += ft.Data[src+k] * cfactor[k]
// 					}
// 				}
// 			}
// 		}
// 		tttt := time.Nanoseconds()
// 		print("\n", (ttt-tt)/1000000, " fill ","\n")
// 		print((tttt-ttt)/1000000, " frac ", "\n")
// 		ft2, ft = ft, ft2
// 	}
// }

func main() {
	//runtime.GOMAXPROCS(1)

	const tableRez = 1000

	mat := Eye(2)
	mat.Scale(.58)
	mat.Set(1, 0, .2)

	const transCount = 3
	trans := make([]*affine.Affine, transCount)
	trans[0] = affine.FromOrigin2(mat, 0, 0)
	trans[1] = affine.FromOrigin2(mat, .5, 1)
	trans[2] = affine.FromOrigin2(mat, 1, 0)

	x1, y1, x2, y2 := fit(trans)
	//print (x1,x2,y1,y2,"\n")

	shift := Zeros(2, 1)
	shift.Set(0, 0, -x1)
	shift.Set(1, 0, -y1)
	scale := (tableRez - 4) / math.Fmax(x2-x1, y2-y1)
	shift.Scale(scale)
	shift.AddDense(Scaled(Ones(2, 1), 2))

	//print (scale," "+shift.String(),"\n")
	for i, t := range trans {
		origin := Scaled(t.GetOrigin(), scale)
		origin.AddDense(shift)
		trans[i] = affine.FromOrigin(t.GetMat(), origin)
	}

	x1, y1, x2, y2 = fit(trans)
	ix1 := int(x1)
	ix2 := int(x2)
	iy1 := int(y1)
	iy2 := int(y2)
	//print (int(x1)," ",int(x2)," ",int(y1)," ",int(y2),"\n")

	rezx := ix2 + 2
	rezy := iy2 + 2

	ft := floatTable.NewFloatTable(rezx, rezy, channelCount)
	ft2 := floatTable.NewFloatTable(rezx, rezy, channelCount)

	ft.Fill(fill)

	t := time.Nanoseconds()
	Render(ix1, ix2, iy1, iy2, ft, ft2, trans)

	t = time.Nanoseconds() - t
	print("Time", "\n")
	print(t/1000000, "\n")

	println("Saving image")
	f, err := os.Open("testFile.png", os.O_WRONLY|os.O_CREAT, 0666)
	println(err)
	MakeImage(f, ft, MakeColorizer(ft))

}
