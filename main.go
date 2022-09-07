package main

import (
	"jwt-demo/transport"
	"log"
	"net/http"

	jwtalias "github.com/golang-jwt/jwt/v4"
)

const (
	SERVER_ADDR = ":8080"
)

var JWT_SIGNING_METHOD = jwtalias.SigningMethodHS256

func main() {

	// user := jwt.UserForToken{
	// 	UserEmail: "test@mail.com",
	// }

	// result, err := tokenManager.GenerateAccessToken(&user)

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println(result)

	http.HandleFunc("/login", transport.Login)
	server := new(http.Server)
	server.Addr = SERVER_ADDR
	log.Println("Starting server at", SERVER_ADDR)
	server.ListenAndServe()
}
