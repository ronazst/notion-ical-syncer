package model

import "github.com/ronazst/notion-ical-syncer/internal/util"

type NotionConfig struct {
	ConfigId     string       `dynamodbav:"config_id"`
	NotionDbId   string       `dynamodbav:"notion_db_id"`
	NotionToken  string       `dynamodbav:"notion_token"`
	FieldMapping FieldMapping `dynamodbav:"field_mapping"`
}

type FieldMapping struct {
	Title       string `dynamodbav:"title"`
	Location    string `dynamodbav:"location"`
	Description string `dynamodbav:"description"`
	EventTime   string `dynamodbav:"event_time"`
}

func (n *NotionConfig) Validate() error {
	if util.IsBlank(n.FieldMapping.Title) || util.IsBlank(n.FieldMapping.EventTime) {
		return util.NewInternalError("The title and event_time field mapping is required")
	}
	return nil
}
