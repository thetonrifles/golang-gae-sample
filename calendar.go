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
	now := time.Now()
	events = append(events, &Event{Title:"first",Timestamp:now.Add(time.Hour)})
	events = append(events, &Event{Title:"second",Timestamp:now.Add(2*time.Hour)})
	events = append(events, &Event{Title:"third",Timestamp:now.Add(3*time.Hour)})
	return events
}