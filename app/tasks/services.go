package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"toDoApp/app/utils"
)

var attachmentFolder string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	attachmentFolder = os.Getenv("ATTACHMENT_FOLDER")
}

func (tApp *TaskApp) allTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var tasks []TaskModel
		tApp.DB.Find(&tasks)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
	} else {
		utils.MethodNotAllowed(w)
	}
}
func (tApp *TaskApp) singleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var task TaskModel
		idVar := r.URL.Query()["id"][0]
		tApp.DB.First(&task, idVar)

		if task.ID == 0 {
			utils.JSONMsg(w, "No task found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	} else {
		utils.MethodNotAllowed(w)
	}
}

func (tApp *TaskApp) createTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		newTask := &TaskModel{}

		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			utils.JSONMsg(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		date, ok := reqBody["DueDate"]
		if ok == false {
			utils.JSONMsg(w, "DueDate not specified", http.StatusBadRequest)
			return
		}
		reqBody["DueDate"] = utils.Dateparser(date.(string))
		jsonBody, _ := json.Marshal(reqBody)
		err = json.Unmarshal(jsonBody, &newTask)
		if err != nil {
			utils.JSONMsg(w, err.Error(), http.StatusBadRequest)
			return
		}

		exc := tApp.DB.Create(&newTask)
		if exc.Error != nil {
			utils.JSONMsg(w, exc.Error.Error(), http.StatusConflict)
			return
		}
		utils.JSONMsg(w, "Task created successfully", http.StatusCreated)

	} else {
		utils.MethodNotAllowed(w)
	}
}

func (tApp *TaskApp) deleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		idVar := r.URL.Query()["id"][0]
		var deltask TaskModel
		tApp.DB.Find(&deltask, idVar)
		if deltask.ID == 0 {
			utils.JSONMsg(w, "No task found", http.StatusNotFound)
			return
		}

		tApp.DB.Unscoped().Delete(TaskModel{}, idVar)
		utils.JSONMsg(w, "Task deleted successfully", http.StatusOK)

	} else {
		utils.MethodNotAllowed(w)
	}
}

func (tApp *TaskApp) updateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPatch {
		idVar := r.URL.Query()["id"][0]
		var updateTask TaskModel
		tApp.DB.Find(&updateTask, idVar)
		if updateTask.ID == 0 {
			utils.JSONMsg(w, "No task found", http.StatusNotFound)
			return
		}
		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			utils.JSONMsg(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		tApp.DB.Model(&updateTask).Updates(reqBody)
		utils.JSONMsg(w, "Task updated successfully", http.StatusOK)
	} else {
		utils.MethodNotAllowed(w)
	}
}

func (tApp *TaskApp) uploadAttachment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		idVar := r.URL.Query()["id"][0]
		var taskToUpdate TaskModel
		tApp.DB.Find(&taskToUpdate, idVar)
		if taskToUpdate.ID == 0 {
			utils.JSONMsg(w, "No task found", http.StatusNotFound)
			return
		}
		file, fh, err := r.FormFile("attachment")
		if err != nil {
			fmt.Println("1")
			utils.JSONMsg(w, err.Error(), http.StatusBadRequest)
			return
		}
		savePath := attachmentFolder + "/" + fh.Filename
		fmt.Println(savePath)
		var _, e = os.Stat(savePath)
		if !os.IsNotExist(e) {
			os.Remove(savePath)
		}
		newFile, err := os.Create(savePath)
		if err != nil {
			fmt.Println("2")
			utils.JSONMsg(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer newFile.Close()

		content, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("3")
			utils.JSONMsg(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		newFile.Write(content)
		taskToUpdate.AttachmentName = fh.Filename
		tApp.DB.Save(&taskToUpdate)
		utils.JSONMsg(w, "File uploaded successfully", http.StatusOK)

	} else {
		utils.MethodNotAllowed(w)
	}
}
