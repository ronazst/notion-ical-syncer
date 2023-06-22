package notion

import (
	"context"
	"errors"
	"github.com/jomei/notionapi"
	"github.com/notion-ical-syncer/internal/model"
	"time"
)

func QueryCalendarData(ctx context.Context, token string, items []model.QueryItem) ([]model.CalendarData, error) {
	var calendarData []model.CalendarData
	date := notionapi.Date(time.Now())
	client := notionapi.NewClient(notionapi.Token(token))

	for _, item := range items {
		response, err := client.Database.Query(ctx, notionapi.DatabaseID(item.DatabaseID), buildQueryRequest(item, date))
		if err != nil {
			return nil, err
		}
		for _, result := range response.Results {
			data, err := buildCalendarData(result)
			if err != nil {
				return nil, err
			}
			calendarData = append(calendarData, *data)
		}
	}

	return calendarData, nil
}

func buildQueryRequest(item model.QueryItem, date notionapi.Date) *notionapi.DatabaseQueryRequest {
	return &notionapi.DatabaseQueryRequest{
		Filter: notionapi.PropertyFilter{
			Property: item.DateFieldKey,
			Date: &notionapi.DateFilterCondition{
				OnOrAfter: &date,
			},
		},
		Sorts: []notionapi.SortObject{{
			Property:  item.DateFieldKey,
			Direction: notionapi.SortOrderASC,
		}},
	}
}

func buildCalendarData(page notionapi.Page) (*model.CalendarData, error) {
	var calData model.CalendarData
	calData.Id = string(page.ID)

	if name, ok := page.Properties["Name"].(*notionapi.TitleProperty); ok {
		calData.Title = name.Title[0].PlainText
	} else {
		return nil, errors.New("failed to get title")
	}
	if date, ok := page.Properties["Date"].(*notionapi.DateProperty); ok {
		parsedTime, err := time.Parse(time.RFC3339, date.Date.Start.String())
		if err != nil {
			return nil, err
		}
		calData.Date = parsedTime
	} else {
		return nil, errors.New("failed to get date")
	}
	return &calData, nil
}
