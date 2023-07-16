package model

import (
	"errors"
	"github.com/ronazst/notion-ical-syncer/internal/util"
)

type NotionConfig struct {
	ConfigId         string       `dynamodbav:"config_id" json:"config_id"`
	NotionDbId       string       `dynamodbav:"notion_db_id" json:"notion_db_id"`
	NotionToken      string       `dynamodbav:"notion_token" json:"notion_token"`
	FieldMapping     FieldMapping `dynamodbav:"field_mapping" json:"field_mapping"`
	ExcludeStatusKey string       `dynamodbav:"exclude_status_key" json:"exclude_status_key"`
	ExcludeStatus    []string     `dynamodbav:"exclude_status" json:"exclude_status"`
}

type FieldMapping struct {
	Title       string `dynamodbav:"title" json:"title"`
	Location    string `dynamodbav:"location"  json:"location"`
	Description string `dynamodbav:"description" json:"description"`
	EventTime   string `dynamodbav:"event_time" json:"event_time"`
}

func (n *NotionConfig) Validate() error {
	if util.IsBlank(n.FieldMapping.Title) || util.IsBlank(n.FieldMapping.EventTime) {
		return errors.New("the title and event_time field mapping is required")
	}
	if len(n.ExcludeStatus) != 0 && util.IsBlank(n.ExcludeStatusKey) {
		return errors.New("exclude status key is required when exclude status is specified")
	}
	return nil
}
