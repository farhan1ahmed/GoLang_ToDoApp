package app

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"
	"toDoApp/app/tasks"
	"toDoApp/app/users"
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

	tasksDB := tasks.TaskApp{db}
	tasksDB.InitTodoModel()
	usersDB := users.UserApp{db}
	usersDB.InitUserModel()
	log.Fatal(http.ListenAndServe(port, nil))
}
