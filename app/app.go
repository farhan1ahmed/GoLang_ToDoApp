package app

import (
	"fmt"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/auth"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/tasks"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/users"
	"github.com/go-co-op/gocron"
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
var scheduler *gocron.Scheduler

func init() {
	databaseLoc = os.Getenv("SQLITE_DATABASE")
	port = os.Getenv("PORT")
	scheduler = gocron.NewScheduler(time.Now().Location())
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

	scheduler.Every(1).Day().At("17:14:00").Do(users.ReminderEmail, taskApp.DB, userApp.DB)
	scheduler.StartAsync()

	authApp := auth.AuthApp{db}
	authApp.InitBlackListModel()

	log.Fatal(http.ListenAndServe(port, nil))
}
