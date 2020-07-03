package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var attachmentFolder string

func init() {
	attachmentFolder = os.Getenv("ATTACHMENT_FOLDER")
}

func (tApp *TaskApp) allTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userID := r.Context().Value("id")
		var tasks []TaskModel
		tApp.DB.Where("user_id = ?", userID).Find(&tasks)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
	} else {
		utils.MethodNotAllowed(w)
	}
}
func (tApp *TaskApp) singleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userID := r.Context().Value("id")
		var task TaskModel
		idVar := r.URL.Query()["id"][0]
		tApp.DB.Where("user_id = ?", userID).First(&task, idVar)

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
		userID := r.Context().Value("id")
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
		reqBody["UserID"] = userID
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
		userID := r.Context().Value("id")
		idVar := r.URL.Query()["id"][0]
		var deltask TaskModel
		tApp.DB.Where("user_id = ?", userID).Find(&deltask, idVar)
		if deltask.ID == 0 {
			utils.JSONMsg(w, "No task found", http.StatusNotFound)
			return
		}
		//Check and delete if there are any attachments for this task
		if deltask.AttachmentName != "" {
			filename := "u" + fmt.Sprintf("%.0f", userID) + "t" + idVar + deltask.AttachmentName
			attachmentPath := attachmentFolder + "/" + filename
			os.Remove(attachmentPath)
		}

		tApp.DB.Unscoped().Delete(TaskModel{}, idVar)
		utils.JSONMsg(w, "Task deleted successfully", http.StatusOK)

	} else {
		utils.MethodNotAllowed(w)
	}
}

func (tApp *TaskApp) updateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPatch {
		userID := r.Context().Value("id")
		idVar := r.URL.Query()["id"][0]
		var updateTask TaskModel
		tApp.DB.Where("user_id = ?", userID).Find(&updateTask, idVar)
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
		userID := r.Context().Value("id")
		idVar := r.URL.Query()["id"][0]
		var taskUploadTo TaskModel
		tApp.DB.Where("user_id = ?", userID).Find(&taskUploadTo, idVar)
		if taskUploadTo.ID == 0 {
			utils.JSONMsg(w, "No task found", http.StatusNotFound)
			return
		}
		file, fh, err := r.FormFile("attachment")
		if err != nil {
			fmt.Println("1")
			utils.JSONMsg(w, err.Error(), http.StatusBadRequest)
			return
		}
		filename := "u" + fmt.Sprintf("%.0f", userID) + "t" + idVar + fh.Filename
		savePath := attachmentFolder + "/" + filename
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
		taskUploadTo.AttachmentName = fh.Filename
		tApp.DB.Save(&taskUploadTo)
		utils.JSONMsg(w, "File uploaded successfully", http.StatusOK)

	} else {
		utils.MethodNotAllowed(w)
	}
}

func (tApp *TaskApp) downloadAttachment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userID := r.Context().Value("id")
		idVar := r.URL.Query()["id"][0]
		var task TaskModel
		tApp.DB.Where("user_id = ?", userID).Find(&task, idVar)
		if task.ID == 0 {
			utils.JSONMsg(w, "No task found", http.StatusNotFound)
			return
		}
		if task.AttachmentName == "" {
			utils.JSONMsg(w, "Task has no attachment", http.StatusNotFound)
			return
		}
		fileCode := "u" + fmt.Sprintf("%.0f", userID) + "t" + idVar
		filepath := attachmentFolder + "/" + fileCode + task.AttachmentName
		readFile, err := os.Open(filepath)
		if err != nil {
			utils.JSONMsg(w, "unable to read file", http.StatusInternalServerError)
			return
		}
		defer readFile.Close()

		fileHeader := make([]byte, 512)
		readFile.Read(fileHeader)
		fileType := http.DetectContentType(fileHeader) // set the type
		fileInfo, _ := readFile.Stat()
		fileSize := fileInfo.Size()

		w.Header().Set("Content-Disposition", "attachment; filename="+task.AttachmentName)
		w.Header().Set("Content-Type", fileType)
		w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		readFile.Seek(0, 0)
		io.Copy(w, readFile)

	} else {
		utils.MethodNotAllowed(w)
	}

}
