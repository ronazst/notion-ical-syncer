package handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func wrap(handlerFunc func(request *http.Request) (string, error)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		result, err := handlerFunc(request)
		if err != nil {
			logrus.WithError(err).Error("Error while handling request")
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(result))
	}

}
