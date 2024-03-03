// Package presentation содержит фабрику роутера, обработчики и схемы валидации
package presentation

import (
	"net/http"
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
			"uri":      uri,
			"method":   method,
			"duration": duration,
		}).Info("Request info")

		log.WithFields(logrus.Fields{
			"status":         recorder.Status,
			"content length": recorder.ContentLength,
		}).Info("Response info")
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
		UserID, err := jose.ParseUserID([]byte(token))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		handlerFn(w, r, UserID)
	})
}
