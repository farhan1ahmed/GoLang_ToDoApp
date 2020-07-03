package tasks

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type TaskApp struct {
	DB *gorm.DB
}

type TaskModel struct {
	gorm.Model
	Title          string    `gorm:"unique_index:uniqueforuser; not null"`
	Description    string    `gorm:"not null"`
	Finished       bool      `gorm:"not null; default:'0'"`
	DueDate        time.Time `gorm:"not null"`
	CompletionDate time.Time `gorm:"default:null"`
	AttachmentName string
	UserID         int `gorm:"unique_index:uniqueforuser; not null"`
}

func (tApp *TaskApp) InitTodoModel() {
	tApp.DB.AutoMigrate(&TaskModel{})
}

func (tApp *TaskApp) InitTaskHandlers() {
	handleRequests(tApp)
}
