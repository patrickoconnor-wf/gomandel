Mandelbrot in Golang
===========

![spiral](https://raw.githubusercontent.com/marijnfs/gomandel/master/spiral.jpg)

Calculates a mandelbrot in double point precision, normalises the values such that the palette is distributed equally.

### Install
Depends on a resize library for resizing after anti aliasing. Install by running:

`go get github.com/nfnt/resize`


### Run
Set x, y and zoom to an interesting point in the mandelbrot and run

`go run mandel.go -xres 1366 -yres 768 -x -.7454 -y 0.1242 -r .005 -aa 4`

This creates the image out.png with 4x anti aliasing. For more options see:

`go run mandel.go --help`

### Concurrency

With help from _jasonmoo_ concurrency is now supported. To make a nice background using all cores, run:

`GOMAXPROCS=0 go run mandel.go  -xres 1920 -yres 1080 -aa 4`

## Todo
  * Allow for deeper zooms

## Examples
   * `./mandel --aa 4 -xres 1680 -yres 1050 -out bac.png -x  -0.745428 -y  0.113009 -r .00003 --palette alternate`
