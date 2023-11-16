package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type InputData struct {
	A int `json:"a"`
	B int `json:"b"`
}

type ResponseData struct {
	AFactorial int `json:"a"`
	BFactorial int `json:"b"`
}

type Calculator struct {
}

func (c *Calculator) calculateFactorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * c.calculateFactorial(n-1)
}

func (c *Calculator) calculateFactorialsAsync(a, b int) (int, int) {
	var wg sync.WaitGroup
	resultChan := make(chan int, 2)

	wg.Add(2)
	go c.calculateFactorialAsync(a, &wg, resultChan)
	go c.calculateFactorialAsync(b, &wg, resultChan)

	wg.Wait()
	close(resultChan)

	return <-resultChan, <-resultChan
}

func (c *Calculator) calculateFactorialAsync(n int, wg *sync.WaitGroup, resultChan chan int) {
	defer wg.Done()
	resultChan <- c.calculateFactorial(n)
}

func validateInput(input InputData) bool {
	return input.A >= 0 && input.B >= 0
}

func calculateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var inputData InputData

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputData)
	if err != nil || !validateInput(inputData) {
		handleError(w, http.StatusBadRequest, "Incorrect input")
		return
	}

	calculator := Calculator{}
	resultA, resultB := calculator.calculateFactorialsAsync(inputData.A, inputData.B)
	result := ResponseData{
		AFactorial: resultA,
		BFactorial: resultB,
	}

	respondWithJSON(w, http.StatusOK, result)
}

func handleError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprint(w, `{"error":"`+message+`"}`)
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	router := httprouter.New()
	router.POST("/calculate", calculateHandler)
	http.ListenAndServe(":8989", router)
}
