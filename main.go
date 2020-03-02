package main

import (
	"math/rand"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

type ClusterManager struct {
	Zone string

	oomCountByHost map[string]int

	OOMCountDesc *prometheus.Desc
}

// Describe simply sends the two Descs in the struct to the channel.
func (c *ClusterManager) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.OOMCountDesc
}

func (c *ClusterManager) Collect(ch chan<- prometheus.Metric) {
	for host, _ := range c.oomCountByHost {
		c.oomCountByHost[host] += rand.Int() % 100
		ch <- prometheus.MustNewConstMetric(
			c.OOMCountDesc,
			prometheus.CounterValue,
			float64(c.oomCountByHost[host]),
			host,
		)
	}
}

func NewClusterManager() (cm *ClusterManager) {
	return &ClusterManager{
		oomCountByHost: map[string]int{
			"host-17": 0,
			"host-23": 0,
		},
		OOMCountDesc: prometheus.NewDesc(
			"clustermanager_oom_crashes_total",
			"Number of OOM crashes.",
			[]string{"host"},
			prometheus.Labels{},
		),
	}
}

func main() {
	var clusterMetricScraper = NewClusterManager()

	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(clusterMetricScraper)
	gatherers := prometheus.Gatherers{
		reg,
	}

	h := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
	log.Infoln("Start server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Errorf("Error occur when start server %v", err)
		os.Exit(1)
	}
}
