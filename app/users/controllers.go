package users

import "net/http"

func handleRequests(uApp *UserApp) {
	http.HandleFunc("/register", uApp.registerUser)
	http.HandleFunc("/confirm/", uApp.confirmUser)
	http.HandleFunc("/login", uApp.loginUser)
	http.HandleFunc("/logout", uApp.logoutUser)
}
