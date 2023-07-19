package handler

import (
	"encoding/json"
	"github.com/ronazst/notion-ical-syncer/internal/awsutil"
	"github.com/ronazst/notion-ical-syncer/internal/model"
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

	config, err := awsutil.QueryNotionConfigs(util.GetOsEnv(util.EnvDdbTable), []string{configId})
	if err != nil {
		return "", err
	}
	if len(config) == 0 {
		return "", util.NewUserInputError("config not found")
	}
	config[0].NotionToken = ""

	data, err := json.Marshal(config[0])
	if err != nil {
		logrus.WithError(err).Error("failed marshal config to json")
		return "", err
	}

	return string(data), nil
}

func changeConfigHandler(request *http.Request) (string, error) {
	config := model.NotionConfig{}
	err := json.NewDecoder(request.Body).Decode(&config)
	if err != nil {
		logrus.WithError(err).Error("failed decode request body to config")
		return "", util.NewUserInputError("failed decode request body to config")
	}

	switch request.Method {
	case http.MethodPost:
		if err = awsutil.AddOrUpdateNotionConfig(util.GetOsEnv(util.EnvDdbTable), config); err != nil {
			return "", err
		}
	case http.MethodDelete:
		if err = awsutil.DeleteNotionConfig(util.GetOsEnv(util.EnvDdbTable), config.ConfigId); err != nil {
			return "", err
		}
	}

	return "", nil
}
