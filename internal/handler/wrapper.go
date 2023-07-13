package handler

import (
	"github.com/ronazst/notion-ical-syncer/internal/util"
	"net/http"
)

func wrap(handlerFunc func(request *http.Request) (string, error)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		result, err := handlerFunc(request)
		if err != nil {
			errorHandler(writer, err)
			return
		}
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(result))
	}
}

func errorHandler(w http.ResponseWriter, err error) {
	switch value := err.(type) {
	case *util.Error:
		w.WriteHeader(value.Code())
		_, _ = w.Write([]byte(value.Error()))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(util.MessageInternalServerError))
	}
}
