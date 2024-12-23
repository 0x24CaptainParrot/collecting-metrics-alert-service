package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

type CompressWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	zw := gzipWriterPool.Get().(*gzip.Writer)
	zw.Reset(w)
	return &CompressWriter{
		ResponseWriter: w,
		zw:             zw,
	}
}

func (cw *CompressWriter) Write(p []byte) (int, error) {
	if cw.Header().Get("Content-Encoding") != "gzip" {
		cw.Header().Set("Content-Encoding", "gzip")
	}
	return cw.zw.Write(p)
}

func (cw *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 && IsCompressibleContentType(cw.Header().Get("Content-Type")) {
		cw.Header().Set("Content-Encoding", "gzip")
	}
	cw.ResponseWriter.WriteHeader(statusCode)
}

func (cw *CompressWriter) Close() {
	_ = cw.zw.Close()
	gzipWriterPool.Put(cw.zw)
}

type CompressReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

var gzipReaderPool = sync.Pool{
	New: func() interface{} {
		return &gzip.Reader{}
	},
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr := gzipReaderPool.Get().(*gzip.Reader)
	if err := zr.Reset(r); err != nil {
		gzipReaderPool.Put(zr)
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
	err := cr.ReadCloser.Close()
	cr.zr.Close()
	gzipReaderPool.Put(cr.zr)
	return err
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		acceptsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

		if acceptsGzip && IsCompressibleContentType(contentType) {
			r.Header.Set("Content-Encoding", "gzip")
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

var compressibleTypes = []string{
	"text/html",
	"application/json",
}

func IsCompressibleContentType(contentType string) bool {
	for _, ct := range compressibleTypes {
		if strings.HasPrefix(contentType, ct) {
			return true
		}
	}
	return false
}
