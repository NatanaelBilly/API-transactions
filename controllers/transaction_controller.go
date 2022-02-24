package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {

	// Connect to database
	db := Connect()
	defer db.Close()

	// create query
	query := "SELECT * FROM transactions"

	// get data
	rows, err := db.Query(query)
	if err != nil {
		// error response
		var response ErrorResponse
		response.Status = 500
		response.Message = "Internal server error"
		log.Println(err)
	}

	var transaction Transaction
	var transactions []Transaction
	for rows.Next() {
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity)
		if err != nil {
			// error response
			var response ErrorResponse
			response.Status = 500
			response.Message = "Internal server error"
			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			log.Fatal(err.Error())
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}

	// success response
	var response TransactionsResponse
	response.Status = 200
	response.Message = "Request success"
	response.Data = transactions

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertNewTransaction(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		var response ErrorResponse
		response.Status = 500
		response.Message = "Internal server error"
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		log.Println(err)
		return
	}
	userId := r.Form.Get("userid")
	productId := r.Form.Get("productid")
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	_, errQuerry := db.Exec("INSERT INTO transactions(userid, productid, quantity) VALUES (?,?,?);", userId, productId, quantity)

	var response UserResponse
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

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
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
	transId := vars["id"]
	data, _ := db.Query(`SELECT * FROM transactions WHERE id = ?;`, transId)

	if data == nil {
		var response ErrorResponse
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", transId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	_, errQuerry := db.Query(`DELETE FROM transactions WHERE id = ?;`, transId)

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

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
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
	transId := vars["id"]

	data, _ := db.Query(`SELECT * FROM transactions WHERE id = ?;`, transId)

	if data == nil {
		var response ErrorResponse
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", transId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	userId := r.Form.Get("userid")
	productId := r.Form.Get("productid")
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	_, errQuerry := db.Exec("UPDATE users set userid = ?, productid = ?, quantity = ? WHERE id = ?;", userId, productId, quantity, transId)

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

func GetUserDetailTransaction(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	query := "select t.id, u.id, u.name, u.age, u.address, p.id, p.name, p.price, t.quantity from transactions t join users u on t.userId=u.id join products p on t.productId=p.id"
	id := r.URL.Query()["id"]
	if id != nil {
		query += "where u.id = " + id[0]
	}
	rows, err := db.Query(query)

	if err != nil {
		var response ErrorResponse
		response.Status = 500
		response.Message = "Internal server error"
		w.WriteHeader(500)
		log.Println(err)
		return
	}

	var listTransaksiDetails []TransactionDetail
	var transactionDetail TransactionDetail
	var user User
	var product Product

	for rows.Next() {
		if err := rows.Scan(&transactionDetail.ID, &user.ID, &user.Name, &user.Age, &user.Address, &product.ID, &product.Name, &product.Price, &transactionDetail.Quantity); err != nil {
			log.Println(err)
			return
		} else {
			transactionDetail.User = user
			transactionDetail.Product = product
			listTransaksiDetails = append(listTransaksiDetails, transactionDetail)
		}
	}

	var response TransactionDetailsResponse
	if len(listTransaksiDetails) > 0 {
		response.Status = 200
		response.Message = "Success Get Data"
		response.Data = listTransaksiDetails
	} else {
		response.Status = 400
		response.Message = "Error Get Data"
		w.WriteHeader(400)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
