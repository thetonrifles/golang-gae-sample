package hw

import (
    "fmt"
    "time"
    "net/http"
    "encoding/json"
)

type HttpResponse struct {
    Status string   `json:"status"`
    Message string  `json:"message"`
    Entity Event    `json:"entity"`
}

func init() {
    http.HandleFunc("/events", GetEventsHandler)
    http.HandleFunc("/event", PostEventHandler)
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	events := GetEvents(r)
	json, _ := json.Marshal(events)
  fmt.Fprint(w, string(json))
}

func PostEventHandler(w http.ResponseWriter, r *http.Request) {
  event := Event{Title:"prova",Timestamp:time.Now()}
  success, _ := PutEvent(r, event)
  response := HttpResponse{Status:"success",Entity:event}
  if !success {
    response = HttpResponse{Status:"failure"}
  }
	json, _ := json.Marshal(response)
  fmt.Fprint(w, string(json))
}
