package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/steffenmllr/prometheus-webpagetest-exporter/webpagetest"
	"github.com/tidwall/gjson"
)

type Metrics struct {
	Key  string
	Help string
	Data []string
}

var RunsLabels = []string{"url", "run", "location"}
var SiteMetrics map[string]*prometheus.GaugeVec

var TotalRunCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "total_run_counter",
		Help: "Counter of the total tests.",
	})

func InitMetrics(metrics []Metrics) {
	SiteMetrics = make(map[string]*prometheus.GaugeVec, len(metrics))

	for _, metric := range metrics {
		SiteMetrics[metric.Key] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: metric.Key,
				Help: metric.Help,
			}, RunsLabels)

		prometheus.MustRegister(SiteMetrics[metric.Key])
	}
}

func ReportResult(result *webpagetest.ResultData, metrics []Metrics, location string) {
	url := gjson.Get(result.Data, "url")

	// Go through all the metrics in the config
	for _, cmetric := range metrics {
		// Get the registered metrics
		if metric, ok := SiteMetrics[cmetric.Key]; ok {
			// for every data element, fetch the value
			for ix, dataKey := range cmetric.Data {
				jsonValue := gjson.Get(result.Data, dataKey)
				run := ix + 1
				// If the value exists, add it
				if jsonValue.Exists() {
					metric.WithLabelValues(url.String(), fmt.Sprint(run), location).Set(jsonValue.Float())
				}
			}

		}
	}
}
