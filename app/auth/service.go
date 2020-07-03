package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/farhan1ahmed/GoLang_ToDoApp/app/utils"
	"net/http"
	"os"
	"strings"
	"time"
)

var jwtSecret string

func init() {
	jwtSecret = os.Getenv("JWT_SECRET")
}

func CreateJWTToken(id uint, username string) (string, error) {
	var err error
	authClaims := jwt.MapClaims{}
	authClaims["id"] = id
	authClaims["userName"] = username
	authClaims["expiry"] = time.Now().Add(time.Minute * 30).Unix()
	auth := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)
	token, err := auth.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return token, nil

}
func GetTokenValue(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	authParsed := strings.Split(authHeader, "Bearer ")
	if len(authParsed) != 2 {
		return "", fmt.Errorf("Access Token Missing")
	}
	tokenVal := authParsed[1]
	return tokenVal, nil
}

func AuthMiddleware(apiEndPoint http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenVal, err := GetTokenValue(r)
		if err != nil {
			utils.JSONMsg(w, err.Error(), http.StatusUnauthorized)
		} else {
			exist := CheckinBlackList(tokenVal)
			if exist {
				utils.JSONMsg(w, "Invalid token", http.StatusUnauthorized)
			} else {
				token, err := jwt.Parse(tokenVal, func(token *jwt.Token) (interface{}, error) {
					return []byte(jwtSecret), nil
				})
				if err != nil {
					utils.JSONMsg(w, err.Error(), http.StatusUnauthorized)
				}
				claims, ok := token.Claims.(jwt.MapClaims)
				if ok && token.Valid {
					ctx := context.WithValue(r.Context(), "id", claims["id"])
					apiEndPoint.ServeHTTP(w, r.WithContext(ctx))
				} else {
					utils.JSONMsg(w, "Invalid token", http.StatusUnauthorized)
				}
			}
		}
	})
}

func CheckinBlackList(val string) bool {
	var token BlackListToken
	tokendb.Where("token_val = ?", val).Find(&token)
	if token.TokenVal != "" {
		return true
	}
	return false
}

func AddToBlackList(tokenVal string) {
	token := BlackListToken{tokenVal}
	tokendb.Create(&token)
}
