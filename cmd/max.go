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
	maxCmd = &cobra.Command{
		Use:   "max",
		Short: "Generate a scatter plot for the maximum value reached recursively for the given function",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, ok := Types[fn]
			if !ok {
				return fmt.Errorf("invalid value for --fn: %s", fn)
			}
			if power < 1 || power > 20 {
				return fmt.Errorf("--k must be in the range 1..20: %d", power)
			}
			limit := uint64(math.Pow10(power))

			p := newPlot()
			p.Title.Text = fmt.Sprintf("%s Maximum Reached 10^%d", info.Title, power)
			p.X.Label.Text = "n"
			p.Y.Label.Text = "Max Reached"
			p.Y.Min = 1
			p.X.Min = 1
			if fn == "" || fn == "f" {
				log.Printf("building f(x) scatter for 10^%d...", power)
				buildMax(p, color.NRGBA{R: 255, G: 0, B: 0, A: 128}, limit, shared.CollatzStoppingTimeF, "f(x)")
			}
			if fn == "" || fn == "g" {
				log.Printf("building g(x) scatter for 10^%d...", power)
				buildMax(p, color.NRGBA{R: 0, G: 255, B: 0, A: 128}, limit, shared.CollatzStoppingTimeG, "g(x)")
			}
			if fn == "" || fn == "h" {
				log.Printf("building h(x) scatter for 10^%d...", power)
				buildMax(p, color.NRGBA{R: 0, G: 0, B: 255, A: 128}, limit, shared.CollatzStoppingTimeH, "h(x)")
			}
			applyConstraintsToPlot(p, minX, minY, maxX, maxY)
			fileName := fmt.Sprintf("max_%s_%d.png", info.File, power)
			return saveToPNG(fileName, 750, 1500, p)
		},
	}
)

func init() {
	maxCmd.Flags().StringVar(&fn, "fn", "", "which function to plot: f, g, h. leave blank for all")
	maxCmd.Flags().IntVar(&power, "k", 7, "examine n up to 10^k")
	maxCmd.Flags().Float64Var(&minX, "min-x", 0, "min x to show on plot. use 0 for min of data")
	maxCmd.Flags().Float64Var(&minY, "min-y", 0, "min y to show on plot. use 0 for min of data")
	maxCmd.Flags().Float64Var(&maxX, "max-x", 10_000, "max x to show on plot. use 0 for max of data")
	maxCmd.Flags().Float64Var(&maxY, "max-y", 100_000, "max y to show on plot. use 0 for max of data")
}

func buildMax(p *plot.Plot, fill color.NRGBA, limit uint64, fn func(n uint64) (uint64, uint64, uint64), title string) {
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
				_, _, a := fn(i)
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
		p.Add(h)
	}
	p.X.Max = max(p.X.Max, maxX)
	p.Y.Max = max(p.Y.Max, maxY)
	p.Legend.Add(fmt.Sprintf("max %s = %d", title, int(maxY)))
}
