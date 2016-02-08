package hw

import (
    "time"
)

type Event struct {
    Title string 		`json:"title"`
    Timestamp time.Time `json:"timestamp"` 
}

func GetEvents() []*Event {
	events := make([]*Event, 0, 10)
	events = append(events, &Event{Title:"first",Timestamp:time.Now()})
	return events
}