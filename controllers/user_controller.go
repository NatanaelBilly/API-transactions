package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Connect to database
	db := Connect()
	defer db.Close()

	query := "SELECT * FROM users"

	name := r.URL.Query()["name"]
	if name != nil {
		query += "WHERE name'" + name[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Print(err)
	}

	var user User
	var users []User

	// get data
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.email, &user.password)
		if err != nil { // error response
			var response ErrorResponse
			response.Status = 400
			response.Message = "Bad Request"
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			log.Println(err)
			return
		} else {
			users = append(users, user)
		}
	}

	// success response
	var response UsersResponse
	response.Status = 200
	response.Message = "Request success"
	response.Data = users

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	db.Close()
}

func InsertNewUser(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		var response ErrorResponse
		response.Status = 500
		response.Message = "Internal server error"
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")

	_, errQuerry := db.Exec("INSERT INTO users(Name, Age, Address) VALUES (?,?,?);", name, age, address)

	var response UserResponse
	if errQuerry == nil {
		response.Status = 200
		response.Message = "success"
		w.WriteHeader(200)
	} else {
		response.Status = 500
		response.Message = "Internal server error"
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		var response ErrorResponse
		response.Status = 500
		response.Message = "Internal server error"
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	vars := mux.Vars(r)
	userId := vars["id"]
	data, _ := db.Query(`SELECT * FROM users WHERE id = ?;`, userId)

	if data == nil {
		var response ErrorResponse
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", userId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	_, errQuerry := db.Query(`DELETE FROM users WHERE id = ?;`, userId)

	var response UserResponse
	if errQuerry == nil {
		response.Status = 200
		response.Message = "success"
		w.WriteHeader(200)
	} else {
		response.Status = 500
		response.Message = "Internal server error"
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()

	if err != nil {
		var response ErrorResponse
		response.Status = 500
		response.Message = "Internal server error"
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	vars := mux.Vars(r)
	userId := vars["id"]

	data, _ := db.Query(`SELECT * FROM users WHERE id = ?;`, userId)

	if data == nil {
		var response ErrorResponse
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", userId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")

	_, errQuerry := db.Exec("UPDATE users set name = ?, age = ?, address = ? WHERE id = ?;", name, age, address, userId)

	var response UserResponse
	if errQuerry == nil {
		response.Status = 200
		response.Message = "success"
		w.WriteHeader(200)
	} else {
		response.Status = 500
		response.Message = "Internal server error"
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()

	if err != nil {
		var response ErrorResponse
		response.Status = 500
		response.Message = "Internal server error"
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	var user User
	user.email = r.Form.Get("email")
	user.password = r.Form.Get("password")
	fmt.Println(user.email)
	fmt.Println(user.password)
	_, errQuerry := db.Query("select * from users  WHERE email = ? AND password = ?;", user.email, user.password)

	var response UserResponse
	if errQuerry == nil {
		response.Status = 200
		response.Message = "login success"
		w.WriteHeader(200)
	} else {
		response.Status = 500
		response.Message = "Internal server error"
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
