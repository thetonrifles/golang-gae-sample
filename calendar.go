package hw

import (
  "fmt"
  "errors"
  "net/http"
  "crypto/md5"
  "google.golang.org/appengine"
  "google.golang.org/appengine/datastore"
)

type Calendar struct {
  Id string       `json:"id"`
  Owner string    `json:"owner"`
  Events []Event `json:"events"`
}

type Event struct {
  Title string 		`json:"title"`
  Timestamp int64 `json:"timestamp"`
}

func GetCalendars(r *http.Request, owner string) []*Calendar {
  context := appengine.NewContext(r)
  q := datastore.NewQuery("calendar").Filter("Owner =", owner)
  var calendars []*Calendar
  _, err := q.GetAll(context, &calendars)
  if err==nil && calendars!=nil {
    return calendars
  } else {
    return []*Calendar{}
  }
}

func GetCalendar(r *http.Request, owner string, calendarId string) *Calendar {
  context := appengine.NewContext(r)
  key := hash(owner + calendarId)
  calendarKey := datastore.NewKey(context, "calendar", key, 0, nil)
  var calendar Calendar
  datastore.Get(context, calendarKey, &calendar)
	return &calendar
}

func PostCalendar(r *http.Request, calendar Calendar) (bool, error) {
  context := appengine.NewContext(r)
  key := hash(calendar.Owner + calendar.Id)
  calendarKey := datastore.NewKey(context, "calendar", key, 0, nil)
  err := datastore.Get(context, calendarKey, &calendar)
  if err != nil {
    _, err := datastore.Put(context, calendarKey, &calendar)
    if err != nil {
      return false, err
    } else {
      return true, nil
    }
  } else {
    return false, errors.New("calendar already exists")
  }
}

func PostEvent(r *http.Request, calendarId string, owner string, event Event) (bool, error) {
  context := appengine.NewContext(r)
  key := hash(owner + calendarId)
  calendarKey := datastore.NewKey(context, "calendar", key, 0, nil)
  var calendar Calendar
  err := datastore.Get(context, calendarKey, &calendar)
  if err != nil {
    return false, errors.New("calendar doesn't exists")
  } else {
    calendar.Events = append(calendar.Events, event)
    _, err := datastore.Put(context, calendarKey, &calendar)
    if err != nil {
      return false, err
    } else {
      return true, nil
    }
  }
}

func hash(s string) string {
  data := []byte(s)
  return fmt.Sprintf("%x", md5.Sum(data))
}
