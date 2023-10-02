package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// calls method the connects to database
	InitDB()

	router := httprouter.New()
	router.GET("/employees", getEmployees)
	router.GET("/employees/:id", getEmployee)
	router.POST("/employees", createEmployee)
	router.PUT("/employees/:id", updateEmployee)
	router.DELETE("/employees/:id", deleteEmployee)

	fmt.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
