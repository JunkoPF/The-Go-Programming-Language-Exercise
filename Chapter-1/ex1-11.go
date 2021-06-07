package main

import (
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

func lissajous(out io.Writer, cycles float64, res float64, size int, nframes int, delay int) {
	var palette = []color.Color{
		color.White,
		color.Black,
	}

	const (
		whiteIndex = 0 // first color in palette
		blackIndex = 1 // next color in palette
	)

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	rq := r.URL.Query()

	var (
		cycles  float64 = 8     // number of complete x oscillator revolutions
		res     float64 = 0.001 // angular resolution
		size    int     = 100   // image canvas covers [-size..+size]
		nframes int     = 32    // number of animation frames
		delay   int     = 8     // delay between frames in 10ms units
	)
	for key, value := range rq {
		switch key {
		case "cycles":
			{
				cycles, _ = strconv.ParseFloat(value[0], 64)
			}
		case "res":
			{
				res, _ = strconv.ParseFloat(value[0], 64)
			}
		case "size":
			{
				size, _ = strconv.Atoi(value[0])
			}
		case "nframes":
			{
				nframes, _ = strconv.Atoi(value[0])
			}
		case "delay":
			{
				delay, _ = strconv.Atoi(value[0])
			}
		}
	}
	lissajous(w, cycles, res, size, nframes, delay)
}
