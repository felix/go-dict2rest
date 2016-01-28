package main

import (
	"log"
	"net/http"
	"time"
)

// Logging middleware
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Start buffered write
		bw := new(Buffer)

		next.ServeHTTP(bw, r)

		// Write out the buffer
		size, _ := bw.Apply(w)

		// Usually milliseconds
		latency := time.Since(startTime).Seconds()

		log.Printf("%s %s %d %d (%v)\n", r.Method, r.URL.Path, bw.Status, size, latency)
	})
}
