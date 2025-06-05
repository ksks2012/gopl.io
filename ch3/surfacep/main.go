// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
// Practice 3.1: Check return value of f(), ignore invalid polygon
// Practice 3.2: Render other functions (eggbox, moguls, saddle)
// Practice 3.3: Color polygons based on height (red for peaks, blue for valleys)

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

// sinc function
func fSinc(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	if r == 0 {           // Avoid division by zero; math.Sin(0)/0 would result in NaN
		return 1.0 // As r approaches 0, sin(r)/r approaches 1
	}
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
	var actualMinZ, actualMaxZ float64 // Used to store the actual Z range; needed if using two-pass scan.
	// In single-pass mode, these values are replaced or ignored by estimated values.

	// Select the corresponding function based on the command-line parameter
	switch *surfaceType {
	case "sinc":
		currentSurfaceFunc = fSinc
		actualMinZ = -0.2 // The Z range of the sinc function is about -0.2 to 1.0
		actualMaxZ = 1.0
	case "eggbox":
		currentSurfaceFunc = fEggBox
		actualMinZ = -0.6 // The range of sin(x)+sin(y) is -2 to 2, multiplied by 0.3 is -0.6 to 0.6
		actualMaxZ = 0.6
	case "moguls":
		currentSurfaceFunc = fMoguls
		actualMinZ = -0.3 // Estimated
		actualMaxZ = 0.3
	case "saddle":
		currentSurfaceFunc = fSaddle
		actualMinZ = -1.8 // (x*x - y*y)/500, with xyrange=30, max is 30*30/500 = 1.8, min is -1.8
		actualMaxZ = 1.8
	default:
		fmt.Fprintf(os.Stderr, "invalid surface type: %s\nAvailable types: sinc, eggbox, moguls, saddle\n", *surfaceType)
		os.Exit(1)
	}

	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az, okA := corner(i+1, j, currentSurfaceFunc)
			bx, by, bz, okB := corner(i, j, currentSurfaceFunc)
			cx, cy, cz, okC := corner(i, j+1, currentSurfaceFunc)
			dx, dy, dz, okD := corner(i+1, j+1, currentSurfaceFunc)

			// Skip this polygon if any corner is invalid
			if !okA || !okB || !okC || !okD {
				continue
			}

			// Calculate the average z value of the polygon (or use one vertex's z, e.g., az)
			avgZ := (az + bz + cz + dz) / 4.0

			// Calculate color based on avgZ
			color := calculateColor(avgZ, actualMinZ, actualMaxZ)

			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='#%s'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, color)
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int, f SurfaceFunc) (float64, float64, float64, bool) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	if math.IsNaN(z) || math.IsInf(z, 0) {
		return 0, 0, 0, false
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z, true
}

// calculateColor function computes the RGB color (HEX string) based on the normalized z value.
// It maps z from the [minZ, maxZ] range to [0, 1], then maps it to a gradient from blue to red.
func calculateColor(z, minZ, maxZ float64) string {
	normalizedZ := (z - minZ) / (maxZ - minZ)
	if normalizedZ < 0 {
		normalizedZ = 0
	} else if normalizedZ > 1 {
		normalizedZ = 1
	}
	// Blue (0, 0, 255) corresponds to normalizedZ = 0 (minZ)
	// Red (255, 0, 0) corresponds to normalizedZ = 1 (maxZ)

	// From blue to red: blue component decreases from 255 to 0, red component increases from 0 to 255
	// Here we directly implement R: from 0 to 255, B: from 255 to 0
	// So when normalizedZ=0, blue=255, red=0
	// When normalizedZ=1, blue=0, red=255
	red := int(normalizedZ * 255)
	blue := int((1 - normalizedZ) * 255)

	// limit red and blue to [0, 255]
	if red < 0 {
		red = 0
	} else if red > 255 {
		red = 255
	}
	if blue < 0 {
		blue = 0
	} else if blue > 255 {
		blue = 255
	}

	return fmt.Sprintf("%02x%02x%02x", red, 0, blue)
}
