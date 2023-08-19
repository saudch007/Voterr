package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/joho/godotenv"
)

func votingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	e := godotenv.Load("../.env")
	if e != nil {
		log.Fatalf("Error loading .env file: %v", e)
	}

	var requestData struct {
		CandidateName string `json:"name"`
		Votes         int    `json:"votes"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := requestData.CandidateName
	votes := strconv.Itoa(requestData.Votes)

	var candidateData = []string{name, votes}
	fmt.Println(candidateData)

}

func main() {

	http.HandleFunc("/votingArena", votingHandler)
	fmt.Println("Sign-Up service is running at port 5000")
	http.ListenAndServe(":5000", nil)

}
