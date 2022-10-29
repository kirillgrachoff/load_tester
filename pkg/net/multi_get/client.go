package multi_get

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/kirillgrachoff/load_tester/pkg/net/xhttp"
)

type MultiGet struct {
	wg    sync.WaitGroup
	count int
	url   []string
	r     *rand.Rand
}

type Logger interface {
	Println(v ...any)
	Printf(format string, args ...any)
}

func NewClient(count int, url []string) *MultiGet {
	return &MultiGet{
		count: count,
		url:   url,
		r:     rand.New(rand.NewSource(42)),
	}
}

func (g *MultiGet) Run(ctx context.Context) error {
	g.wg.Add(g.count)

	for i := 0; i < g.count; i++ {
		logger := log.New(os.Stdout, fmt.Sprintf("(index: %4d) ", i+1), log.Flags())
		go g.worker(ctx, logger)
	}

	g.wg.Wait()
	return nil
}

func (g *MultiGet) worker(ctx context.Context, logger Logger) {
	defer g.wg.Done()
	logger.Printf("start querying")
	count := 0
	for {
		index := g.r.Intn(len(g.url))
		select {
		case <-ctx.Done():
			return
		case resp := <-xhttp.Get(g.url[index]):
			if resp.Err != nil {
				logger.Printf("error while query: %s | time: %s", resp.Err, resp.Time)
				return
			}
			count++
			logger.Printf("total count: %6d | status: %s | time: %s", count, resp.Response.Status, resp.Time)
		}
	}
}
