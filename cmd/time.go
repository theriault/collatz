package cmd

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"runtime"
	"sync"

	"github.com/spf13/cobra"
	"github.com/theriault/collatz/shared"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

var (
	timeCmd = &cobra.Command{
		Use:   "time",
		Short: "Generate a scatter plot or histogram for the total stopping time of x",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, ok := Types[fn]
			if !ok {
				return fmt.Errorf("invalid value for --fn: %s", fn)
			}
			if power < 1 || power > 20 {
				return fmt.Errorf("--k must be in the range 1..20: %d", power)
			}
			limit := uint64(math.Pow10(power))
			graphType, err := cmd.Flags().GetString("graph")
			if err != nil {
				return err
			}

			p := newPlot()
			p.Title.Text = fmt.Sprintf("%s Total Stopping Time 10^%d", info.Title, power)
			if graphType == "histogram" {
				p.X.Label.Text = "Stopping Time"
				p.Y.Label.Text = "Count"
			} else if graphType == "scatter" {
				p.X.Label.Text = "x"
				p.Y.Label.Text = "Total Stopping Time"
			}
			p.Y.Min = 0
			p.X.Min = 1
			if graphType == "histogram" {
				if fn == "" || fn == "f" {
					log.Printf("building f(x) histogram for 10^%d...", power)
					buildHistogram(p, color.NRGBA{R: 255, G: 0, B: 0, A: 128}, limit, shared.CollatzStoppingTimeF, "f(x)")
				}
				if fn == "" || fn == "g" {
					log.Printf("building g(x) histogram for 10^%d...", power)
					buildHistogram(p, color.NRGBA{R: 0, G: 255, B: 0, A: 128}, limit, shared.CollatzStoppingTimeG, "g(x)")
				}
				if fn == "" || fn == "h" {
					log.Printf("building h(x) histogram for 10^%d...", power)
					buildHistogram(p, color.NRGBA{R: 0, G: 0, B: 255, A: 128}, limit, shared.CollatzStoppingTimeH, "h(x)")
				}
			} else if graphType == "scatter" {
				if fn == "" || fn == "f" {
					log.Printf("building f(x) scatter for 10^%d...", power)
					buildTime(p, color.NRGBA{R: 255, G: 0, B: 0, A: 128}, limit, shared.CollatzStoppingTimeF, "f(x)")
				}
				if fn == "" || fn == "g" {
					log.Printf("building g(x) scatter for 10^%d...", power)
					buildTime(p, color.NRGBA{R: 0, G: 255, B: 0, A: 128}, limit, shared.CollatzStoppingTimeG, "g(x)")
				}
				if fn == "" || fn == "h" {
					log.Printf("building h(x) scatter for 10^%d...", power)
					buildTime(p, color.NRGBA{R: 0, G: 0, B: 255, A: 128}, limit, shared.CollatzStoppingTimeH, "h(x)")
				}
			} else {
				return fmt.Errorf("unexpected value for --graph: %s", graphType)
			}
			applyConstraintsToPlot(p, minX, minY, maxX, maxY)
			fileName := fmt.Sprintf("time_%s_%s_%d.png", graphType, info.File, power)
			return saveToPNG(fileName, 1500, 900, p)
		},
	}
)

func init() {
	timeCmd.Flags().StringVar(&fn, "fn", "", "which function to plot: f (standard), g (reduced), h (main result). leave blank to plot all")
	timeCmd.Flags().IntVar(&power, "k", 5, "examine n up to 10^k")
	timeCmd.Flags().String("graph", "", "graph type: scatter | histogram")
	timeCmd.Flags().Float64Var(&minX, "min-x", 0, "min x to show on plot. use 0 for min of data")
	timeCmd.Flags().Float64Var(&minY, "min-y", 0, "min y to show on plot. use 0 for min of data")
	timeCmd.Flags().Float64Var(&maxX, "max-x", 0, "max x to show on plot. use 0 for max of data")
	timeCmd.Flags().Float64Var(&maxY, "max-y", 0, "max y to show on plot. use 0 for max of data")
}

func buildTime(p *plot.Plot, fill color.NRGBA, limit uint64, fn func(n uint64) (uint64, uint64, uint64), title string) {
	var wg sync.WaitGroup
	workers := uint64(runtime.GOMAXPROCS(0))
	completed := make(chan plotter.XYs, workers)
	for w := uint64(1); w <= workers; w++ {
		wg.Add(1)
		go (func(worker uint64, workerCount uint64, limit uint64, completed chan<- plotter.XYs) {
			defer wg.Done()
			xys := make(plotter.XYs, limit/workerCount+1)
			j := 0
			for i := worker; i < limit; i += workerCount {
				a, _, _ := fn(i)
				xys[j].X = float64(i)
				xys[j].Y = float64(a)
				j++
			}
			completed <- xys
		})(w, workers, limit, completed)
	}
	wg.Wait()
	close(completed)

	maxX := float64(0)
	maxY := float64(0)
	for v := range completed {
		for _, xy := range v {
			maxX = max(maxX, xy.X)
			maxY = max(maxY, xy.Y)
		}
		h, err := plotter.NewScatter(v)
		if err != nil {
			panic(err)
		}
		h.Color = fill
		p.Legend.TextStyle.Font.Size = 20
		p.Add(h)
	}
	p.X.Max = max(p.X.Max, maxX)
	p.Y.Max = max(p.Y.Max, maxY)
	p.Legend.Add(fmt.Sprintf("max %s = %d", title, int(maxY)))
}

func buildHistogram(p *plot.Plot, fill color.NRGBA, limit uint64, fn func(n uint64) (uint64, uint64, uint64), title string) {
	var wg sync.WaitGroup
	workers := uint64(runtime.GOMAXPROCS(0))
	completed := make(chan []uint64, workers)
	for w := uint64(1); w <= workers; w++ {
		wg.Add(1)
		go (func(worker uint64, workerCount uint64, limit uint64, completed chan<- []uint64) {
			defer wg.Done()
			values := make([]uint64, 100_000)
			for i := worker; i < limit; i += workerCount {
				a, _, _ := fn(i)
				values[a]++
			}
			completed <- values
		})(w, workers, limit, completed)
	}
	wg.Wait()
	close(completed)

	maxX := 0
	values := make([]uint64, 100_000)
	for v := range completed {
		for i := 0; i < len(v); i++ {
			if v[i] > 0 {
				maxX = i
				values[i] += v[i]
			}
		}
	}

	p.X.Max = max(p.X.Max, float64(maxX))
	xys := make(plotter.XYs, maxX)
	for i := 0; i < maxX; i++ {
		xys[i].X = float64(i)
		xys[i].Y = float64(values[i])
	}
	h, err := plotter.NewHistogram(xys, maxX)
	if err != nil {
		panic(err)
	}
	h.LineStyle.Width = 0
	h.Width = 1
	h.FillColor = fill
	p.Legend.Add(fmt.Sprintf("max %s = %d", title, int(maxX)))
	p.Add(h)
}
