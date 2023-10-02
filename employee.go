package main

type Employee struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Salary  int    `json:"salary"`
	Address string `json:"address"`
}
