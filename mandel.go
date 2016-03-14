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
	"flag"

	"github.com/nfnt/resize"
)

var IT, xres, yres, aa int
var xpos, ypos, radius float64
var out_filename, palette_string string
var invert bool

func init() {
	flag.IntVar(&IT, "IT", 512, "maximum number of iterations")
	flag.IntVar(&xres, "xres", 500, "x resolution")
	flag.IntVar(&yres, "yres", 500, "y resolution")
	flag.IntVar(&aa, "aa", 1, "anti alias, e.g. set aa=4 for 4xAA")
	flag.Float64Var(&xpos, "x", -.75, "real coordinate")
	flag.Float64Var(&ypos, "y", 0.0, "imaginary coordinate")
	flag.Float64Var(&radius, "r", 3.0, "radius")
	flag.StringVar(&out_filename, "out", "out.png", "output file")
	flag.StringVar(&palette_string, "palette", "plan9", "One of: plan9|websafe|gameboy|retro")
	flag.BoolVar(&invert, "invert", false, "Inverts colouring")
	flag.Parse()
}

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
	//width, height := 1366, 768
	width, height := xres*aa, yres*aa
	ratio := float64(height) / float64(width)
	//xpos, ypos, zoom_width := -.748, 0.1, .003
	//xpos, ypos, zoom_width := -.235125, .827214, 4.0e-5
	//xpos, ypos, zoom_width := -.16070135, 1.0375665, 1.0e-7
	//xpos, ypos, zoom_width := -.7453, .1127, 6.5e-4
	//xpos, ypos, zoom_width := 0.45272105023, 0.396494224267,  5E-9
	//xpos, ypos, zoom_width := -.160568374422, 1.037894847008, .000001
	//xpos, ypos, zoom_width := .232223859135, .559654166164, .00000000004
	xmin, xmax := xpos - radius / 2.0, xpos + radius / 2.0
	ymin, ymax := ypos - radius * ratio / 2.0, ypos + radius * ratio / 2.0
	
	
	single_values := make([]float64, width * height)
	
	fmt.Println("Mandelling...")

	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			a := (float64(x) / float64(width)) * (xmax - xmin) + xmin 
			b := (float64(y) / float64(height)) * (ymax - ymin) + ymin
			stop_it, norm := it(a, b)
			smooth_val := float64(IT - stop_it) + math.Log(norm)
			//smooth_val /= IT
			if invert {
				single_values[i] = smooth_val
			} else {
				single_values[i] = -smooth_val
			}
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
	palette_map := make(map[string][]color.Color)
	palette_map["plan9"] = palette.Plan9
	palette_map["websafe"] = palette.WebSafe
	palette_map["gameboy"] = Gameboy
	palette_map["retro"] = Retro
	
	pal = palette_map[palette_string]
	//pal := palette.WebSafe
	//pal := Gameboy
	//pal := Retro
	split_values := make([]float64, len(pal)-1)
	for i := range split_values {
		split_values[i] = sorted_values[(i+1) * len(sorted_values) / len(pal)]
	}
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

	image_resized := resize.Resize(uint(xres), uint(yres), image, resize.Lanczos3)
	out_file, _ := os.Create(out_filename)
	png.Encode(out_file, image_resized)
}
