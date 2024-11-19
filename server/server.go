package server

import (
	"crud-basico/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Creates an user and persist in the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Fail do read the request body!"))
		return
	}

	var user user
	if err = json.Unmarshal(requestBody, &user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to read the user data"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to connect with DataBase"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO usuarios (nome, email) VALUES ($1, $2)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error wihile trying to prepare"))
		return
	}
	defer statement.Close()

	insert, err := statement.Exec(user.Name, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to save user"))
		return
	}
	rowsAffected, err := insert.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error " + err.Error()))
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User are created successfully Number of rows afected: %d", rowsAffected)))


}
