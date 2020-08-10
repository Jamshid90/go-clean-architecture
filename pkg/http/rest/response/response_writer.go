package response

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	code int
}

func NewResponseWriter(w http.ResponseWriter, statusCode int) *ResponseWriter {
	return &ResponseWriter{w, statusCode}
}

func (w *ResponseWriter) StatusCode() int {
	return w.code
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
