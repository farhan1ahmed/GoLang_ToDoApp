package tasks

import (
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/auth"
	"net/http"
)

func handleRequests(tApp *TaskApp) {
	http.Handle("/alltasks", auth.AuthMiddleware(http.HandlerFunc(tApp.allTasks)))
	http.Handle("/task", auth.AuthMiddleware(http.HandlerFunc(tApp.singleTask)))
	http.Handle("/create", auth.AuthMiddleware(http.HandlerFunc(tApp.createTask)))
	http.Handle("/delete", auth.AuthMiddleware(http.HandlerFunc(tApp.deleteTask)))
	http.Handle("/update", auth.AuthMiddleware(http.HandlerFunc(tApp.updateTask)))
	http.Handle("/upattach", auth.AuthMiddleware(http.HandlerFunc(tApp.uploadAttachment)))
	http.Handle("/downattach", auth.AuthMiddleware(http.HandlerFunc(tApp.downloadAttachment)))
	http.Handle("/deleteattach", auth.AuthMiddleware(http.HandlerFunc(tApp.deleteAttachment)))
	http.Handle("/similartasks", auth.AuthMiddleware(http.HandlerFunc(tApp.similarTasks)))

}
