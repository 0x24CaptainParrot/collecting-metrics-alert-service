package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type CompressWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

func (cw *CompressWriter) Write(p []byte) (int, error) {
	if cw.Header().Get("Content-Encoding") != "gzip" {
		cw.Header().Set("Content-Encoding", "gzip")
	}
	return cw.zw.Write(p)
}

func (cw *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		cw.Header().Set("Content-Encoding", "gzip")
	}
	cw.ResponseWriter.WriteHeader(statusCode)
}

func (cw *CompressWriter) Close() error {
	return cw.zw.Close()
}

type CompressReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &CompressReader{
		ReadCloser: r,
		zr:         zr,
	}, nil
}

func (cr *CompressReader) Read(p []byte) (int, error) {
	return cr.zr.Read(p)
}

func (cr *CompressReader) Close() error {
	if err := cr.ReadCloser.Close(); err != nil {
		return err
	}
	return cr.zr.Close()
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			cw := NewCompressWriter(w)
			defer cw.Close()
			w = cw
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			cr, err := NewCompressReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decode gzip", http.StatusInternalServerError)
				return
			}
			defer cr.Close()
			r.Body = cr
		}

		h.ServeHTTP(w, r)
	})
}
