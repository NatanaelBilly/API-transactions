package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {

	// Connect to database
	db := Connect()
	defer db.Close()

	// create query
	query := "SELECT * FROM products"

	// get data
	rows, err := db.Query(query)
	if err != nil {
		// error response
		var response ErrorResponse
		response.Status = 400
		response.Message = "Bad Request"
		log.Println(err)
		return
	}

	var product Product
	var products []Product
	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			// error response
			var response ErrorResponse
			response.Status = 400
			response.Message = "Bad Request"
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			log.Fatal(err.Error())
			return
		} else {
			products = append(products, product)
		}
	}

	// success response
	var response ProductsResponse
	response.Status = 200
	response.Message = "Request success"
	response.Data = products

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertNewProduct(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		var response ErrorResponse
		response.Status = 500
		response.Message = "Internal server error"
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		log.Println(err)
		return
	}
	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	_, errQuerry := db.Exec("INSERT INTO products(name, price) VALUES (?,?);", name, price)

	var response ProductResponse
	if errQuerry == nil {
		response.Status = 200
		response.Message = "success"
	} else {
		response.Status = 500
		response.Message = "Internal server error"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
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
	prodId := vars["id"]
	data, _ := db.Query(`SELECT * FROM products WHERE id = ?;`, prodId)

	if data == nil {
		var response ErrorResponse
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", prodId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	_, errQuerry := db.Query(`DELETE FROM products WHERE id = ?;`, prodId)

	var response ProductResponse
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

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
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
	prodId := vars["id"]

	data, _ := db.Query(`SELECT * FROM products WHERE id = ?;`, prodId)

	if data == nil {
		var response ErrorResponse
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", prodId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	_, errQuerry := db.Exec("UPDATE users set name = ?, price = ? WHERE id = ?;", name, price, prodId)

	var response ProductResponse
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
