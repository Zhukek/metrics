package compress

import (
	"net/http"
	"strings"

	"github.com/Zhukek/metrics/internal/service/gzip"
)

const (
	gzipformat = "gzip"
)

func GzipMiddleware(handler http.Handler) http.HandlerFunc {
	gzipFunc := func(w http.ResponseWriter, r *http.Request) {
		resWriter := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportFormat := strings.Contains(acceptEncoding, gzipformat)

		if supportFormat {
			w.Header().Set("Content-Encoding", "gzip")
			compressWriter := gzip.NewGzipWriter(w)
			resWriter = compressWriter
			defer compressWriter.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		contentFormat := strings.Contains(contentEncoding, gzipformat)
		if contentFormat {
			compressReader, err := gzip.NewGzipReader(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			defer compressReader.Close()
			r.Body = compressReader
		}

		handler.ServeHTTP(resWriter, r)
	}

	return http.HandlerFunc(gzipFunc)
}
