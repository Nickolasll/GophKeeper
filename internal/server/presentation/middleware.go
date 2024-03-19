package presentation

import (
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type responseRecorder struct {
	http.ResponseWriter
	Status        int
	ContentLength int
}

func logging(handler http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		recorder := &responseRecorder{
			ResponseWriter: w,
			Status:         200,
			ContentLength:  0,
		}
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		handler.ServeHTTP(recorder, r)

		duration := time.Since(start)

		log.WithFields(logrus.Fields{
			"uri":            uri,
			"method":         method,
			"duration_ms":    duration.Milliseconds(),
			"status":         recorder.Status,
			"content_length": recorder.ContentLength,
		}).Info("Request info")
	}

	return http.HandlerFunc(logFn)
}

func auth(handlerFn authenticatedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}
		UserID, err := joseService.ParseUserID([]byte(token))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		handlerFn(w, r, UserID)
	})
}

func compress(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, reader *http.Request) {
		originalWriter := writer

		acceptEncoding := reader.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			compressWriter := newCompressWriter(writer)
			originalWriter = compressWriter
			defer func() {
				err := compressWriter.Close()
				if err != nil {
					log.Error(err)
				}
			}()
		}

		contentEncoding := reader.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(reader.Body)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)

				return
			}
			reader.Body = cr
			defer func() {
				err := cr.Close()
				if err != nil {
					log.Error(err)
				}
			}()
		}

		handler.ServeHTTP(originalWriter, reader)
	})
}
