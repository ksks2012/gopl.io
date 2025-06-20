// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
// Practice 8.5: Modify the surface program to use goroutines to compute the polygons concurrently.
package main

import (
	"math"
	"runtime"
	"sync"
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

func main() {
	basic()
	concurrent()
}

func saveAsSVG(fileName string, polygons [][8]float64) {
	// This function would save the polygons to an SVG file.
	// Implementation is not shown here, as it is not part of the exercise.
}

func concurrent() {
	polygons := make([][8]float64, 0, cells*cells)
	polygonCh := make(chan [8]float64, cells*cells)

	numWorkers := runtime.NumCPU()
	if numWorkers == 0 {
		numWorkers = 1
	}

	var wg sync.WaitGroup // Use sync.WaitGroup to wait for all worker goroutines to finish

	// Distribute the cells*cells tasks among numWorkers goroutines
	// Each worker processes a portion of the rows

	rowsPerWorker := cells / numWorkers
	if cells%numWorkers != 0 {
		rowsPerWorker++
	}

	for k := 0; k < numWorkers; k++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			startRow := workerID * rowsPerWorker
			endRow := startRow + rowsPerWorker
			if endRow > cells {
				endRow = cells
			}

			for i := startRow; i < endRow; i++ {
				for j := 0; j < cells; j++ {
					ax, ay := corner(i+1, j)
					bx, by := corner(i, j)
					cx, cy := corner(i, j+1)
					dx, dy := corner(i+1, j+1)
					polygonCh <- [8]float64{ax, ay, bx, by, cx, cy, dx, dy}
				}
			}
		}(k)
	}

	wg.Wait()
	close(polygonCh)

	for p := range polygonCh {
		polygons = append(polygons, p)
	}

	saveAsSVG("concurrent.svg", polygons)
}

func basic() {
	polygons := make([][8]float64, 0, cells*cells)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			polygons = append(polygons, [8]float64{ax, ay, bx, by, cx, cy, dx, dy})
		}
	}
	saveAsSVG("basic.svg", polygons)
}

func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

//!-
