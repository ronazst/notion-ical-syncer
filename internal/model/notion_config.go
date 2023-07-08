package model

type NotionConfig struct {
	ConfigId     string       `dynamodbav:"config_id"`
	NotionDbId   string       `dynamodbav:"notion_db_id"`
	NotionToken  string       `dynamodbav:"notion_token"`
	FieldMapping FieldMapping `dynamodbav:"field_mapping"`
}

type FieldMapping struct {
	Summary     string `dynamodbav:"summary"`
	Location    string `dynamodbav:"location"`
	Date        string `dynamodbav:"date"`
	CreatedTime string `dynamodbav:"created_time"`
	Description string `dynamodbav:"description"`
	URL         string `dynamodbav:"url"`
}
