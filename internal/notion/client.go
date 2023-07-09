package notion

import (
	"context"
	"fmt"
	ics "github.com/arran4/golang-ical"
	"github.com/jomei/notionapi"
	"github.com/ronazst/notion-ical-syncer/internal/model"
	"github.com/ronazst/notion-ical-syncer/internal/util"
	"github.com/sirupsen/logrus"
	"time"
)

func QueryEvents(ctx context.Context, configs []model.NotionConfig) ([]*ics.VEvent, error) {
	var events []*ics.VEvent
	date := notionapi.Date(time.Now())

	for _, config := range configs {
		logger := logrus.WithField("config_id", config.ConfigId).WithField("notion_db_id", config.NotionDbId)
		logger.Info("Start to query calendar data with notion config")

		client := notionapi.NewClient(notionapi.Token(config.NotionToken))
		response, err := client.Database.Query(ctx, notionapi.DatabaseID(config.NotionDbId), buildQueryRequest(config.FieldMapping, date))
		if err != nil {
			logger.WithError(err).Error("Failed to query calendar data with notion config")
			return nil, err
		}
		for _, result := range response.Results {
			event, err := buildIcsVEvent(result, config.FieldMapping)
			if err != nil {
				return nil, err
			}
			events = append(events, event)
		}
	}

	return events, nil
}

func buildQueryRequest(mapping model.FieldMapping, date notionapi.Date) *notionapi.DatabaseQueryRequest {
	return &notionapi.DatabaseQueryRequest{
		Filter: notionapi.PropertyFilter{
			Property: mapping.EventTime,
			Date: &notionapi.DateFilterCondition{
				OnOrAfter: &date,
			},
		},
		Sorts: []notionapi.SortObject{{
			Property:  mapping.EventTime,
			Direction: notionapi.SortOrderASC,
		}},
	}
}

func buildIcsVEvent(page notionapi.Page, mapping model.FieldMapping) (*ics.VEvent, error) {
	event := ics.NewEvent(page.ID.String())
	event.SetCreatedTime(page.CreatedTime)
	event.SetModifiedAt(page.LastEditedTime)
	event.SetURL(page.URL)

	if prop, ok := page.Properties[mapping.Title]; ok && prop.GetType() == "title" {
		event.SetSummary(prop.(*notionapi.TitleProperty).Title[0].PlainText)
	} else {
		return nil, util.NewInternalError("Failed to get title, please check your title field mapping")
	}

	if value, err := getTextPropValue(page.Properties, mapping.Title); err != nil {
		return nil, err
	} else {
		event.SetSummary(value)
	}

	if err := setVEventTime(event, page.Properties, mapping.EventTime); err != nil {
		return nil, err
	}

	if !util.IsBlank(mapping.Description) {
		value, err := getTextPropValue(page.Properties, mapping.Description)
		if err != nil {
			return nil, err
		}
		event.SetDescription(value)
	}
	if !util.IsBlank(mapping.Location) {
		value, err := getTextPropValue(page.Properties, mapping.Location)
		if err != nil {
			return nil, err
		}
		event.SetLocation(value)
	}

	return event, nil
}

func setVEventTime(event *ics.VEvent, properties notionapi.Properties, key string) error {
	prop, ok := properties[key]
	if !ok || prop.GetType() != "date" {
		return util.NewInternalError("Failed to get event time, please check your event time field mapping")
	}

	startTime := (*time.Time)(prop.(*notionapi.DateProperty).Date.Start)
	endTime := (*time.Time)(prop.(*notionapi.DateProperty).Date.End)
	if startTime != nil && endTime != nil {
		event.SetStartAt(*startTime)
		event.SetEndAt(*endTime)
	} else if startTime != nil {
		event.SetAllDayStartAt(*startTime)
	} else {
		event.SetAllDayEndAt(*endTime)
	}

	return nil
}

func getTextPropValue(properties notionapi.Properties, key string) (string, error) {
	prop, ok := properties[key]
	if !ok {
		return "", util.NewInternalError(fmt.Sprintf("The field %s does not in notion properties", key))
	}
	switch v := prop.(type) {
	case *notionapi.RichTextProperty:
		return getRichText(v.RichText), nil
	case *notionapi.TitleProperty:
		return getRichText(v.Title), nil
	default:
		return "", util.NewInternalError(fmt.Sprintf("The type %s does not support yet", prop.GetType()))
	}
}

func getRichText(text []notionapi.RichText) string {
	if len(text) == 0 {
		return ""
	}
	return text[0].PlainText
}
