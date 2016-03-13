package main

import (
	"os"
	"fmt"
	img "image"
	//"image/color"
	"image/color/palette"
	"image/png"
	"math"
	"sort"
)

const IT = 2000

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

func main() {
	//width, height := 1366, 768
	width, height := 1366*4, 768*4
	ratio := float64(height) / float64(width)
	//xpos, ypos, zoom_width := -.748, 0.1, .003
	xpos, ypos, zoom_width := -.235125, .827214, 4.0e-5
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
			single_values[i] = smooth_val
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

	pal := palette.Plan9
	split_values := make([]float64, len(pal)-1)
	for i := range split_values {
		split_values[i] = sorted_values[i * len(sorted_values) / len(split_values)]
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
	out_file, _ := os.Create("out.png")
	png.Encode(out_file, image)
}
