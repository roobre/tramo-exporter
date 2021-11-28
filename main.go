package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("addr", ":8080", "Address to listen at")
	flag.Parse()

	log.Printf("Starting HTTP server in %s", *addr)
	if err := http.ListenAndServe(*addr, http.HandlerFunc(periodHandler)); err != nil {
		log.Fatalf("Cannot start http server: %v", err)
	}
}

func periodHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "text/plain")

	fmt.Fprintln(rw, "# HELP The current period.")
	fmt.Fprintln(rw, "# TYPE period gauge")

	t := periodAt(time.Now())
	fmt.Fprintf(rw, "period{name=\"valley\"} %d\n", asInt(t == periodValley))
	fmt.Fprintf(rw, "period{name=\"plain\"} %d\n", asInt(t == periodPlain))
	fmt.Fprintf(rw, "period{name=\"peak\"} %d\n", asInt(t == periodPeak))
}

type period int

const (
	periodValley = iota
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

func asInt(b bool) int {
	if b {
		return 1
	}

	return 0
}
