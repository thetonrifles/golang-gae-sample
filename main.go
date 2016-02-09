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
  r.HandleFunc("/calendar", PostCalendarHandler).Methods("POST")
  http.Handle("/", r)
}

func GetCalendarHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  calendar := GetCalendar(r, vars["id"])
  if calendar != nil {
    json, _ := json.Marshal(calendar)
    fmt.Fprint(w, string(json))
  } else {
    errorHandler(w, r, http.StatusNotFound)
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
      json, _ := json.Marshal(calendar)
      fmt.Fprint(w, string(json))
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
  } else {
    message = "bad request"
  }
  response := HttpResponse{Status:"failure",Message:message}
  json, _ := json.Marshal(response)
  fmt.Fprint(w, string(json))
}
