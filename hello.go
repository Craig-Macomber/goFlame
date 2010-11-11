package main

import (
	//"runtime"
	"time"
	"/affine"
	. "gomatrix.googlecode.com/hg/matrix"
	"/floatTable"
	"math"
)

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
	v = math.Log2(v / math.MinFloat32)
	v = v * 255 / maxv
	return uint8(math.Fmax(0, math.Fmin(255, v)))
}


func main() {
	//runtime.GOMAXPROCS(1)


	const tableRez = 1000

	mat := Eye(2)
	mat.Scale(.58)
	mat.Set(1, 0, .2)

	const transCount = 3
	print("YYYYYY", "\n")
	trans := make([]*affine.Affine, transCount)
	print("xxxx", "\n")
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
	const channelCount = 3
	ft := floatTable.NewFloatTable(rezx, rezy, channelCount)
	ft2 := floatTable.NewFloatTable(rezx, rezy, channelCount)
	ft.Fill(math.MinFloat32)
	t := time.Nanoseconds()
	const iterCount = 30

	cfactor := make([]float32, channelCount)
	nyx := make([]int, rezy)
	nyy := make([]int, rezy)

	var niy1 int = iy1
	var niy2 int = iy2

	for i := 0; i < iterCount; i++ {
		print("Iter-")

		refit := i%5 == 0
		if refit {
			niy1 = iy2
			niy2 = iy1
		}

		ft2.Fill(0)
		for tn, t := range trans {
			for k := 0; k < channelCount; k++ {
				cfactor[k] = .5
			}
			cfactor[tn] = 4.0
			a11 := t.GetMat().Get(0, 0)
			a12 := t.GetMat().Get(0, 1)
			a21 := t.GetMat().Get(1, 0)
			a22 := t.GetMat().Get(1, 1)
			shiftx := t.GetShift().Get(0, 0)
			shifty := t.GetShift().Get(1, 0)
			for y := iy1; y <= iy2; y++ {
				nyx[y] = int(a12*float64(y)) << 4
				nyy[y] = int(a22*float64(y)) << 4
			}

			for x := ix1; x <= ix2; x++ {
				nnx := int(a11*float64(x)+shiftx) << 4
				nny := int(a21*float64(x)+shifty) << 4

				for y := iy1; y <= iy2; y++ {
					nx := (nnx + nyx[y]) >> 4
					ny := (nny + nyy[y]) >> 4
					out := ft2.GetCellStart(nx, ny)
					src := ft.GetCellStart(x, y)

					if refit {
						if niy1 > int(ny) {
							for k := 0; k < ft.CellLength; k++ {
								v := ft.Data[src+k]
								if v > 0 {
									niy1 = ny
								}
							}
						} else if niy2 < ny {
							for k := 0; k < ft.CellLength; k++ {
								v := ft.Data[src+k]
								if v > 0 {
									niy2 = ny
								}
							}
						}
					}

					for k := 0; k < ft.CellLength; k++ {
						v := ft.Data[src+k]
						ft2.Data[out+k] += v * cfactor[k]
					}
				}
			}
		}
		ft2, ft = ft, ft2
		iy1, iy2 = niy1, niy2
	}
	t = time.Nanoseconds() - t
	print("Time", "\n")
	print(t/1000000, "\n")

}
