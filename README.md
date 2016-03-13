Mandelbrot in Golang
===========

Calculates a mandelbrot in double point precision, normalises the values such that the palette is distributed equally.

## Run
Set x, y and zoom to an interesting point in the mandelbrot and run
> go run mandel.go
This creates the image out.png. To do antialiasing, generate the image at a higher resolution (e.g. in examples 1366*4 x 768*4) and use imagemagic to reduce resolution
> convert out.png -resize 1366x768 out2.png

## Todo
  * Use Flags to make it easily parametrisable
  * Allow for deeper zooms

