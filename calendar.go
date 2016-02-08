package hw

import (
    "time"
    "net/http"
    "google.golang.org/appengine"
    "google.golang.org/appengine/datastore"
)

type Calendar struct {
  Id string       `json:"id"`
  Owner string    `json:"owner"`
  Events []*Event `json:"events"`
}

type Event struct {
  Title string 		    `json:"title"`
  Timestamp time.Time `json:"timestamp"`
}

func GetCalendar(r *http.Request, id string) *Calendar {
  context := appengine.NewContext(r)
  q := datastore.NewQuery("calendar").Filter("Id=",id)
  var calendar Calendar
  keys, _ := q.GetAll(context, &calendar)
  if len(keys) == 0 {
    return nil
  }
	return &calendar
}

func PostCalendar(r *http.Request, calendar Calendar) (bool, error) {
  context := appengine.NewContext(r)
  _, err := datastore.Put(context, datastore.NewIncompleteKey(context, "calendar", nil), &calendar)
  if err != nil {
    return false, err
  } else {
    return true, nil
  }
}

/*
func PutEvent(r *http.Request, calendarId string, event Event) (bool, error) {
  context := appengine.NewContext(r)
  _, err := datastore.Put(context, datastore.NewIncompleteKey(context, "event", nil), &event)
  if err != nil {
    return false, err
  } else {
    return true, nil
  }
}
*/
