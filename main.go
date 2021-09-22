package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

type Product struct {
	ID int `json:"id" form:"id"`
	Code string `json:"code" form:"code"`
	Name string `json:"name" form:"name"`
	Price decimal.Decimal `json:"price" form:"price" sql:"type:decimal(16, 2)"`
}

type Result struct {
	Code int `json:"code"`
	Data interface{} `json:"data"`
	Message string `json:"message"`
}

func main() {
	db, err := gorm.Open("mysql", "root:root@/simple_api?charset=utf8&parseTime=True")

	checkErr(err)

	log.Println("Connection database successfully")

	db.AutoMigrate(&Product{})
	handleRequest()
}

func checkErr(err error)  {
	if err != nil{
		log.Println(err)
	}
	return
}

func handleRequest()  {
	log.Println("Start the development server at http://127.0.0.1:9999")

	router := mux.NewRouter().StrictSlash(true)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		res := Result{Code: 404, Message: "Method not found"}
		respone, err := json.Marshal(res)

		checkErr(err)
		w.Write(respone)
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		res := Result{Code: 405, Message: "Method not allowed"}
		response, _ := json.Marshal(res)

		w.Write(response)
	})

	router.HandleFunc("/api/products", createProduct).Methods("POST")
	router.HandleFunc("/api/products", getProducts).Methods("GET")
	router.HandleFunc("/api/products/{id}", getProductById).Methods("GET")
	router.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	router.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9999", router))
}

func createProduct(w http.ResponseWriter, r *http.Request)  {
	payloads, err := ioutil.ReadAll(r.Body)

	checkErr(err)

	var product Product
	json.Unmarshal(payloads, &product)

	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Created product successfully"}
	result, _ := json.Marshal(res)


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProducts(w http.ResponseWriter, r *http.Request)  {
	products := []Product{}
	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Getted product successfully"}
	result, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProductById(w http.ResponseWriter, r *http.Request)  {
	param := mux.Vars(r)
	id := param["id"]

	payloads, err := ioutil.ReadAll(r.Body)

	checkErr(err)

	var productUpdate Product
	json.Unmarshal(payloads, &productUpdate)

	var product Product
	db.First(&product, id)
	db.Model(&product).Updates(productUpdate)

	res := Result{Code: 200, Data: product, Message: "Updated product successfully"}
	result, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request)  {
	param := mux.Vars(r)
	id := param["id"]

	var product Product

	db.First(&product, id)
	db.Delete(&product)

	res := Result{Code: 200, Data: product, Message: "Deleted product successfully"}
	result, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request)  {
	param := mux.Vars(r)
	id := param["id"]

	payload, _ := ioutil.ReadAll(r.Body)

	var productUpdate Product
	json.Unmarshal(payload, &productUpdate)

	var product Product
	db.First(&product, id)
	db.Model(&product).Updates(productUpdate)

	res := Result{Code: 200, Data: product, Message: "Updated product successfully"}
	result, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}