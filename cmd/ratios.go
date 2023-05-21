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
	"gonum.org/v1/plot/vg"
)

var (
	ratiosCmd = &cobra.Command{
		Use:   "ratios",
		Short: "Generate a line graph or histogram between the ratio of h(x) over f(x) or g(x)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fn != "f" && fn != "g" {
				return fmt.Errorf("--fn should be f or g")
			}
			info, ok := Types[fn]
			if !ok {
				return fmt.Errorf("invalid value for --fn: %s", fn)
			}
			if power < 1 || power > 20 {
				return fmt.Errorf("--k must be in the range 1..20: %d", power)
			}
			limit := uint64(math.Pow10(power))
			group, err := cmd.Flags().GetUint64("group")
			if err != nil {
				return err
			}
			graphType, err := cmd.Flags().GetString("graph")
			if err != nil {
				return err
			}

			p := newPlot()
			p.Title.Text = fmt.Sprintf("Σh(x)/Σ%s 10^%d", info.Title, power)
			if graphType == "line" {
				p.X.Label.Text = "x"
				p.Y.Label.Text = fmt.Sprintf("Σh(x)/Σ%s", info.Title)
				log.Printf("building line graph for 10^%d...", power)
				if fn == "g" {
					buildRatioLine(p, color.NRGBA{R: 0, G: 255, B: 255, A: 128}, limit, group, shared.CollatzStoppingTimeH, shared.CollatzStoppingTimeG, "g(x)")
				} else if fn == "f" {
					buildRatioLine(p, color.NRGBA{R: 0, G: 0, B: 255, A: 128}, limit, group, shared.CollatzStoppingTimeH, shared.CollatzStoppingTimeF, "f(x)")
				} else {
					return fmt.Errorf("unexpected value for --fn: %s", fn)
				}
			} else if graphType == "histogram" {
				p.X.Label.Text = fmt.Sprintf("Σh(x)/Σ%s", info.Title)
				p.Y.Label.Text = "Count"
				log.Printf("building histogram for 10^%d...", power)
				if fn == "g" {
					buildRatioHistogram(p, color.NRGBA{R: 0, G: 255, B: 0, A: 128}, limit, group, shared.CollatzStoppingTimeH, shared.CollatzStoppingTimeG, "g(x)")
				} else if fn == "f" {
					buildRatioHistogram(p, color.NRGBA{R: 0, G: 0, B: 255, A: 128}, limit, group, shared.CollatzStoppingTimeH, shared.CollatzStoppingTimeF, "f(x)")
				} else {
					return fmt.Errorf("unexpected value for --fn: %s", fn)
				}
			} else {
				return fmt.Errorf("unexpected value for --graph: %s", graphType)
			}
			applyConstraintsToPlot(p, minX, minY, maxX, maxY)
			fileName := fmt.Sprintf("ratios_%s_%s_%d.png", graphType, info.File, power)
			return saveToPNG(fileName, 1500, 900, p)
		},
	}
)

func init() {
	ratiosCmd.Flags().StringVar(&fn, "fn", "", "which function to compare h to: f, g")
	ratiosCmd.Flags().IntVar(&power, "k", 5, "examine n up to 10^k")
	ratiosCmd.Flags().String("graph", "", "plot using line or histogram")
	ratiosCmd.Flags().Uint64("group", 5000, "number of x to group into each data point")
	ratiosCmd.Flags().Float64Var(&minX, "min-x", 0, "min x to show on plot. use 0 for min of data")
	ratiosCmd.Flags().Float64Var(&minY, "min-y", 0, "min y to show on plot. use 0 for min of data")
	ratiosCmd.Flags().Float64Var(&maxX, "max-x", 0, "max x to show on plot. use 0 for max of data")
	ratiosCmd.Flags().Float64Var(&maxY, "max-y", 0, "max y to show on plot. use 0 for max of data")
}

func buildRatioLine(p *plot.Plot, fill color.NRGBA, limit uint64, group uint64, fnN func(n uint64) (uint64, uint64, uint64), fnD func(n uint64) (uint64, uint64, uint64), title string) {
	var wg sync.WaitGroup
	workers := uint64(runtime.GOMAXPROCS(0))
	numerator := make([]uint64, limit/group)
	denominator := make([]uint64, limit/group)
	for w := uint64(0); w < workers; w++ {
		wg.Add(1)
		go (func(worker uint64, workerCount uint64, limit uint64) {
			defer wg.Done()
			for i := worker * group; i < limit; i += workerCount * group {
				for j := uint64(0); j < group; j++ {
					a, _, _ := fnN(1 + i + j)
					b, _, _ := fnD(1 + i + j)
					numerator[i/group] += a
					denominator[i/group] += b
				}
			}
		})(w, workers, limit)
	}
	wg.Wait()
	numeratorSum := uint64(0)
	denominatorSum := uint64(0)
	xys := make(plotter.XYs, limit/group)
	for i := 0; i < int(limit/group); i++ {
		numeratorSum += numerator[i]
		denominatorSum += denominator[i]
		xys[i].X = float64(uint64(i) * group)
		xys[i].Y = float64(numeratorSum) / float64(denominatorSum)
	}
	h, err := plotter.NewLine(xys)
	if err != nil {
		panic(err)
	}
	h.LineStyle.Width = vg.Points(1.5)
	h.Color = fill
	p.Legend.TextStyle.Font.Size = 20
	p.Add(h)
	applyConstraintsToPlot(p, 1, float64(limit), xys[len(xys)-1].Y, xys[0].Y)
}

func buildRatioHistogram(p *plot.Plot, fill color.NRGBA, limit uint64, group uint64, fnN func(n uint64) (uint64, uint64, uint64), fnD func(n uint64) (uint64, uint64, uint64), title string) {
	var wg sync.WaitGroup
	workers := uint64(runtime.GOMAXPROCS(0))
	numerator := make([]uint64, limit)
	denominator := make([]uint64, limit)
	for w := uint64(0); w < workers; w++ {
		wg.Add(1)
		go (func(worker uint64, workerCount uint64, limit uint64) {
			defer wg.Done()
			for i := worker; i < limit; i += workerCount {
				a, _, _ := fnN(1 + i)
				b, _, _ := fnD(1 + i)
				numerator[i] += a
				denominator[i] += b
			}
		})(w, workers, limit)
	}
	wg.Wait()
	numeratorSum := uint64(0)
	denominatorSum := uint64(0)
	minX := float64(1)
	maxX := float64(0)
	xys := make(plotter.Values, group+1)
	for i := uint64(1); i < limit; i++ {
		numeratorSum += numerator[i]
		denominatorSum += denominator[i]
		k := int(float64(numeratorSum) / float64(denominatorSum) * float64(group))
		xys[k]++
	}
	filteredXys := make(plotter.XYs, 0)
	maxY := float64(0)
	for i := 0; i < len(xys); i++ {
		if xys[i] > math.Log10(float64(group)) {
			x := float64(i) / float64(group)
			filteredXys = append(filteredXys, plotter.XY{X: x, Y: xys[i]})
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if xys[i] > maxY {
				maxY = xys[i]
			}
		}
	}
	h, err := plotter.NewHistogram(filteredXys, len(filteredXys))
	if err != nil {
		panic(err)
	}
	h.LineStyle.Width = 0
	h.FillColor = fill
	p.Legend.TextStyle.Font.Size = 20
	p.Add(h)
	applyConstraintsToPlot(p, float64(minX), float64(maxX), 0, maxY)
}
