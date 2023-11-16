package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func calculateFactorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * calculateFactorial(n-1)
}

func calculateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var inputData InputData

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"Incorrect input"}`)
		return
	}

	if inputData.A < 0 || inputData.B < 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"Incorrect input"}`)
		return
	}

	resultChan := make(chan ResponseData)
	go func() {
		resultChan <- ResponseData{AFactorial: calculateFactorial(inputData.A), BFactorial: calculateFactorial(inputData.B)}
	}()

	result := <-resultChan

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func main() {
	router := httprouter.New()
	router.POST("/calculate", calculateHandler)
	http.ListenAndServe(":8989", router)
}
