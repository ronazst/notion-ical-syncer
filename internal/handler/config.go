package handler

import (
	"encoding/json"
	"github.com/ronazst/notion-ical-syncer/internal/awsutil"
	"github.com/ronazst/notion-ical-syncer/internal/util"
	"github.com/sirupsen/logrus"
	"net/http"
)

func getConfigHandler(request *http.Request) (string, error) {
	configId := request.URL.Query().Get("config_id")
	if util.IsBlank(configId) {
		logrus.Error("get config with blank config_id")
		return "", util.NewUserInputError("config_id is required and can no be blank")
	}

	configs, err := awsutil.QueryNotionConfigs(util.GetOsEnv(util.EnvDdbTable), []string{configId})
	if err != nil {
		return "", err
	}
	for i := 0; i < len(configs); i++ {
		configs[i].NotionToken = ""
	}

	data, err := json.Marshal(configs)
	if err != nil {
		logrus.WithError(err).Error("failed marshal configs to json")
		return "", err
	}

	return string(data), nil
}

func addConfigHandler(request *http.Request) (string, error) {
	return "add", nil
}

func updateConfigHandler(request *http.Request) (string, error) {
	return "update", nil
}

func deleteConfigHandler(request *http.Request) (string, error) {
	return "delete", nil
}
