package cmd

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"

	"github.com/spf13/cobra"
	"github.com/theriault/collatz/shared"
)

var (
	compareCmd = &cobra.Command{
		Use:   "compare",
		Short: "Compare g (reduced) or h (main result) to f to empirically verify equality",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fn == "" {
				return fmt.Errorf("--fn is required")
			}
			if power < 1 || power > 20 {
				return fmt.Errorf("--k must be in the range 1..20: %d", power)
			}
			limit := uint64(math.Pow10(power))
			if fn == "g" {
				log.Printf("comparing g(x) to f(x) from 1..10^%d", power)
				compare(limit, shared.CollatzStoppingTimeG, shared.CollatzStoppingTimeF)
			}
			if fn == "h" {
				log.Printf("comparing h(x) to f(x) from 1..10^%d", power)
				compare(limit, shared.CollatzStoppingTimeH, shared.CollatzStoppingTimeF)
			}
			return nil
		},
	}
)

func init() {
	compareCmd.Flags().StringVar(&fn, "fn", "", "which function to plot: g, h")
	compareCmd.Flags().IntVar(&power, "k", 5, "examine n up to 10^k")
}

func compare(limit uint64, fnA func(n uint64) (uint64, uint64, uint64), fnB func(n uint64) (uint64, uint64, uint64)) {
	var wg sync.WaitGroup
	workers := uint64(runtime.GOMAXPROCS(0))
	A := make([]uint64, limit)
	B := make([]uint64, limit)
	for w := uint64(1); w <= workers; w++ {
		wg.Add(1)
		go (func(worker uint64, workerCount uint64, limit uint64) {
			defer wg.Done()
			for i := worker; i < limit; i += workerCount {
				_, a, _ := fnA(i)
				b, _, _ := fnB(i)
				A[i] = a
				B[i] = b
			}
		})(w, workers, limit)
	}
	wg.Wait()
	for i := 1; i < int(limit); i++ {
		if A[i] != B[i] {
			fmt.Printf("%d: %d != %d\n", i, A[i], B[i])
			break
		}
	}
}
