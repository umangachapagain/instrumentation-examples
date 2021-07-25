/*
	Directly instrumented random number generator.

	This is a sample application used to demonstrate direct instrumentation
	of code using Prometheus client_golang.
*/
package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

var (
	/*
		NewGauge() returns a Gauge based on provided GaugeOpts.

		Gauge is a Metric that represents a single numerical value that can
		arbitrarily go up and down.

		A Gauge is typically used for measured values like temperatures or current
		memory usage, but also "counts" that can go up and down, like the number of
		running goroutines.

		Gauge implements Collector interface, so it can be used by Prometheus to collect metrics.
	*/
	randomNumber = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "random",
		Subsystem: "number",
		Name:      "generated",
		Help:      "A randomly generated number",
		ConstLabels: prometheus.Labels{
			"app": "random_number_generator",
		},
	})
)

func init() {
	/*
		Collectors must be registered for collection.
		MustRegister() registers collectors to registry.

		You can use the default registry or create a custom registry for registration.
		Default registry comes with pre-registered collectors.

		Registry implements Registerer and Gatherer interface.

		Here, randomNumber (which is a collector) is registered to default registry.
	*/
	prometheus.MustRegister(randomNumber)
}

func main() {
	flag.Parse()

	/*
		goroutine that writes a random number to metrics value
		and sleeps for 10s.
	*/
	go func() {
		for {
			rnd := rand.Float64()
			// Set sets the randomNumber (of type Gauge) to an arbitrary float64 value.
			randomNumber.Set(rnd)
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{},
	))
	log.Printf("Running server at %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

/*
Output:

	Terminal 1:
	go run main.go

	Terminal 2:
	for i in {1..50};do curl -s localhost:8080/metrics|grep random_number_generated; sleep 3; echo "";done

	# HELP random_number_generated A randomly generated number
	# TYPE random_number_generated gauge
	random_number_generated{app="random_number_generator"} 0.6046602879796196

	# HELP random_number_generated A randomly generated number
	# TYPE random_number_generated gauge
	random_number_generated{app="random_number_generator"} 0.6046602879796196

	# HELP random_number_generated A randomly generated number
	# TYPE random_number_generated gauge
	random_number_generated{app="random_number_generator"} 0.6046602879796196

	# HELP random_number_generated A randomly generated number
	# TYPE random_number_generated gauge
	random_number_generated{app="random_number_generator"} 0.9405090880450124

	# HELP random_number_generated A randomly generated number
	# TYPE random_number_generated gauge
	random_number_generated{app="random_number_generator"} 0.9405090880450124

	# HELP random_number_generated A randomly generated number
	# TYPE random_number_generated gauge
	random_number_generated{app="random_number_generator"} 0.9405090880450124
*/

/*
References:
	https://prometheus.io/
	https://www.oreilly.com/library/view/prometheus-up/9781492034131/
	https://github.com/prometheus/client_golang/
*/
