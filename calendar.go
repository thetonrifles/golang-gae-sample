package hw

import (
  "time"
  "net/http"
  "hash/fnv"
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

func GetCalendar(r *http.Request, owner string, calendarId string) *Calendar {
  context := appengine.NewContext(r)

  key := hash(owner + calendarId)
  calendarKey := datastore.NewKey(context, "calendar", key, 0, nil)
  var calendar Calendar
  datastore.Get(context, calendarKey, &calendar)

  //q := datastore.NewQuery("calendar").Filter("Id=",id)
  //var calendar Calendar
  //keys, _ := q.GetAll(context, &calendar)
  //log.Println(keys)
  //if len(keys) == 0 {
  //  return nil
  //}
	return &calendar
}

func PostCalendar(r *http.Request, calendar Calendar) (bool, error) {
  context := appengine.NewContext(r)
  key := hash(calendar.Owner + calendar.Id)
  _, err := datastore.Put(context, datastore.NewKey(context, "calendar", key, 0, nil), &calendar)
  if err != nil {
    return false, err
  } else {
    return true, nil
  }
}

func hash(s string) string {
  h := fnv.New32a()
  h.Write([]byte(s))
  return string(h.Sum32())
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
