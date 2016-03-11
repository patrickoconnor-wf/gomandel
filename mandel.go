package main

import (
	"os"
	"fmt"
	img "image"
	"image/color"
	"image/png"
	"math"
)

const IT = 512

func it(ca, cb float64) (int, float64) {
	var a, b float64 = 0, 0
	for i := 0; i < IT; i++ {
		as, bs := a*a, b*b
		if as + bs > 6 {
			return i, as + bs
		}
		//if as + bs < .00001 {
		//	return .00001
		//}
		a, b = as - bs + ca, 2 * a * b + cb
	}
	return IT, a * a + b * b
}

func main() {
	width, height := 4000, 2000
	ratio := float64(height) / float64(width)
	xpos, ypos, zoom_width := .275, .4775, .01
	xmin, xmax := xpos - zoom_width / 2.0, xpos + zoom_width / 2.0
	ymin, ymax := ypos - zoom_width * ratio / 2.0, ypos + zoom_width * ratio / 2.0
	

	fmt.Println("hello")
	image := img.NewRGBA(img.Rectangle{img.Point{0, 0}, img.Point{width, height}})
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			a := (float64(x) / float64(width)) * (xmax - xmin) + xmin 
			b := (float64(y) / float64(height)) * (ymax - ymin) + ymin
			stop_it, norm := it(a, b)
			smooth_val := IT + 1 - (math.Log(norm) + float64(stop_it))
			smooth_val /= IT
			//fmt.Println(norm)
			r, g, b := HuslToRGB(smooth_val * 3600., 100.0, smooth_val * 100.0)
			//fmt.Println(norm, stop_it, smooth_val)
			c := color.RGBA{uint8(255. * r), uint8(255. * g), uint8(255. * b), 255}
			
			image.Set(x, y, c)
		}
	}

	out_file, _ := os.Create("out.png")
	png.Encode(out_file, image)
}
