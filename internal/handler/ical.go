package handler

import (
	"errors"
	ics "github.com/arran4/golang-ical"
	"github.com/ronazst/notion-ical-syncer/internal/awsutil"
	"github.com/ronazst/notion-ical-syncer/internal/notion"
	"github.com/ronazst/notion-ical-syncer/internal/util"
	"net/http"
	"strings"
)

func iCalHandler(request *http.Request) (string, error) {
	configIds, err := getConfigIds(request)
	if err != nil {
		return "", err
	}

	notionConfigs, err := awsutil.QueryNotionConfigs(util.GetOsEnv(util.EnvDdbTable), configIds)
	if err != nil {
		return "", err
	}

	vEvents, err := notion.QueryEvents(request.Context(), notionConfigs)
	if err != nil {
		return "", err
	}

	cal := ics.NewCalendar()
	cal.SetProductId("notion-ical-syncer")
	cal.SetVersion("2.0")
	for _, vEvent := range vEvents {
		cal.AddVEvent(vEvent)
	}

	return cal.Serialize(), nil
}

func getConfigIds(request *http.Request) ([]string, error) {
	configIds := request.URL.Query().Get(util.QueryArgConfigIds)
	if util.IsBlank(configIds) {
		return nil, errors.New("user request without config_ids")
	}
	return strings.Split(configIds, ","), nil
}
