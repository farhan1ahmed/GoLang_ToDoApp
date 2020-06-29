package users

import "net/http"

func handleRequests(uApp *UserApp) {
	http.HandleFunc("/register", uApp.registerUser)
	http.HandleFunc("/confirm/", uApp.confirmUser)
}
