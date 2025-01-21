package basics

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Image struct {
	Width, Height int
}

func (img Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (img Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.Width, img.Height)
}

func (img Image) At(x, y int) color.Color {
	fx := float64(x) / float64(img.Width)
	fy := float64(y) / float64(img.Height)

	r := uint8(math.Sin(fx*10+math.Cos(fy*20))*127 + 128)
	g := uint8(math.Sin(fy*10+fx)*127 + 128)
	b := uint8(math.Cos(fx*20+fy*10)*127 + 128)

	return color.RGBA{r, g, b, 255}
}

func IGmain() {
	fmt.Printf("\n%s\n", ("-----------------image generator-----------------"))
	width := 510
	height := 520
	img := Image{Width: width, Height: height}

	f, err := os.Create("awesome_image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	png.Encode(f, img)
	println("Image saved to awesome_image.png")
}
