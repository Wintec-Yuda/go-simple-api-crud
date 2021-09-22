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
	db, err := gorm.Open("mysql", "root:@/name_database?charset=utf8&parseTime=True")

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

func checkErrResponse(err error)  {
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleRequest()  {
	log.Println("Start the development server at http://127.0.0.1:9999")

	router := mux.NewRouter().StrictSlash(true)

	router.NotFoundHandler = http.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		res := Result{Code: 404, Message: "Method not found"}
		respone, err := json.Marshal(res)

		checkErr(err)
		w.Write(respone)
	})
	router.MethodNotAllowedHandler = http.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		res := Result{Code: 405, Message: "Method not allowed"}
		response, err := json.Marshal(res)

		checkErr(err)
		w.Write(response)
	})

	router.HandleFunc("/api/products", createProduct).Method("POST")
	router.HandleFunc("/api/products", getProducts).Method("GET")
	router.HandleFunc("/api/products/{id}", getProductById).Method("GET")
	router.HandleFunc("/api/products/{id}", updateProduct).Method("PUT")
	router.HandleFunc("/api/products/{id}", deleteProduct).Method("DELETE")

	log.Fatal(http.ListenAndServe(":9999", router))
}

func createProduct(w http.ResponseWriter, r *http.Request)  {
	payloads, err := ioutil.ReadAll(r.Body)

	checkErr(err)

	var product Product
	json.Unmarshal(payloads, &product)

	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Created product successfully"}
	result, err := json.Marshal(res)

	checkErrResponse(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProducts(w http.ResponseWriter, r *http.Request)  {
	products := []Product{}
	db.Find(&products)

	res := Result{Code: 200, Data: product, Message: "Getted product successfully"}
	result, err := json.Marshal(res)

	checkErrResponse(err)

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
	result, err := json.Marshal(res)

	checkErrResponse(err)

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
	result, err := json.Marshal(res)

	checkErrResponse(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}