package web

import (
	"embed"
	"github.com/sirupsen/logrus"
	"io/fs"
)

//go:embed assets
var embedFS embed.FS

func GetWebAssetsFs() (fs.FS, error) {
	webFiles, err := fs.Sub(embedFS, "assets")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open sub FS of static files")
		return nil, err
	}
	return webFiles, nil
}
