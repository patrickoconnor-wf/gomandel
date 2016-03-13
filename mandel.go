package main

import (
	"os"
	"fmt"
	img "image"
	"image/color"
	"image/color/palette"
	"image/png"
	"math"
	"sort"
)

const IT = 512

func it(ca, cb float64) (int, float64) {
	var a, b float64 = 0, 0
	for i := 0; i < IT; i++ {
		as, bs := a*a, b*b
		if as + bs > 4 {
			return i, as + bs
		}
		//if as + bs < .00001 {
		//	return .00001
		//}
		a, b = as - bs + ca, 2 * a * b + cb
	}
	return IT, a * a + b * b
}

var Gameboy = []color.Color{
	color.RGBA{14, 55, 15, 255},
	color.RGBA{47, 97, 48, 255},
	color.RGBA{138, 171, 25, 255},
	color.RGBA{154, 187, 27, 255},
}

var Retro = []color.Color{
	color.RGBA{0x00, 0x04, 0x0f, 0xff},
	color.RGBA{0x03, 0x26, 0x28, 0xff},
	color.RGBA{0x07, 0x3e, 0x1e, 0xff},
	color.RGBA{0x18, 0x55, 0x08, 0xff},
	color.RGBA{0x5f, 0x6e, 0x0f, 0xff},
	color.RGBA{0x84, 0x50, 0x19, 0xff},
	color.RGBA{0x9b, 0x30, 0x22, 0xff},
	color.RGBA{0xb4, 0x92, 0x2f, 0xff},
	color.RGBA{0x94, 0xca, 0x3d, 0xff},
	color.RGBA{0x4f, 0xd5, 0x51, 0xff},
	color.RGBA{0x66, 0xff, 0xb3, 0xff},
	color.RGBA{0x82, 0xc9, 0xe5, 0xff},
	color.RGBA{0x9d, 0xa3, 0xeb, 0xff},
	color.RGBA{0xd7, 0xb5, 0xf3, 0xff},
	color.RGBA{0xfd, 0xd6, 0xf6, 0xff},
	color.RGBA{0xff, 0xf0, 0xf2, 0xff},
}

func main() {
	//width, height := 600, 600
	//width, height := 1366*4, 768*4
	width, height := 1680*4, 1050*4
	ratio := float64(height) / float64(width)
	//xpos, ypos, zoom_width := -.748, 0.1, .003
	//xpos, ypos, zoom_width := -.235125, .827214, 4.0e-5
	//xpos, ypos, zoom_width := -.16070135, 1.0375665, 1.0e-7
	//xpos, ypos, zoom_width := -.7453, .1127, 6.5e-4
	xpos, ypos, zoom_width := 0.45272105023, 0.396494224267,  .3E-9
	//xpos, ypos, zoom_width := -.160568374422, 1.037894847008, .000001
	//xpos, ypos, zoom_width := .232223859135, .559654166164, .00000000004
	xmin, xmax := xpos - zoom_width / 2.0, xpos + zoom_width / 2.0
	ymin, ymax := ypos - zoom_width * ratio / 2.0, ypos + zoom_width * ratio / 2.0
	
	
	single_values := make([]float64, width * height)
	
	fmt.Println("Mandelling...")

	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			a := (float64(x) / float64(width)) * (xmax - xmin) + xmin 
			b := (float64(y) / float64(height)) * (ymax - ymin) + ymin
			stop_it, norm := it(a, b)
			smooth_val := IT + 1 - (math.Log(norm) + float64(stop_it))
			smooth_val /= IT
			single_values[i] = 1.0 - smooth_val
			//fmt.Println(norm)
			//r, g, b := HuslToRGB(100. + 100. * smooth_val, 88.7, 44.3 + 20. * smooth_val)
			//r, g, b := smooth_val, .4 * smooth_val, -smooth_val
			//fmt.Println(norm, stop_it, smooth_val)
			//c := color.RGBA{uint8(255. * r), uint8(255. * g), uint8(255. * b), 255}
			
			//image.Set(x, y, c)
			i++
		}
	}

	sorted_values := make([]float64, len(single_values))
	for i := range sorted_values {
		sorted_values[i] = single_values[i]
	}
	sort.Float64s(sorted_values)
	//fmt.Println(sorted_values[0:10])

	var pal []color.Color
	if true {
		//pal = palette.Plan9
		pal = palette.WebSafe
	} else {
		pal = Retro
	}
	//pal := palette.WebSafe
	//pal := Gameboy
	//pal := Retro
	split_values := make([]float64, len(pal)-1)

	
	factor := .98
	start := .9
	for i := range split_values {
		//index := (i+1) * len(sorted_values) / len(pal)
		index := int(float64(len(sorted_values)-1) * (1.0 - start))
		fmt.Println(index, len(sorted_values))
		split_values[i] = sorted_values[index]
		start *= factor
	}
	sort.Float64s(split_values)
	//fmt.Println(split_values)
	

	image := img.NewRGBA(img.Rectangle{img.Point{0, 0}, img.Point{width, height}})
	i = 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			color_index := sort.Search(len(split_values), func(j int) bool {return single_values[i] < split_values[j]})
			//fmt.Println(color_index)
			image.Set(x, y, pal[color_index])
			i++
		}
	}
	out_file, _ := os.Create("out.png")
	png.Encode(out_file, image)
}
