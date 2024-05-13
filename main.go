package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// representasi Json
type Products struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	//1. buat route multiplexer
	mux := http.NewServeMux()

	//3. tambahkan handler
	mux.HandleFunc("GET /products", listProduct)
	mux.HandleFunc("POST /products", createProduct)
	mux.HandleFunc("PUT /products/{id}", updateProduct)
	mux.HandleFunc("DELETE /products/{id}", deleteProduct)

	//4. buat server
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	//5. jalankan server
	server.ListenAndServe()
}

var database = map[int]Products{}

var lastId = 0

// 2. fungsi handler
func listProduct(w http.ResponseWriter, r *http.Request) {
	//slice untuk response
	var products []Products

	//itterasi / looping data product dan tambahkan kedalam slice
	for _, v := range database {
		products = append(products, v)
	}

	//kita ubah menjadi json
	data, err := json.Marshal(products)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("Terjadi kesalahan"))
	}

	//lempar json yanng sudah kita ubah tadi ke dalam response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
func createProduct(w http.ResponseWriter, r *http.Request) {

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Kesalahan dalam request"))
	}
	var products Products
	err = json.Unmarshal(bodyByte, &products)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Kesalahan dalam request"))
	}

	//increment no ururt
	lastId++

	products.ID = lastId

	//tambahkan ke db
	database[products.ID] = products
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte("Request berhasil di proses dan ditambahkan"))
}
func updateProduct(w http.ResponseWriter, r *http.Request) {
	//baca id
	productID := r.PathValue("id")
	productInt, err := strconv.Atoi(productID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Kesalahan dalam request"))
	}

	//baca req body
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Kesalahan dalam request"))
	}
	var products Products
	err = json.Unmarshal(bodyByte, &products)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Kesalahan dalam request"))
	}

	//supaya id tidak tergantikan dari requst body
	products.ID = productInt
	database[productInt] = products

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)

}
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")

	productInt, err := strconv.Atoi(productID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Kesalahan dalam request"))
	}

	delete(database, productInt)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
}
