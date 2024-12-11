package server

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"strconv"

	"github.com/Joao-lucas-felix/BasicCRUD/database"
	"github.com/gorilla/mux"
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

// Find all users in the data base
func FindAllUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to connect with DataBase"))
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM usuarios u")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to get the users in DataBase!"))
		return
	}
	defer rows.Close()

	var users []user

	for rows.Next() {
		var user user
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error while trying to convert the user of DataBase response!"))
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to convert the user to json"))
		return
	}

}

// Search one user in the data base
func FindUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to read the ID"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error"))
		return
	}
	defer db.Close()

	row, err := db.Query("SELECT * FROM usuarios WHERE id = $1", ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error, while trying to get the user in the data base"))
		return
	}
	var u user
	if row.Next() {
		if err := row.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Error, while trying to scan the user"))
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to convert the user to json"))
		return
	}
}

// Find and updates an user in the data base
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to read the ID"))
		return
	}
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Fail do read the request body!"))
		return
	}

	var userToUpdate user
	if err = json.Unmarshal(requestBody, &userToUpdate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to read the user data"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to Connect with the database"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("UPDATE usuarios  SET nome = $1, email = $2 WHERE id = $3")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to Prepare the statment" + err.Error()))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(userToUpdate.Name, userToUpdate.Email, ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to update the user"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to read the params"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to connect with the database"))
		return
	}
	statement, err := db.Prepare("DELETE FROM usuarios  WHERE id = $1")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to prepare the statement"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to delete the user"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
