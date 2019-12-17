package main

import (
  "os"
  "time"
  "strconv"
  "net/http"
  vegeta "github.com/AkshatM/vegeta/lib"
  "github.com/dchest/uniuri"
)

// This maps a number to a list of headers. The header *values* are generated
// dynamically. The Nth header key is populated by chekcing if the Nth bit is
// set in `profile`. Since `profile` is a 64-bit integer, that means upto 64
// unique headers can be added.
func generateHeaders(profile uint64) http.Header {
	var headers http.Header
	if (profile & (1 << 1)) > 0 {
		headers.Add("user-agent", uniuri.New())
	}
	if (profile & (1 << 2)) > 0 {
		headers.Add("server", uniuri.New())
	}
	if (profile & (1 << 3)) > 0 {
		headers.Add("x-client-trace-id", uniuri.New())
	}
	if (profile & (1 << 4)) > 0 {
		headers.Add("x-envoy-downstream-service-cluster", uniuri.New())
	}
	if (profile & (1 << 5)) > 0 {
		headers.Add("x-envoy-downstream-service-node", uniuri.New())
	}
	if (profile & (1 << 6)) > 0 {
		headers.Add("x-envoy-external-address", uniuri.New())
	}
	if (profile & (1 << 7)) > 0 {
		headers.Add("x-envoy-force-trace", uniuri.New())
	}
	if (profile & (1 << 8)) > 0 {
		headers.Add("x-envoy-internal", uniuri.New())
	}
	return headers
}

// Use Vegeta to hit Envoy with desired rate, duration and headers.
// This generates latency metrics in HDR format (http://hdrhistogram.github.io/HdrHistogram/plotFiles.html)
// to stdout.
func main() {

  url := os.Args[1]
  profile, err := strconv.ParseUint(os.Args[2], 10, 64)
  if err != nil {
	  panic(err)
  }

  frequency, err := strconv.Atoi(os.Args[3])
  if err != nil {
	  panic(err)
  }

  length, err := strconv.Atoi(os.Args[4])
  if err != nil {
	  panic(err)
  }

  rate := vegeta.Rate{Freq: frequency, Per: time.Second}
  duration := time.Duration(length) * time.Second
  targeter := vegeta.NewStaticTargeter(vegeta.Target{
    Method: "GET",
    URL:    url,
    BodyContainsTimestamp: true,
    Header: generateHeaders(profile),
  })
  attacker := vegeta.NewAttacker()

  var metrics vegeta.Metrics
  for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
    metrics.Add(res)
  }
  metrics.Close()

  reporter := vegeta.NewHDRHistogramPlotReporter(&metrics)
  reporter.Report(os.Stdout)
}
