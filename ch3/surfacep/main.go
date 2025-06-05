// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
// Practice 3.1: Check return value of f(), ignore invalid polygon
// Practice 3.2: Add command-line flag to choose surface type
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

type SurfaceFunc func(x, y float64) float64

func fSinc(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

// Egg Box pattern
func fEggBox(x, y float64) float64 {
	return 0.3 * (math.Sin(x) + math.Sin(y))
}

// Moguls pattern (a bit more complex, can be varied)
func fMoguls(x, y float64) float64 {
	return 0.1*(math.Sin(x*0.5)*math.Cos(y*0.5)+math.Sin(x*0.8)*math.Cos(y*0.8)) +
		0.05*math.Sin(x*2)
}

// Saddle pattern
func fSaddle(x, y float64) float64 {
	return (x*x - y*y) / 500
}

var surfaceType = flag.String("type", "sinc", "Choose surface type: sinc, eggbox, moguls, saddle")

func main() {
	flag.Parse()

	var currentSurfaceFunc SurfaceFunc

	switch *surfaceType {
	case "sinc":
		currentSurfaceFunc = fSinc
	case "eggbox":
		currentSurfaceFunc = fEggBox
	case "moguls":
		currentSurfaceFunc = fMoguls
	case "saddle":
		currentSurfaceFunc = fSaddle
	default:
		fmt.Fprintf(
			os.Stderr, "invalid surface type: %s\nAvailable types: sinc, eggbox, moguls, saddle\n", *surfaceType)
		os.Exit(1)
	}

	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, errA := corner(i+1, j, currentSurfaceFunc)
			bx, by, errB := corner(i, j, currentSurfaceFunc)
			cx, cy, errC := corner(i, j+1, currentSurfaceFunc)
			dx, dy, errD := corner(i+1, j+1, currentSurfaceFunc)

			// If any point is invalid, skip this polygon
			if errA != nil || errB != nil || errC != nil || errD != nil {
				continue
			}
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int, f SurfaceFunc) (float64, float64, error) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)
	// Check if z is NaN or Inf, and return 0,0 if so.
	if math.IsNaN(z) || math.IsInf(z, 0) {
		return 0, 0, fmt.Errorf("invalid z value: %f at (%g, %g)", z, x, y)
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, nil
}
