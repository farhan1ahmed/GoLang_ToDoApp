package users

import (
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/tasks"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type UserApp struct {
	DB *gorm.DB
}

type UserModel struct {
	gorm.Model
	UserName  string            `gorm:"unique; notnull"`
	Email     string            `gorm:"unique; notnull"`
	Password  string            `gorm:"notnull"`
	FBuser    bool              `gorm:"default:'0'"`
	Confirmed bool              `gorm:"default:'0'"`
	Tasks     []tasks.TaskModel `gorm:"ForeignKey:UserID"`
}

func (uApp *UserApp) InitUserModel() {
	uApp.DB.AutoMigrate(&UserModel{})
}
func (uApp *UserApp) InitUserHandlers() {
	handleRequests(uApp)
}
