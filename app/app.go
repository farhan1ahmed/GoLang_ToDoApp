package app

import (
	"context"
	"fmt"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/auth"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/tasks"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/users"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/utils"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

var db *gorm.DB
var err error

var databaseLoc string
var port string
var START_TIME string
var DELAY_TIME time.Duration
var ctx context.Context

func init() {
	databaseLoc = os.Getenv("SQLITE_DATABASE")
	port = os.Getenv("PORT")
	START_TIME = time.Now().Format("2006-01-02")
	DELAY_TIME = time.Hour * 24   // 1 day
	ctx = context.Background()
}

func Start() {
	db, err = gorm.Open("sqlite3", databaseLoc)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database!")
	}
	defer db.Close()

	taskApp := tasks.TaskApp{db}
	taskApp.InitTodoModel()
	taskApp.InitTaskHandlers()

	userApp := users.UserApp{db}
	userApp.InitUserModel()
	userApp.InitUserHandlers()

	authApp := auth.AuthApp{db}
	authApp.InitBlackListModel()

	go func() {
		for trigger := range utils.CronJob(ctx, utils.Dateparser(START_TIME), DELAY_TIME) {
			fmt.Println(trigger)
			users.ReminderEmail(taskApp.DB, userApp.DB)
		}
	}()
	log.Fatal(http.ListenAndServe(port, nil))
}
