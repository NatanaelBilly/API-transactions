package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	controllers "modul2-APIGet/controllers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/users", controllers.InsertNewUser).Methods("POST")
	router.HandleFunc("/users/{id}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE")

	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	router.HandleFunc("/products", controllers.InsertNewProduct).Methods("POST")
	router.HandleFunc("/products/{id}", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", controllers.DeleteProduct).Methods("DELETE")

	router.HandleFunc("/transactions", controllers.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions", controllers.InsertNewTransaction).Methods("POST")
	router.HandleFunc("/transactions/{id}", controllers.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{id}", controllers.DeleteProduct).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8000")
	log.Println("Connected to port 8000")
	http.ListenAndServe(":8000", router)
}
