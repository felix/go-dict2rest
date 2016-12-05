package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"github.com/alexedwards/stack"
)

// Gzip Compression
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(ctx *stack.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") ||
		   ctx.Get("handled").(bool) == true {
			// move on to next handler in the chain
			next.ServeHTTP(w, req)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		ctx.Put("handled", true)
		defer ctx.Delete("handled")
		next.ServeHTTP(gzw, req)
	})
}
