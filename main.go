package hw

import (
    "fmt"
    "net/http"
    "encoding/json"
)

func init() {
    http.HandleFunc("/events", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	events := GetEvents()
	json, _ := json.Marshal(events)
    fmt.Fprint(w, string(json))
}