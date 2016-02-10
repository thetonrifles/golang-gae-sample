package hw

import (
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
)

type HttpResponse struct {
  Status string   `json:"status"`
  Message string  `json:"message"`
}

func init() {
  r := mux.NewRouter()
  r.HandleFunc("/calendar/{id}", auth(GetCalendarHandler)).Methods("GET")
  r.HandleFunc("/calendars", auth(GetCalendarsHandler)).Methods("GET")
  r.HandleFunc("/calendar", auth(PostCalendarHandler)).Methods("POST")
  r.HandleFunc("/event/{calendarId}", auth(PostEventHandler)).Methods("POST")
  http.Handle("/", r)
}

func GetCalendarsHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  calendars := GetCalendars(r, owner)
  for _, calendar := range calendars {
    if calendar.Events == nil {
      calendar.Events = []Event{}
    }
  }
  responseHandler(w, calendars)
}

func GetCalendarHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  vars := mux.Vars(r)
  calendar := GetCalendar(r, owner, vars["id"])
  if calendar != nil {
    if calendar.Events == nil {
      calendar.Events = []Event{}
    }
    responseHandler(w, calendar)
  } else {
    errorHandler(w, r, http.StatusNotFound, "not found")
  }
}

func PostCalendarHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  decoder := json.NewDecoder(r.Body)
  var calendar Calendar
  err := decoder.Decode(&calendar)
  if err == nil {
    calendar.Owner = owner
    if calendar.Events == nil {
      calendar.Events = []Event{}
    }
    success, err := PostCalendar(r, calendar)
    if success {
      responseHandler(w, calendar)
    } else {
      errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
    }
  } else {
    errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
  }
}

func PostEventHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  vars := mux.Vars(r)
  decoder := json.NewDecoder(r.Body)
  var event Event
  err := decoder.Decode(&event)
  if err == nil {
    success, err := PostEvent(r, vars["calendarId"], owner, event)
    if success {
      responseHandler(w, event)
    } else {
      errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
    }
  } else {
    errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
  }
}

func auth(fn http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    owner := r.Header.Get("Authorization")
    if owner == "" {
      errorHandler(w, r, http.StatusUnauthorized, "unauthorized")
    } else {
      fn(w, r)
    }
  }
}

func responseHandler(w http.ResponseWriter, v interface{}) {
  w.Header().Set("Content-Type", "application/json")
  encoder := json.NewEncoder(w)
  encoder.Encode(v)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)
  response := HttpResponse{Status:"failure",Message:message}
  encoder := json.NewEncoder(w)
  encoder.Encode(response)
}
