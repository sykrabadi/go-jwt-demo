package transport

import (
	"encoding/json"
	"io"
	"jwt-demo/util/jwt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func TestValidate(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	log.Println("Good to go")
	authToken := strings.Replace(authorizationHeader, "Bearer", "", -1)
	log.Println(authToken)
	// TODO: Use jwt.VerifyAccessToken() to verify the access token
	//jwt.VerifyAccessToken()
}

func Login(w http.ResponseWriter, r *http.Request) {
	// read from file
	userFile, err := os.Open("users.json")
	if err != nil {
		log.Println("Error opening users.json")
		return
	}
	defer userFile.Close()

	byteJson, err := io.ReadAll(userFile)
	if err != nil {
		log.Println("Error reading users.json")
		return
	}

	var userDataFromFile map[string]interface{}

	err = json.Unmarshal([]byte(byteJson), &userDataFromFile)
	if err != nil {
		log.Println("Error unmarshall users.json")
		return
	}
	// read request body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userInfo jwt.UserForToken

	err = json.Unmarshal(b, &userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO : Fix checking mechanism
	if (userDataFromFile["email"] != userInfo.UserEmail) || (userDataFromFile["password"] != userInfo.Password) {
		log.Println("Email from file and request are unmatch")
		http.Error(w, "Email from file and request are unmatch", http.StatusBadRequest)
		return
	}

	accessToken, err := jwt.GenerateAccessToken(&userInfo)

	if err != nil {
		log.Println("Error Generating Access Token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(&userInfo)

	if err != nil {
		log.Println("Error Generating Refresh Token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := jwt.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	response, err := json.Marshal(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(response)
}
