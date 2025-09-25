package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseData struct {
	size   int
	status int
}

type loggingResponseWriter struct {
	responseData *responseData
	http.ResponseWriter
	headerWritten bool
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size += size
	return size, err
}

func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	if !l.headerWritten {
		l.ResponseWriter.WriteHeader(statusCode)
		l.headerWritten = true
	}
	l.responseData.status = statusCode
}

type Slogger struct {
	slogger *zap.SugaredLogger
}

func (s *Slogger) WithLogging(handler http.Handler) http.Handler {
	logFunc := func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		uri := req.URL
		method := req.Method

		responseData := responseData{0, 0}
		loggingResponseWriter := loggingResponseWriter{
			responseData:   &responseData,
			ResponseWriter: res,
		}

		handler.ServeHTTP(&loggingResponseWriter, req)

		duration := time.Since(start)

		s.slogger.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", loggingResponseWriter.responseData.status,
			"size", loggingResponseWriter.responseData.size,
		)
	}

	return http.HandlerFunc(logFunc)
}

func (s *Slogger) ErrLog(err error) {
	s.slogger.Errorln(err)
}

func (s *Slogger) Info(str string) {
	s.slogger.Infoln(str)
}

func (s *Slogger) Sync() {
	s.slogger.Sync()
}

func NewSlogger() (slogger *Slogger, err error) {
	logger, err := zap.NewDevelopment()

	if err != nil {
		return
	}
	slogger = &Slogger{logger.Sugar()}
	return
}
