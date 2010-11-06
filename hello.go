package main

import (
    //"runtime"
    "time"
    "os"
    "image"
    "image/png"
    //"/bitTable"
    "/affine"
    . "gomatrix.googlecode.com/hg/matrix"
    "/floatTable"
    //"/lazyImage"
    "math"
    //"unsafe"
     )

func fit(trans []*affine.Affine) (x1, y1, x2, y2 float64){
    org:=trans[0].GetOrigin()
    x2=org.Get(0,0)
    x1=x2
    y2=org.Get(1,0)
    y1=y2
    nx1, ny1, nx2, ny2:=x1, y1, x2, y2
    done:=false
    for !done{
        //print (x1,x2,y1,y2,"\n")
        for i:=0; i<4; i++ {
            xx:=x1
            yy:=y1
            if i%2==0 {
                xx=x2
            }
            if i>1 {
                yy=y2
            }
            for _, t := range trans {
                x, y:=t.Trans(xx,yy)
                nx1=math.Fmin(nx1,x)
                nx2=math.Fmax(nx2,x)
                ny1=math.Fmin(ny1,y)
                ny2=math.Fmax(ny2,y)
            }
        }
        if nx1>x1 && nx2<x2 && ny1>y1 && ny2<y2 {
            done=true
        }
        x1, y1, x2, y2=nx1, ny1, nx2, ny2
        dx:=nx2-nx1
        const f=.001
        x1-=dx*f
        x2+=dx*f
        dy:=ny2-ny1
        y1-=dy*f
        y2+=dy*f
        
    }
    return
}

func col(v, maxv float64) uint8 {
	    v=math.Log2(v/math.MinFloat32)
        v=v*255/maxv
	    return uint8(math.Fmax(0,math.Fmin(255,v)))
}


func main() {
    //runtime.GOMAXPROCS(1)
    
	
	//print("Table", "\n")
	
	const tableRez=500
	
	mat:=Eye(2)
	mat.Scale(.58)
	mat.Set(1,0,.2)
	
	const transCount=3
	trans:=make([]*affine.Affine, transCount)
	trans[0]=affine.FromOrigin2(mat,0,0)
    trans[1]=affine.FromOrigin2(mat,.5,1)
    trans[2]=affine.FromOrigin2(mat,1,0)
	
    x1, y1, x2, y2:=fit(trans)
    //print (x1,x2,y1,y2,"\n")
    
    shift := Zeros(2, 1)
    shift.Set(0, 0, -x1)
    shift.Set(1, 0, -y1)
    scale:=(tableRez-4)/math.Fmax(x2-x1,y2-y1)
    shift.Scale(scale)
    shift.AddDense(Scaled(Ones(2,1),2))
    
    //print (scale," "+shift.String(),"\n")
    for i, t := range trans {
        origin:=Scaled(t.GetOrigin(),scale)
        origin.AddDense(shift)
        trans[i]=affine.FromOrigin(t.GetMat(),origin)
    }
	
	
	x1, y1, x2, y2=fit(trans)
	ix1:=int(x1)
	ix2:=int(x2)
	iy1:=int(y1)
	iy2:=int(y2)
	//print (int(x1)," ",int(x2)," ",int(y1)," ",int(y2),"\n")
	
	rezx:=ix2+2
	rezy:=iy2+2
	ft:=floatTable.NewFloatTable(uint(rezx),uint(rezy),3)
	ft2:=floatTable.NewFloatTable(uint(rezx),uint(rezy),3)
	ft.Fill(math.MinFloat32)
	t:=time.Nanoseconds()
	const iterCount=30
	for i:=0; i<iterCount ; i++ {
	    print ("Iter-")
	    ft2.Fill(0)
	    for tn, t := range trans {
	        utn:=uint(tn)
	        a11:=t.GetMat().Get(0,0)
	        a12:=t.GetMat().Get(0,1)
	        a21:=t.GetMat().Get(1,0)
	        a22:=t.GetMat().Get(1,1)
	        shiftx:=t.GetShift().Get(0,0)
	        shifty:=t.GetShift().Get(1,0)
            for x:=ix1; x<=ix2 ; x++ {
                nnx:=a11*float64(x)+shiftx
                nny:=a21*float64(x)+shifty
                for y:=iy1; y<=iy2 ; y++ {
                    nx:=nnx+a12*float64(y)
                    ny:=nny+a22*float64(y)
                    out:=ft2.GetCellStart(uint(nx),uint(ny))
                    src:=ft.GetCellStart(uint(x),uint(y))
                    //sum:=float32(0.0)
                    //for k := uint(0) ; k<ft.CellLength ; k++{
                    //    sum+=ft.Data[src+k]
                    //}
                    //same:=ft.Data[src+uint(tn)]
                    for k := uint(0) ; k<ft.CellLength ; k++{
//                         b:=k==utn
//                         ptr:=unsafe.Pointer(&b)
//                         v:=*ptr
//                         s:=float32(.5+float32(int(v))*3.5)
                        s:=float32(.5)
                        if k==utn {
                            s=4.0
                        }
                        ft2.Data[out+k]+=ft.Data[src+k]*s

                    }
                }
            }
        }
        ft2,ft=ft,ft2
	}
	t=time.Nanoseconds()-t
	print ("Time", "\n")
	print (t/1000000, "\n")
	
	//const rez=tableRez
	m:=image.NewNRGBA(rezx,rezy)
	maxv:=float64(0.0)
	for  _, t := range ft.Data {
	   maxv=math.Fmax(maxv,float64(t))
	}
	//print ("\n")
	//print ("mv=")
	//print (maxv,"\n")
	maxv=math.Log2(maxv/math.MinFloat32)
	//print (maxv,"\n")

    for x := 0; x<rezx ; x++ {
	    for y := 0; y<rezy ; y++ {
	        f:=ft.GetCell(uint(x),uint(y))
	        c:=new (image.NRGBAColor)
	        
	        v:=float64(f[0])
	        //print (v,"\n")
	        //v/=maxv
	        
	        v=math.Log2(v/math.MinFloat32)
	        v=v*255/maxv
	        //v=v*1
	        c.R=uint8(math.Fmax(0,math.Fmin(255,v)))
	        c.G=col(float64(f[1]),maxv)
	        c.B=col(float64(f[2]),maxv)
	        c.A=255
            m.Pix[y*m.Stride+x]=*c
	    }
	}
	
	
	
	iTrans:=make([]*affine.Affine, len(trans))
	for i, t := range trans {
        iTrans[i]=t.Inverse()
    }
	
	//print("Image", "\n")
	
	//print("Writing File", "\n")
	f,_:=os.Open("testFile.png",os.O_WRONLY,0666)
	//print(f,"\n")
	
	//print(e,"\n")
	
	png.Encode(f,m)
	
	
}
