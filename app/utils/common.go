package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

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

func CronJob(ctx context.Context, startTriggerTime time.Time, delay time.Duration) <- chan time.Time {
	delayTriggerChan := make(chan time.Time, 1)
	waitFor := time.Until(startTriggerTime)
	//Make sure that the start time is not zero
	if !startTriggerTime.IsZero(){
		// if the starting time has passed already, set the starting time to the next triggering event
		if waitFor < 0 {
			timeTillNext := waitFor - delay
			floorTimes := timeTillNext / delay * - 1
			startTriggerTime = startTriggerTime.Add(floorTimes * delay)
		}
	}

	//Go routine for cron job task
	go func() {
		//Start the scheduled task for the first time
		t := <-time.After(time.Until(startTriggerTime))
		delayTriggerChan <- t

		//After the first triggering event, start a ticker to check for future triggering events
		ticker := time.NewTicker(delay)
		defer ticker.Stop()

		for{
			select{
			case trigger := <-ticker.C:
				delayTriggerChan <- trigger
			case <-ctx.Done():
				close(delayTriggerChan)
				return
			}
		}
	}()
	return delayTriggerChan
}