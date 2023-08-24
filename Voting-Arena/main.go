package main

import (
	"fmt"
	"net/http"
)

func votingHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {

	http.HandleFunc("/votingArena", votingHandler)
	fmt.Println("Sign-Up service is running at port 5000")
	http.ListenAndServe(":5000", nil)

}
