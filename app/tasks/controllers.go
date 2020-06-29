package tasks

import (
	"net/http"
)

func handleRequests(tApp *TaskApp) {
	http.HandleFunc("/alltasks", tApp.allTasks)
	http.HandleFunc("/task", tApp.singleTask)
	http.HandleFunc("/create", tApp.createTask)
	http.HandleFunc("/delete", tApp.deleteTask)
	http.HandleFunc("/update", tApp.updateTask)
	http.HandleFunc("/upattach", tApp.uploadAttachment)
	//http.HandleFunc("/downloadAttach", tApp.updateTask)

}
