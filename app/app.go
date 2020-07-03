package app

import (
	"fmt"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/auth"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/tasks"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/users"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"
)

var db *gorm.DB
var err error

var databaseLoc string
var port string

func init() {
	databaseLoc = os.Getenv("SQLITE_DATABASE")
	port = os.Getenv("PORT")
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

	log.Fatal(http.ListenAndServe(port, nil))
}
