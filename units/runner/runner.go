package runner

import (
	"context"
	"log"
	"sync"
	"time"
)

type Config struct {
	Concurrency int           // Amount of async requests runs
	Requests    int           // If 0 — work in duration mode
	Duration    time.Duration // If > 0 — work in time mode
	Timeout     time.Duration // Timout for each request
	RateLimit   int           // Ratel limit for duration testing
}

type Result struct {
	Total     int
	Success   int
	Failures  int
	Latencies []time.Duration
	Errors    map[string]int
}

func RunScenario(scenario func(ctx context.Context) error, cfg Config) Result {
	var mu sync.Mutex
	result := Result{Errors: map[string]int{}}
	sem := make(chan struct{}, cfg.Concurrency)
	var wg sync.WaitGroup

	// Case 1: If Requests specified — use fixed-request logic
	if cfg.Requests > 0 {
		for i := 0; i < cfg.Requests; i++ {
			wg.Add(1)
			sem <- struct{}{}
			go func() {
				defer func() {
					<-sem
					wg.Done()
				}()

				runOneRequest(scenario, cfg, &result, &mu)
			}()
		}
		wg.Wait()
		return result
	}

	// Case 2: If Duration specified — run until timeout
	if cfg.Duration > 0 {
		ticker := &time.Ticker{}
		if cfg.RateLimit > 0 {
			ticker = time.NewTicker(time.Second / time.Duration(cfg.RateLimit))
		}

		done := time.After(cfg.Duration)
	loop:
		for {
			select {
			case <-done:
				break loop
			default:
				if cfg.RateLimit > 0 {
					<-ticker.C // wait for next tick
				}

				sem <- struct{}{}
				wg.Add(1)
				go func() {
					defer func() {
						<-sem
						wg.Done()
					}()
					runOneRequest(scenario, cfg, &result, &mu)
				}()
			}
		}

		ticker.Stop()
		wg.Wait()
		return result
	}

	log.Fatal("either Requests or Duration must be set")
	return result
}

func runOneRequest(scenario func(ctx context.Context) error, cfg Config, result *Result, mu *sync.Mutex) {
	ctx := context.Background()
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	start := time.Now()
	err := scenario(ctx)
	latency := time.Since(start)

	mu.Lock()
	defer mu.Unlock()
	result.Total++
	result.Latencies = append(result.Latencies, latency)
	if err != nil {
		result.Failures++
		result.Errors[err.Error()]++
	} else {
		result.Success++
	}
}
