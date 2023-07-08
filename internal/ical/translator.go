package ical

import (
	"errors"
	"github.com/arran4/golang-ical"
	"github.com/ronazst/notion-ical-syncer/internal/model"
	"time"
)

func TransToICal(calDataList []model.CalendarData) (string, error) {
	if len(calDataList) == 0 {
		return "", errors.New("empty calendar calDataList")
	}

	cal := ics.NewCalendar()
	cal.SetProductId("notion-ical-syncer")
	for _, data := range calDataList {
		event := cal.AddEvent(data.Id)
		event.SetCreatedTime(time.Now())
		event.SetSummary(data.Title)
		event.SetAllDayStartAt(data.Date)
	}

	return cal.Serialize(), nil
}
