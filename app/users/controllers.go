package users

import (
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/auth"
	"net/http"
)

func handleRequests(uApp *UserApp) {
	http.HandleFunc("/register", uApp.registerUser)
	http.HandleFunc("/confirm/", uApp.confirmUser)
	http.HandleFunc("/login", uApp.loginUser)
	http.Handle("/logout", auth.AuthMiddleware(http.HandlerFunc(uApp.logoutUser)))
}
