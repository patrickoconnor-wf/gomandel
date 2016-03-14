Mandelbrot in Golang
===========

![spiral](https://raw.githubusercontent.com/marijnfs/gomandel/master/spiral.jpg)

Calculates a mandelbrot in double point precision, normalises the values such that the palette is distributed equally.

## Run
Set x, y and zoom to an interesting point in the mandelbrot and run

`go run mandel.go -xres 1366 -yres 768 -x -.7454 -y 0.1242 -r .005 -aa 4`

This creates the image out.png with 4x anti aliasing. For more options see:

`go run mander.go --help`

`convert out.png -resize 1366x768 out2.png`

## Install
Depends on a resize library for resizing after anti aliasing. Install by running:

`go get github.com/nfnt/resize`

## Todo
  * Allow for deeper zooms

