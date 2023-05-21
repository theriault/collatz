package cmd

import (
	"image/png"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// shared options
var fn string
var power int
var minX float64
var minY float64
var maxX float64
var maxY float64

type Info struct {
	Title string
	File  string
}

// Types are the main stopping time functions we will examine
var Types = map[string]Info{
	"": {
		Title: "Combined",
		File:  "combined",
	},
	"f": {
		Title: "f(x)",
		File:  "f",
	},
	"g": {
		Title: "g(x)",
		File:  "g",
	},
	"h": {
		Title: "h(x)",
		File:  "h",
	},
}

// saveToPNG is a helper function to save a given plot to the filesystem
func saveToPNG(fileName string, width, height int, p *plot.Plot) error {
	fullPath := "results/" + fileName
	log.Printf("writing to %s...\n", fullPath)
	img := vgimg.New(vg.Points(float64(width)), vg.Points(float64(height)))
	dc := draw.New(img)
	p.Draw(dc)
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, img.Image()); err != nil {
		return err
	}
	return nil
}

// newPlot is a helper function to instantiate a plot with font sizes already scaled up
func newPlot() *plot.Plot {
	p := plot.New()
	p.X.Tick.Label.Font.Size = 20
	p.Y.Tick.Label.Font.Size = 20
	p.Title.TextStyle.Font.Size = 30
	p.X.Label.TextStyle.Font.Size = 25
	p.Y.Label.TextStyle.Font.Size = 25
	p.Legend.TextStyle.Font.Size = 25
	p.Legend.Top = true
	return p
}

// applyConstraintsToPlot is a helper function set the min/max for the plot
func applyConstraintsToPlot(p *plot.Plot, minX, maxX, minY, maxY float64) {
	if minX != 0 {
		p.X.Min = minX
	}
	if minY != 0 {
		p.Y.Min = minY
	}
	if maxX != 0 {
		p.X.Max = maxX
	}
	if maxY != 0 {
		p.Y.Max = maxY
	}
}

func max[T float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}
