package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const defaultTariffLocation = "Europe/Madrid"

func main() {
	addr := flag.String("addr", ":8080", "Address to listen at")
	location := flag.String("tariff-location", defaultTariffLocation, "Location for the tariff where the timezone is")
	flag.Parse()

	log.Printf("Starting HTTP server in %s", *addr)
	if err := http.ListenAndServe(*addr, http.HandlerFunc(periodHandler(*location))); err != nil {
		log.Fatalf("Cannot start http server: %v", err)
	}
}

func periodHandler(locationStr string) func(rw http.ResponseWriter, r *http.Request) {
	location, err := time.LoadLocation(locationStr)
	if err != nil {
		log.Fatalf("Error loading timezone: %v", err)
	}

	log.Printf("Using timezone %s", location)

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("content-type", "text/plain")

		p := periodAt(time.Now().In(location))

		fmt.Fprintln(rw, "# HELP The current period.")
		fmt.Fprintln(rw, "# TYPE period gauge")
		fmt.Fprintf(rw, "period{name=\"valley\"} %d\n", p.GaugeEquals(periodValley))
		fmt.Fprintf(rw, "period{name=\"plain\"} %d\n", p.GaugeEquals(periodPlain))
		fmt.Fprintf(rw, "period{name=\"peak\"} %d\n", p.GaugeEquals(periodPeak))

		fmt.Fprintln(rw, "# HELP The current period as a number.")
		fmt.Fprintln(rw, "# TYPE period_value gauge")
		fmt.Fprintf(rw, "period_value %d\n", p)
	}
}

type period int

func (p period) GaugeEquals(other period) int {
	if p == other {
		return 1
	}

	return 0
}

const (
	_ = iota
	periodValley
	periodPlain
	periodPeak
)

func periodAt(t time.Time) period {
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return periodValley
	}

	hour := t.Hour()
	switch true {
	case hour < 8:
		return periodValley
	case hour < 14:
		return periodPeak
	case hour < 18:
		return periodPlain
	case hour < 22:
		return periodPeak
	case hour < 24:
		return periodPlain
	default:
		panic("unreachable hour")
	}
}
