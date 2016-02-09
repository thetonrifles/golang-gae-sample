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
  r.HandleFunc("/calendar/{id}", GetCalendarHandler).Methods("GET")
  r.HandleFunc("/calendars", GetCalendarsHandler).Methods("GET")
  r.HandleFunc("/calendar", PostCalendarHandler).Methods("POST")
  r.HandleFunc("/event/{calendarId}", PostEventHandler).Methods("POST")
  http.Handle("/", r)
}

func GetCalendarsHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  if owner == "" {
    errorHandler(w, r, http.StatusUnauthorized, "unauthorized")
  } else {
    calendars := GetCalendars(r, owner)
    for _, calendar := range calendars {
      if calendar.Events == nil {
        calendar.Events = []*Event{}
      }
    }
    encoder := json.NewEncoder(w)
    encoder.Encode(calendars)
  }
}

func GetCalendarHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  if owner == "" {
    errorHandler(w, r, http.StatusUnauthorized, "unauthorized")
  } else {
    vars := mux.Vars(r)
    calendar := GetCalendar(r, owner, vars["id"])
    if calendar != nil {
      if calendar.Events == nil {
        calendar.Events = []*Event{}
      }
      encoder := json.NewEncoder(w)
      encoder.Encode(calendar)
    } else {
      errorHandler(w, r, http.StatusNotFound, "not found")
    }
  }
}

func PostCalendarHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  if owner == "" {
    errorHandler(w, r, http.StatusUnauthorized, "unauthorized")
  } else {
    decoder := json.NewDecoder(r.Body)
    var calendar Calendar
    err := decoder.Decode(&calendar)
    if err == nil {
      calendar.Owner = owner
      if calendar.Events == nil {
        calendar.Events = []*Event{}
      }
      success, err := PostCalendar(r, calendar)
      if success {
        encoder := json.NewEncoder(w)
        encoder.Encode(calendar)
      } else {
        errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
      }
    } else {
      errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
    }
  }
}

func PostEventHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  if owner == "" {
    errorHandler(w, r, http.StatusUnauthorized, "unauthorized")
  } else {
    vars := mux.Vars(r)
    decoder := json.NewDecoder(r.Body)
    var event Event
    err := decoder.Decode(&event)
    if err == nil {
      success, err := PostEvent(r, vars["calendarId"], owner, event)
      if success {
        encoder := json.NewEncoder(w)
        encoder.Encode(event)
      } else {
        errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
      }
    } else {
      errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
    }
  }
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
  w.WriteHeader(status)
  response := HttpResponse{Status:"failure",Message:message}
  encoder := json.NewEncoder(w)
  encoder.Encode(response)
}
