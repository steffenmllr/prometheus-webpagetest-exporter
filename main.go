package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/steffenmllr/prometheus-webpagetest-exporter/exporter"
	"github.com/steffenmllr/prometheus-webpagetest-exporter/webpagetest"
)

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type Config struct {
	Key     string
	Host    string
	Port    string
	Timer   duration
	Metrics []exporter.Metrics
	Sites   []struct {
		Name       string
		Url        string
		Location   string
		Connection string
	}
}

var queue *exporter.ListQueue

func main() {

	configFile := os.Args[1]
	if configFile == "" {
		fmt.Printf("Please set the config file")
		os.Exit(1)
	}

	// Read Config
	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	// Start a queue
	queue = exporter.NewListQueue()

	// Init the Metrics
	exporter.InitMetrics(config.Metrics)

	// Start client
	wpt, err := webpagetest.NewClient(config.Host)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	runChecks := func() {
		for _, site := range config.Sites {
			// Increment total runs
			exporter.TotalRunCounter.Inc()

			// Get Test Results
			testResult, err := wpt.RunTestAndWait(webpagetest.TestSettings{
				URL:      site.Url,
				Location: site.Location,
				Runs:     1,
			}, queue.UpdateStatus)

			if err != nil {
				fmt.Printf("Error: %v", err)
			} else {
				// Report Results
				exporter.ReportResult(testResult, config.Metrics, site.Location)
				queue.AddTestResult(testResult, config.Host)
			}
		}
	}

	go func() {
		// Run frist check on startup
		runChecks()

		// Check within duration
		for _ = range time.Tick(config.Timer.Duration) {
			runChecks()
		}
	}()

	// Expose the registered metrics via HTTP.
	http.Handle("/", queue)
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting on Port :%v", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))

}
