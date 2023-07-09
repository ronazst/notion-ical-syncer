package model

import "time"

type CalendarData struct {
	Id             string
	Title          string
	Description    string
	Date           time.Time
	CreatedTime    time.Time
	LastEditedTime time.Time
}
