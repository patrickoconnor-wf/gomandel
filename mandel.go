package mandel

import (
	"fmt"
	img "image"
	"image/color"
	"image/color/palette"
	"math"
	"sort"
	"sync"

	"github.com/nfnt/resize"
)

var Gray []color.Color
var Gameboy []color.Color
var Retro []color.Color
var Alternate []color.Color
var BlackWhite []color.Color

func Init() {

	// flag.IntVar(&IT, "IT", 512, "maximum number of iterations")
	// flag.IntVar(&xres, "xres", 500, "x resolution")
	// flag.IntVar(&yres, "yres", 500, "y resolution")
	// flag.IntVar(&aa, "aa", 1, "anti alias, e.g. set aa=4 for 4xAA")
	// flag.Float64Var(&xpos, "x", -.75, "real coordinate")
	// flag.Float64Var(&ypos, "y", 0.0, "imaginary coordinate")
	// flag.Float64Var(&radius, "r", 3.0, "radius")
	// flag.StringVar(&out_filename, "out", "out.png", "output file")
	// flag.StringVar(&palette_string, "palette", "plan9", "One of: plan9|websafe|gameboy|retro|alternate")
	// flag.StringVar(&focusstring, "focus", "", "sequence of focus command. Select quadrant (numbered 1-4). e.g.: 1423. Read code to understand")
	// flag.BoolVar(&invert, "invert", false, "Inverts colouring")
	// flag.Parse()

	Gray = make([]color.Color, 255*3)
	for i := 0; i < 255*3; i++ {
		Gray[i] = color.RGBA{uint8(i / 3), uint8((i + 1) / 3), uint8((i + 2) / 3), 255}
	}

	Alternate = make([]color.Color, 20)
	for i := 0; i < len(Alternate); i++ {
		switch i % 6 {
		case 0:
			Alternate[i] = color.RGBA{0x18, 0x4d, 0x68, 255}
		case 1:
			Alternate[i] = color.RGBA{0x31, 0x80, 0x9f, 255}
		case 2:
			Alternate[i] = color.RGBA{0xfb, 0x9c, 0x6c, 255}
		case 3:
			Alternate[i] = color.RGBA{0xd5, 0x51, 0x21, 255}
		case 4:
			Alternate[i] = color.RGBA{0xcf, 0xe9, 0x90, 255}
		case 5:
			Alternate[i] = color.RGBA{0xea, 0xfb, 0xc5, 255}
		}
	}

	BlackWhite = make([]color.Color, 0)
	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			BlackWhite = append(BlackWhite, color.RGBA{0, 0, 0, 255})
		} else {
			BlackWhite = append(BlackWhite, color.RGBA{255, 255, 255, 255})
		}
	}
	Gameboy = []color.Color{
		color.RGBA{14, 55, 15, 255},
		color.RGBA{47, 97, 48, 255},
		color.RGBA{138, 171, 25, 255},
		color.RGBA{154, 187, 27, 255},
	}

	Retro = []color.Color{
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
}

func it(ca, cb float64, iter int) (int, float64) {
	var a, b float64 = 0, 0
	for i := 0; i < iter; i++ {
		as, bs := a*a, b*b
		if as+bs > 4 {
			return i, as + bs
		}
		//if as + bs < .00001 {
		//	return .00001
		//}
		a, b = as-bs+ca, 2*a*b+cb
	}
	return iter, a*a + b*b
}

func Create(xres int,
	yres int,
	aa int,
	iter int,
	xpos float64,
	ypos float64,
	radius float64,
	paletteString string,
	invert bool,
	focusString string,
	filename string) img.Image {

	var focusstring string
	width, height := xres*aa, yres*aa
	ratio := float64(height) / float64(width)
	//xpos, ypos, zoom_width := -.16070135, 1.0375665, 1.0e-7
	//xpos, ypos, zoom_width := -.7453, .1127, 6.5e-4
	//xpos, ypos, zoom_width := 0.45272105023, 0.396494224267,  .3E-9
	//xpos, ypos, zoom_width := -.160568374422, 1.037894847008, .000001
	//xpos, ypos, zoom_width := .232223859135, .559654166164, .00000000004
	yRadius := float64(radius * ratio)

	tempRadius, tempYRadius := radius/4.0, yRadius/4.0
	for _, command := range focusstring {
		switch string(command) {
		case "1":
			xpos -= tempRadius
			ypos += tempRadius
		case "2":
			xpos += tempRadius
			ypos += tempRadius
		case "3":
			xpos -= tempRadius
			ypos -= tempRadius
		case "4":
			xpos += tempRadius
			ypos -= tempRadius
		case "w":
			ypos += tempRadius
		case "s":
			ypos -= tempRadius
		case "a":
			xpos -= tempRadius
		case "d":
			xpos += tempRadius
		case "r":
			tempRadius, tempYRadius = radius/4, yRadius/4
		case "z":
			radius /= 2
			yRadius /= 2
			tempRadius, tempYRadius = radius/4, yRadius/4
		default:
			return nil
		}
		tempRadius /= 2
		tempYRadius /= 2
	}

	xmin, xmax := xpos-radius/2.0, xpos+radius/2.0
	ymin, ymax := ypos-yRadius/2.0, ypos+yRadius/2.0

	singleValues := make([]float64, width*height)

	fmt.Print("Mandelling...")

	var wg sync.WaitGroup

	for y := 0; y < height; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := 0; x < width; x++ {
				a := (float64(x)/float64(width))*(xmax-xmin) + xmin
				b := (float64(y)/float64(height))*(ymin-ymax) + ymax
				stopIt, norm := it(a, b, iter)
				smoothVal := float64(iter-stopIt) + math.Log(norm)
				i := y*width + x
				if invert {
					singleValues[i] = smoothVal
				} else {
					singleValues[i] = -smoothVal
				}
			}
		}(y)
	}
	wg.Wait()
	fmt.Println("Done")
	fmt.Print("Sorting...")
	sortedValues := make([]float64, len(singleValues))
	for i := range sortedValues {
		sortedValues[i] = singleValues[i]
	}
	sort.Float64s(sortedValues)

	fmt.Println("Done")

	cont := make([]color.Color, 10000)
	for i := range cont {
		//val := float64(i) / float64(len(cont))
		val := i * 256 / len(cont)
		cont[i] = color.RGBA{uint8(val), 0, uint8(255 - val), uint8(255)}
	}

	var pal []color.Color
	paletteMap := make(map[string][]color.Color)
	paletteMap["plan9"] = palette.Plan9
	paletteMap["websafe"] = palette.WebSafe
	paletteMap["gameboy"] = Gameboy
	paletteMap["retro"] = Retro
	paletteMap["gray"] = Gray
	paletteMap["cont"] = cont
	paletteMap["alternate"] = Alternate
	paletteMap["blackwhite"] = BlackWhite

	pal = paletteMap[paletteString]

	splitValues := make([]float64, len(pal)-1)

	// factor := .98
	// start := .9
	for i := range splitValues {
		index := (i + 1) * len(sortedValues) / len(pal)
		//index := int(float64(len(sortedValues)-1) * (1.0 - start))
		splitValues[i] = sortedValues[index]
		// start *= factor
	}
	sort.Float64s(splitValues)

	image := img.NewRGBA(img.Rectangle{img.Point{0, 0}, img.Point{width, height}})

	fmt.Print("Filling...")

	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			splitValues := sort.Search(len(splitValues), func(j int) bool { return singleValues[i] < splitValues[j] })
			image.Set(x, y, pal[splitValues])
			i++
		}
	}
	fmt.Println("Done")

	fmt.Println("Resizing...")
	imageResized := resize.Resize(uint(xres), uint(yres), image, resize.Lanczos3)
	fmt.Println("Done")
	return imageResized

	// outFile, _ := os.Create(filename)
	// png.Encode(outFile, imageResized)
	// fmt.Println("Finished writing to:", filename)
	// fmt.Printf("--r %v --x %v --y %v\n", radius, (xmin+xmax)/2, (ymin+ymax)/2)
}
