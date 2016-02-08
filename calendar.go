package hw

import (
    "time"
    "net/http"
    "google.golang.org/appengine"
    "google.golang.org/appengine/datastore"
)

type Event struct {
    Title string 		`json:"title"`
    Timestamp time.Time `json:"timestamp"`
}

func GetEvents(r *http.Request) []Event {
  context := appengine.NewContext(r)
  q := datastore.NewQuery("event")
  var events []Event
  q.GetAll(context, &events)
	return events
}

func PutEvent(r *http.Request, event Event) (bool, error) {
  context := appengine.NewContext(r)
  _, err := datastore.Put(context, datastore.NewIncompleteKey(context, "event", nil), &event)
  if err != nil {
    return false, err
  } else {
    return true, nil
  }
}
