// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 21.

// Server3 is an "echo" server that displays request parameters.
// Practice 1.12: Add http://localhost:8000/?cycles=20 for dinamic cycles
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

var palette = []color.Color{
	color.White,                        // 0: white background
	color.RGBA{0x00, 0xff, 0x00, 0xff}, // 1: bright green
	color.RGBA{0x00, 0xcc, 0x00, 0xff}, // 2: darker green
	color.RGBA{0x00, 0x99, 0x00, 0xff}, // 3: even darker green
	color.RGBA{0x00, 0x00, 0xff, 0xff}, // 4: bright blue
	color.RGBA{0x00, 0x00, 0xcc, 0xff}, // 5: darker blue
	color.RGBA{0x00, 0x00, 0x99, 0xff}, // 6: even darker blue
	color.RGBA{0x80, 0x00, 0x80, 0xff}, // 7: purple
	color.RGBA{0xff, 0xa5, 0x00, 0xff}, // 8: orange
	color.RGBA{0xff, 0x00, 0x00, 0xff}, // 9: red
}

func lissajous(out io.Writer, inputCycles int) error {
	const (
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	cycles := float64(inputCycles) // number of complete x oscillator revolutions
	freq := rand.Float64() * 3.0   // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference

	numDrawingColors := len(palette) - 1

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)

		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)

			// Select color index based on the value of 't'
			// Map t to a value between 1 and numDrawingColors
			// (t / (cycles * 2 * math.Pi)) gets the proportion of t in the whole cycle (0.0 to 1.0)
			// Multiply by numDrawingColors and add 1 to get the corresponding color index
			colorIndex := uint8(1 + int(t*freq)%numDrawingColors)

			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				colorIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	if err := gif.EncodeAll(out, &anim); err != nil {
		log.Printf("lissajous: encoding error: %v", err)
		return fmt.Errorf("lissajous: encoding error: %v", err)
	}

	return nil
}

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Parse the "cycles" query parameter, default to 5 if not provided
		cycles := 5
		w.Header().Set("Content-Type", "image/gif")
		if s := r.URL.Query().Get("cycles"); s != "" {
			if n, err := strconv.Atoi(s); err != nil {
				http.Error(w, fmt.Sprintf("Invalid cycles parameter: %q. Must be an integer.", s), http.StatusBadRequest)
				log.Printf("Invalid cycles parameter: %q (error: %v)", s, err)
				return
			} else if n <= 0 {
				http.Error(w, fmt.Sprintf("Invalid cycles parameter: %d. Must be a positive integer.", n), http.StatusBadRequest)
				log.Printf("Invalid cycles parameter: %d (must be > 0)", n)
				return
			} else {
				cycles = n
			}
		}

		err := lissajous(w, cycles)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating lissajous: %v", err), http.StatusInternalServerError)
			log.Printf("Error generating lissajous: %v", err)
			return
		}
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
