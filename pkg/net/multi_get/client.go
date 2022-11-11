package multi_get

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/kirillgrachoff/load_tester/pkg/net/xhttp"
)

type MultiGet struct {
	count     int
	url       []string
	keepAlive bool
	r         *rand.Rand
}

type Logger interface {
	Println(v ...any)
	Printf(format string, args ...any)
}

func NewClient(count int, url []string, keepAlive bool) *MultiGet {
	return &MultiGet{
		count:     count,
		url:       url,
		keepAlive: keepAlive,
		r:         rand.New(rand.NewSource(42)),
	}
}

func (g *MultiGet) Run(ctx context.Context) (err error) {
	group, ctx := errgroup.WithContext(ctx)

	for i := 0; i < g.count; i++ {
		logger := log.New(os.Stdout, fmt.Sprintf("(index: %4d) ", i+1), log.Flags())
		group.Go(func() error {
			return g.worker(ctx, logger)
		})
	}

	err = group.Wait()
	if ctx.Err() != nil {
		log.Println(ctx.Err())
	}
	return
}

func (g *MultiGet) worker(ctx context.Context, logger Logger) error {
	logger.Printf("start querying")
	count := 0
	for {
		index := g.r.Intn(len(g.url))
		select {
		case <-ctx.Done():
			return nil
		case resp := <-xhttp.Get(g.url[index]):
			if resp.Err != nil {
				logger.Printf("total count: %6d | error while query: %s | time: %s", count, resp.Err, resp.Time)
			} else {
				if !g.keepAlive {
					resp.Response.Body.Close()
				}
				logger.Printf("total count: %6d | status: %s | time: %s", count, resp.Response.Status, resp.Time)
			}
			count++
		}
	}
}
