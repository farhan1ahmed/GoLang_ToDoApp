package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type JsonResponse struct {
	MSG string `json:"msg"`
}

func MethodNotAllowed(w http.ResponseWriter) {
	message := JsonResponse{"Method not allowed"}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusMethodNotAllowed)
	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func JSONMsg(w http.ResponseWriter, msg string, code int) {
	message := JsonResponse{msg}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Dateparser(dateISO string) time.Time {
	const layoutISO = "2006-01-02 00:00:00 MST"
	const timeLayout = " 00:00:00 "
	t := time.Now()
	zone, _ := t.Zone()
	dt := dateISO + timeLayout + zone
	date, err := time.Parse(layoutISO, dt)
	if err != nil {
		fmt.Println(err.Error())
	}
	return date
}
