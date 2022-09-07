package transport

import (
	"encoding/json"
	"io"
	"jwt-demo/util/jwt"
	"log"
	"net/http"
)

type Message struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func Login(w http.ResponseWriter, r *http.Request) {
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

	token, err := jwt.GenerateAccessToken(&userInfo)

	if err != nil {
		log.Println("Error Generating Token")
	}

	output, err := json.Marshal(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}
