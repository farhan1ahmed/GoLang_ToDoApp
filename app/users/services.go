package users

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"toDoApp/app/utils"
)

var confirmURL string
var senderMail string
var senderPass string
var smtpAddr string
var smtpHost string

func init() {
	confirmURL = os.Getenv("CONFIRM_URL")
	senderMail = os.Getenv("SENDER_MAIL")
	senderPass = os.Getenv("SENDER_PASS")
	smtpAddr = os.Getenv("SMTP_ADDR")
	smtpHost = os.Getenv("SMTP_HOST")
}

func generateConfirmURL(email string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(email))
	return confirmURL + "/" + encoded
}
func sendConfirmEmail(email string) {
	from := fmt.Sprintf("From: <%s>\r\n", senderMail)
	to := fmt.Sprintf("From: <%s>\r\n", email)
	subject := "Subject: Confirm Todo Account\r\n"
	body := "Confirm: " + generateConfirmURL(email) + "?confirm=1" +
		"\nNot you? Click: " + generateConfirmURL(email) + "?confirm=0"
	msg := from + to + subject + "\r\n" + body

	auth := smtp.PlainAuth("", senderMail, senderPass, smtpHost)
	err := smtp.SendMail(smtpAddr, auth, senderMail, []string{email}, []byte(msg))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (uApp *UserApp) registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newUser UserModel
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newUser)
		if err != nil {
			utils.JSONMsg(w, err.Error(), http.StatusBadRequest)
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
		utils.JSONMsg(w, "User created successfully. Confirmation email sent", http.StatusCreated)
		sendConfirmEmail(newUser.Email)
	} else {
		utils.MethodNotAllowed(w)
	}
}
func decodeURL(url string) string {
	re := regexp.MustCompile(`/confirm/(.*?)\?confirm=.$`)
	encodedPart := re.FindStringSubmatch(url)[1]
	decoded, err := base64.URLEncoding.DecodeString(encodedPart)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(decoded)
}
func (uApp *UserApp) confirmUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		email := decodeURL(r.URL.String())
		var user UserModel
		uApp.DB.Unscoped().Where("email = ?", email).Find(&user)
		if user.ID == 0 {
			utils.JSONMsg(w, "No such user exists in database", http.StatusOK)
			return
		}
		confirm := r.URL.Query()["confirm"][0]
		if confirm == "0" {
			uApp.DB.Unscoped().Where("email = ?", email).Delete(&user)
			utils.JSONMsg(w, "User deleted", http.StatusOK)
			return
		}
		uApp.DB.Model(&user).Where("email = ?", email).Update("confirmed", true)
		utils.JSONMsg(w, "User confirmed successfully", http.StatusOK)
	} else {
		utils.MethodNotAllowed(w)
	}
}