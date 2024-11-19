package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type user struct{
	ID uint32 `json:"id"`
	Nome string `json:"nome"`
	Email string `json:"email"`
}

// Creates an user and persist in the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil{
		w.Write([]byte("Fail do read the request body!"))
		return
	}
	var user user
	if err = json.Unmarshal(requestBody, &user); err != nil{
		w.Write([]byte("Error while trying to read the user data"))
		return
	}
	fmt.Println(user)

}
