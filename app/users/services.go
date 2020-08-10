package users

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/auth"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"sync"
	"time"
)

var confirmURL string
var resetURL string
var senderMail string
var senderPass string
var smtpAddr string
var smtpHost string
var jwtSecret string

func init() {
	confirmURL = os.Getenv("CONFIRM_URL")
	resetURL = os.Getenv("RESET_URL")
	senderMail = os.Getenv("SENDER_MAIL")
	senderPass = os.Getenv("SENDER_PASS")
	smtpAddr = os.Getenv("SMTP_ADDR")
	smtpHost = os.Getenv("SMTP_HOST")
	jwtSecret = os.Getenv("JWT_SECRET")
}

func generateEncodedURL(email string, endpoint string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(email))
	return endpoint + "/" + encoded
}
func sendEmail(toEmail string, msg string, wg *sync.WaitGroup) {
	from := fmt.Sprintf("From: <%s>\r\n", senderMail)
	to := fmt.Sprintf("From: <%s>\r\n", toEmail)
	emailBody := from + to + msg

	auth := smtp.PlainAuth("", senderMail, senderPass, smtpHost)
	err := smtp.SendMail(smtpAddr, auth, senderMail, []string{toEmail}, []byte(emailBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	wg.Done()
}

func (uApp *UserApp) registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newUser UserModel
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newUser)
		if err != nil {
			utils.JSONMsg(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		defer r.Body.Close()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		newUser.Password = string(hashedPassword)

		exc := uApp.DB.Create(&newUser)
		if exc.Error != nil {
			utils.JSONMsg(w, exc.Error.Error(), http.StatusConflict)
			return
		}

		subject := "Subject: Confirm Todo Account\r\n"
		link := generateEncodedURL(newUser.Email, confirmURL)
		body := "Confirm: " + link + "?confirm=1" +
			"\nNot you? Click: " + link + "?confirm=0"
		msg := subject + "\r\n" + body

		var wg sync.WaitGroup
		wg.Add(1)
		go sendEmail(newUser.Email, msg, &wg)
		wg.Wait()
		utils.JSONMsg(w, "User created successfully. Confirmation email sent", http.StatusCreated)

	} else {
		utils.MethodNotAllowed(w)
	}
}
func decodeURL(url string, endpoint string) string {
	var re *regexp.Regexp
	switch endpoint {
	case confirmURL:
		re = regexp.MustCompile(`/confirm/(.*?)\?confirm=.$`)
	case resetURL:
		re = regexp.MustCompile(`/resetpassword/(.*?)$`)
	}
	encodedPart := re.FindStringSubmatch(url)[1]
	decoded, err := base64.URLEncoding.DecodeString(encodedPart)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(decoded)
}
func (uApp *UserApp) confirmUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		email := decodeURL(r.URL.String(), confirmURL)
		var user UserModel
		uApp.DB.Where("email = ?", email).Find(&user)
		if user.ID == 0 {
			utils.JSONMsg(w, "No such user exists in database", http.StatusNotFound)
			return
		}
		confirm := r.URL.Query().Get("confirm")
		if confirm == "0" {
			uApp.DB.Unscoped().Where("email = ?", email).Delete(&user)
			utils.JSONMsg(w, "User deleted", http.StatusOK)
			return
		}
		if user.Confirmed == false {
			uApp.DB.Model(&user).Where("email = ?", email).Update("confirmed", true)
			utils.JSONMsg(w, "User confirmed successfully", http.StatusOK)
		} else {
			utils.JSONMsg(w, "User already confirmed", http.StatusOK)
		}
	} else {
		utils.MethodNotAllowed(w)
	}
}

func (uApp *UserApp) loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var reqBody UserModel
		var loginUser UserModel
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			utils.JSONMsg(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		err = uApp.DB.Where("email = ?", reqBody.Email).Find(&loginUser).Error
		if err != nil {
			utils.JSONMsg(w, "No user with given email", http.StatusUnauthorized)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(reqBody.Password))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			utils.JSONMsg(w, "Invalid Password", http.StatusUnauthorized)
			return
		}
		token, err := auth.CreateJWTToken(loginUser.ID, loginUser.UserName)
		if err != nil {
			utils.JSONMsg(w, err.Error(), http.StatusInternalServerError)
			return
		}
		accessCookie := http.Cookie{Name: "accessCookie", Value: token}
		http.SetCookie(w, &accessCookie)
		utils.JSONMsg(w, "User logged in successfully", http.StatusOK)

	} else {
		utils.MethodNotAllowed(w)
	}
}

func (uApp *UserApp) logoutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tokenVal, err := auth.GetTokenValue(r)
		if err != nil {
			utils.JSONMsg(w, err.Error(), http.StatusBadRequest)
			return
		}
		auth.AddToBlackList(tokenVal)
		accessCookie := http.Cookie{Name: "accessCookie", Value: "", MaxAge: -1}
		http.SetCookie(w, &accessCookie)
		utils.JSONMsg(w, "User logged out successfully", http.StatusOK)

	} else {
		utils.MethodNotAllowed(w)
	}
}

func (uApp *UserApp) forgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var reqBody map[string]string
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			utils.JSONMsg(w, err.Error(), http.StatusBadRequest)
			return
		}
		var user UserModel
		email := reqBody["email"]
		uApp.DB.Where("email=?", email).Find(&user)
		if user.ID == 0 {
			utils.JSONMsg(w, "No such email exists in the database", http.StatusBadRequest)
			return
		}
		if user.FBuser {
			utils.JSONMsg(w, "Can not change password: This is an OAuth registered account", http.StatusBadRequest)
			return
		}

		subject := "Subject: Reset Password\r\n"
		link := generateEncodedURL(email, resetURL)
		body := "Follow the link to reset Password: " + link
		msg := subject + "\r\n" + body

		var wg sync.WaitGroup
		wg.Add(1)
		go sendEmail(email, msg, &wg)
		wg.Wait()
		utils.JSONMsg(w, "Link to reset password sent on the provided email.", http.StatusOK)

	} else {
		utils.MethodNotAllowed(w)
	}
}

func (uApp *UserApp) resetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := decodeURL(r.URL.String(), resetURL)
		var user UserModel
		uApp.DB.Where("email = ?", email).Find(&user)
		if user.ID == 0 {
			utils.JSONMsg(w, "No such user exists in database", http.StatusNotFound)
			return
		}
		var reqBody map[string]string
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			utils.JSONMsg(w, "Password not provided", http.StatusBadRequest)
			return
		}
		newPassword := reqBody["password"]
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
		uApp.DB.Save(&user)
		utils.JSONMsg(w, "Password changed successfully", http.StatusOK)
	} else {
		utils.MethodNotAllowed(w)
	}
}

func ReminderEmail(taskAppDB *gorm.DB, userAppDB *gorm.DB) {
	var noOfThreadsSem = make(chan int, 3)
	users, _ := userAppDB.Table("user_models").Select("id, email").Rows()
	defer users.Close()
	var id int
	var email string
	var title string
	dateTomorrow := time.Now().Add(time.Hour*24).Format("2006-01-02")
	var wg sync.WaitGroup
	for users.Next() {
		users.Scan(&id, &email)
		tasks, err := taskAppDB.Table("task_models").Select("title").Where("user_id=? AND finished=? AND due_date LIKE ?", id, false, "%"+dateTomorrow+"%").Rows()
		if err != nil {
			fmt.Println(err.Error())
		}
		subject := "Subject: ToDo Reminder\r\n"
		body := "Following tasks are due today: \n"
		msg := subject + "\r\n"
		taskExists := false
		for tasks.Next(){
			taskExists = true
			tasks.Scan(&title)
			body = body + title + "\n"
		if taskExists {
			wg.Add(1)
			noOfThreadsSem <- 1
			go func() {
				msg = msg + body
				fmt.Println("Sending reminder Email to " + email)
				sendEmail(email, msg, &wg)
				<- noOfThreadsSem
			}()
			}
		}
	wg.Wait()
}}
