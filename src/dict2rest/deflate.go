package main

import (
	"compress/flate"
	"io"
	"net/http"
	"strings"
)

// DEFLATE Compression
type flateResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w flateResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Deflate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "deflate") {
			// If deflate is unsupported, revert to standard handler.
			next.ServeHTTP(w, req)
			return
		}
		w.Header().Set("Content-Encoding", "deflate")
		fl, err := flate.NewWriter(w, -1) // Use default compression level
		if err != nil {
			next.ServeHTTP(w, req)
			return
		}
		defer fl.Close()
		flw := flateResponseWriter{Writer: fl, ResponseWriter: w}
		next.ServeHTTP(flw, req)
	})
}
