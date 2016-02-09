package hw

import (
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
  r.HandleFunc("/calendar", PostCalendarHandler).Methods("POST")
  http.Handle("/", r)
}

func GetCalendarHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  owner := r.Header.Get("Authorization")
  if owner == "" {
    errorHandler(w, r, http.StatusUnauthorized)
  } else {
    calendar := GetCalendar(r, owner, vars["id"])
    if calendar != nil {
      if calendar.Events == nil {
        calendar.Events = []*Event{}
      }
      encoder := json.NewEncoder(w)
      encoder.Encode(calendar)
    } else {
      errorHandler(w, r, http.StatusNotFound)
    }
  }
}

func PostCalendarHandler(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var calendar Calendar
  err := decoder.Decode(&calendar)
  if err == nil {
    if calendar.Events == nil {
      calendar.Events = []*Event{}
    }
    success, _ := PostCalendar(r, calendar)
    if success {
      encoder := json.NewEncoder(w)
      encoder.Encode(calendar)
    } else {
      errorHandler(w, r, http.StatusInternalServerError)
    }
  } else {
    errorHandler(w, r, http.StatusInternalServerError)
  }
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
  w.WriteHeader(status)
  var message string
  if status == http.StatusNotFound {
    message = "not found"
  } else if status == http.StatusUnauthorized {
    message = "unauthorized"
  } else {
    message = "bad request"
  }
  response := HttpResponse{Status:"failure",Message:message}
  encoder := json.NewEncoder(w)
  encoder.Encode(response)
}
