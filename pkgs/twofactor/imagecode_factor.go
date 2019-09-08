package twofactor

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/zouxinjiang/le/pkgs/lib"
	"golang.org/x/image/font/gofont/goregular"
	"math/rand"
	"os"
	"time"
)

type ImageFactor struct {
}

func init() {
	register("imagecode", NewImageCodeFactor)
}

func NewImageCodeFactor() TwoFactor {
	return &ImageFactor{}
}

func (i ImageFactor) Do(params map[string]string) (addr string, code string, err error) {
	code = params["code"]
	fpath := os.TempDir() + "/" + lib.RandStr(10) + ".png"
	err = i.generate(code, fpath)
	return fpath, code, err
}

func (self ImageFactor) generate(code, fpath string) error {
	var heigth, width = 40, len(code) * 25
	dc := gg.NewContext(width, heigth)
	dc.SetRGBA(1, 1, 1, 1)
	dc.Clear()
	tf, _ := truetype.Parse(goregular.TTF)
	face := truetype.NewFace(tf, &truetype.Options{Size: 32})
	dc.SetFontFace(face)
	rand.Seed(time.Now().UnixNano())
	for i, v := range []rune(code) {
		r := rand.Float64()
		g := rand.Float64() * r
		b := rand.Float64() * r * 2
		a := rand.Float64()*0.5 + 0.5
		dc.SetRGBA(r, g, b, a)
		dc.DrawString(string(v), float64(10+i*24), float64(heigth-32)/2+32)
	}
	return dc.SavePNG(fpath)
}
