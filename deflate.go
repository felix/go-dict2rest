package main

import (
	"compress/flate"
	"io"
	"net/http"
	"strings"
	"github.com/alexedwards/stack"
)

// DEFLATE Compression
type flateResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w flateResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Deflate(ctx *stack.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "deflate") ||
		   ctx.Get("handled").(bool) == true {
			next.ServeHTTP(w, req)
			return
		}
		fl, err := flate.NewWriter(w, -1) // Use default compression level
		if err != nil {
			next.ServeHTTP(w, req)
			return
		}
		w.Header().Set("Content-Encoding", "deflate")
		defer fl.Close()
		flw := flateResponseWriter{Writer: fl, ResponseWriter: w}
		ctx.Put("handled", true)
		defer ctx.Delete("handled")
		next.ServeHTTP(flw, req)
	})
}
