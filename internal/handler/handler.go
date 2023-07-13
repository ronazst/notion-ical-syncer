package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ronazst/notion-ical-syncer/internal/util"
	"github.com/ronazst/notion-ical-syncer/internal/web"
	"net/http"
)

func NewHandler() (http.Handler, error) {
	if err := util.CheckRequiredEnv(); err != nil {
		return nil, err
	}
	webFs, err := web.GetWebAssetsFs()
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	subRouter := router.PathPrefix(fmt.Sprintf("/%s", util.GetOsEnv(util.EnvStackId))).Subrouter()
	subRouter.Methods(http.MethodGet).Path("/webui").Handler(http.StripPrefix("/webui", http.FileServer(http.FS(webFs))))
	subRouter.Methods(http.MethodGet).Path("/ical").HandlerFunc(wrap(iCalHandler))

	return router, nil
}
