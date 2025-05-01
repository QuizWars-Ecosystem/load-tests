package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/QuizWars-Ecosystem/load-tests/units/report"
	"github.com/QuizWars-Ecosystem/load-tests/units/runner"
	"github.com/QuizWars-Ecosystem/load-tests/units/scenario"
	_ "github.com/QuizWars-Ecosystem/load-tests/units/scenario/samples"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "Address of the gRPC server")
	scenarioList := flag.String("scenarios", "", "Comma-separated list of scenario keys")
	concurrency := flag.Int("concurrency", 10, "Concurrent workers per scenario")
	requests := flag.Int("requests", 100, "Total requests per scenario")
	durationStr := flag.String("duration", "", "How long to run scenario (e.g., 10s)")
	rpcStr := flag.String("rpc", "", "Request per second for duration tests")
	timeoutStr := flag.String("timeout", "3s", "Per-request timeout")
	reportFormat := flag.String("report", "print", "Report format, print|json|html")
	flag.Parse()

	if addr == nil || *addr == "" {
		log.Fatalf("-addr flag is required")
	}

	scenario.SetGlobalAddr(*addr)

	var duration time.Duration
	var rpc int
	var timeout time.Duration

	if *durationStr != "" {
		duration, _ = time.ParseDuration(*durationStr)
	}

	if *rpcStr != "" {
		rpc, _ = strconv.Atoi(*rpcStr)
	}

	timeout, _ = time.ParseDuration(*timeoutStr)

	cfg := runner.Config{
		Concurrency: *concurrency,
		Requests:    *requests,
		Duration:    duration,
		Timeout:     timeout,
		RateLimit:   rpc,
	}

	if *scenarioList == "" {
		log.Fatal("No scenarios specified")
	}

	names := strings.Split(*scenarioList, ",")
	var wg sync.WaitGroup

	scenariosToRun, err := scenario.ResolveScenarioOrder(names)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range scenariosToRun {
		log.Printf("Running scenario %s", s.Name)

		if s.Init != nil {
			if err = s.Init(); err != nil {
				log.Fatalf("Error initializing scenario %s: %v", s.Name, err)
			}
		}

		wg.Add(1)
		go func(s *scenario.Scenario) {
			defer wg.Done()
			result := runner.RunScenario(s.Call, cfg)
			switch *reportFormat {
			case "print":
				report.Print(s.Name, result)
			case "json":
				report.SaveJSON(s.Name, result)
			case "html":
				report.SaveHTML(s.Name, result)
			default:
				report.Print(s.Name, result)
			}
		}(s)
	}

	wg.Wait()
	fmt.Println("✅ All scenarios completed.")
}
