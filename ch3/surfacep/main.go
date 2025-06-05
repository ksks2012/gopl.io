// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
// Practice 3.1: Check return value of f(), ignore invalid polygon
// Practice 3.2: Render other functions (eggbox, moguls, saddle)
// Practice 3.3: Color polygons based on height (red for peaks, blue for valleys)
// Practice 3.4: Create a web server to render SVG surface data.

package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	width, height = 600, 320                     // canvas size in pixels
	cells         = 100                          // number of grid cells
	xyrange       = 30.0                         // axis ranges (-xyrange..+xyrange)
	xyscale       = float64(width) / 2 / xyrange // pixels per x or y unit (according to width)
	zscale        = float64(height) * 0.4        // pixels per z unit (according to height)
	angle         = math.Pi / 6                  // angle of x, y axes (=30°)

	minZ_est = -1.8
	maxZ_est = 1.8
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

func main() {
	http.HandleFunc("/", handleSurface)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleSurface(w http.ResponseWriter, r *http.Request) {
	// 1. Set Content-Type header
	w.Header().Set("Content-Type", "image/svg+xml")

	// 2. Get configuration from URL query parameters (optional, adds flexibility)
	// Example: /?width=800&height=400&type=saddle
	var currentWidth, currentHeight int = 600, 320            // default values
	var currentSurfaceFunc SurfaceFunc = fSinc                // default function
	var currentMinZ, currentMaxZ float64 = minZ_est, maxZ_est // default Z range

	if wStr := r.URL.Query().Get("width"); wStr != "" {
		if w, err := strconv.Atoi(wStr); err == nil && w > 0 {
			currentWidth = w
		}
	}
	if hStr := r.URL.Query().Get("height"); hStr != "" {
		if h, err := strconv.Atoi(hStr); err == nil && h > 0 {
			currentHeight = h
		}
	}

	surfaceType := r.URL.Query().Get("type")
	switch surfaceType {
	case "sinc":
		currentSurfaceFunc = fSinc
		currentMinZ, currentMaxZ = -0.2, 1.0
	case "eggbox":
		currentSurfaceFunc = fEggBox
		currentMinZ, currentMaxZ = -0.6, 0.6
	case "moguls":
		currentSurfaceFunc = fMoguls
		currentMinZ, currentMaxZ = -0.3, 0.3
	case "saddle":
		currentSurfaceFunc = fSaddle
		currentMinZ, currentMaxZ = -1.8, 1.8
	default:
		fmt.Fprintf(w, "Invalid surface type: %s\n", surfaceType)
		w.WriteHeader(http.StatusBadRequest)
		// return
	}

	xyscale := float64(currentWidth) / 2 / xyrange // pixels per x or y unit
	zscale := float64(currentHeight) * 0.4         // pixels per z unit

	// 3. Write SVG data to the ResponseWriter
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; stroke-width: 0.7' "+
		"width='%d' height='%d'>", currentWidth, currentHeight)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az, okA := corner(i+1, j, currentSurfaceFunc, float64(currentWidth), float64(currentHeight), xyscale, zscale)
			bx, by, bz, okB := corner(i, j, currentSurfaceFunc, float64(currentWidth), float64(currentHeight), xyscale, zscale)
			cx, cy, cz, okC := corner(i, j+1, currentSurfaceFunc, float64(currentWidth), float64(currentHeight), xyscale, zscale)
			dx, dy, dz, okD := corner(i+1, j+1, currentSurfaceFunc, float64(currentWidth), float64(currentHeight), xyscale, zscale)

			if !okA || !okB || !okC || !okD {
				continue
			}

			avgZ := (az + bz + cz + dz) / 4.0
			color := calculateColor(avgZ, currentMinZ, currentMaxZ)

			fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='#%s'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, color)
		}
	}
	fmt.Fprintln(w, "</svg>")
}
func calculateColor(z, minZ, maxZ float64) string {
	normalizedZ := (z - minZ) / (maxZ - minZ)
	if normalizedZ < 0 {
		normalizedZ = 0
	} else if normalizedZ > 1 {
		normalizedZ = 1
	}

	red := int(normalizedZ * 255)
	blue := int((1 - normalizedZ) * 255)

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

func corner(i, j int, f SurfaceFunc, canvasWidth, canvasHeight, xyscale, zscale float64) (float64, float64, float64, bool) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z := f(x, y)

	if math.IsNaN(z) || math.IsInf(z, 0) {
		return 0, 0, 0, false
	}

	sx := canvasWidth/2 + (x-y)*cos30*xyscale
	sy := canvasHeight/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z, true
}
