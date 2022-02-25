/*
	Random number generator.

	This is a sample application used to demonstrate instrumentation
	using exporters.
*/
package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

/*
	randomNumber is a collector as it implements prometheus.Collector interface.
*/
type randomNumber struct {
	randomNumber *prometheus.Desc
}

var _ prometheus.Collector = randomNumber{}

var (
	/*
		prometheus.NewDesc returns a new prometheus.Desc.
		prometheus.Desc is a descriptor used by Prometheus to store immutable metric metadata.
	*/
	randomNumberDesc = prometheus.NewDesc(
		// Name of the metric.
		prometheus.BuildFQName("random", "number", "generated"),
		// Help text for the metric.
		"A randomly generated number",
		// Variable label names.
		nil,
		// Constant label values.
		prometheus.Labels{
			"app": "random_number_generator",
		},
	)
)

func main() {
	flag.Parse()

	// Create a new randomNumber collector.
	collector := randomNumber{
		randomNumber: randomNumberDesc,
	}
	// Register the custom collector.
	prometheus.MustRegister(collector)

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{},
	))
	log.Printf("Running server at %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// Describe is a required method for prometheus.Collector.
func (r randomNumber) Describe(ch chan<- *prometheus.Desc) {
	ch <- r.randomNumber
}

var collectCounter int

// Collect is a required method for prometheus.Collector.
func (r randomNumber) Collect(ch chan<- prometheus.Metric) {
	log.Printf("Collect Counter: %d\n", collectCounter)
	/*
		collectCounter is used to count the number of times the Collect method is called.
		This is used to demonstrate how stale metrics are not exposed once we stop calling MustNewConstMetric.
	*/
	if collectCounter < 5 {
		/*
			MustNewConstMetric returns a prometheus.Metric with one fixed value that cannot be changed.
			It is a throw-away metric that is generated on the fly.
		*/
		ch <- prometheus.MustNewConstMetric(randomNumberDesc, prometheus.GaugeValue, rand.Float64())
	}
	collectCounter++
}

/*
Output:

Terminal 1
go run main.go

2022/02/25 14:28:03 Running server at :8080
2022/02/25 14:28:07 Collect Counter: 0
2022/02/25 14:28:09 Collect Counter: 1
2022/02/25 14:28:11 Collect Counter: 2
2022/02/25 14:28:13 Collect Counter: 3
2022/02/25 14:28:15 Collect Counter: 4
2022/02/25 14:28:17 Collect Counter: 5
2022/02/25 14:28:19 Collect Counter: 6
2022/02/25 14:28:21 Collect Counter: 7
2022/02/25 14:28:23 Collect Counter: 8
2022/02/25 14:28:25 Collect Counter: 9


Terminal 2
for i in {1..10};do curl -s localhost:8080/metrics|grep random_number_generated; sleep 2; echo "";done

# HELP random_number_generated A randomly generated number
# TYPE random_number_generated gauge
random_number_generated{app="random_number_generator"} 0.6046602879796196

# HELP random_number_generated A randomly generated number
# TYPE random_number_generated gauge
random_number_generated{app="random_number_generator"} 0.9405090880450124

# HELP random_number_generated A randomly generated number
# TYPE random_number_generated gauge
random_number_generated{app="random_number_generator"} 0.6645600532184904

# HELP random_number_generated A randomly generated number
# TYPE random_number_generated gauge
random_number_generated{app="random_number_generator"} 0.4377141871869802

# HELP random_number_generated A randomly generated number
# TYPE random_number_generated gauge
random_number_generated{app="random_number_generator"} 0.4246374970712657

*/

/*
References:
	https://prometheus.io/
	https://www.oreilly.com/library/view/prometheus-up/9781492034131/
	https://github.com/prometheus/client_golang/
*/
