package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	e := godotenv.Load("../.env")
	if e != nil {
		log.Fatalf("Error loading .env file: %v", e)
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPass := os.Getenv("ADMIN_PASS")

	var requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recievedEmail := requestData.Email
	recievedPassword := requestData.Password

	if recievedEmail == adminEmail && recievedPassword == adminPass {
		const successMsg string = "Welcome Admin"
		err := json.NewEncoder(w).Encode(successMsg)
		if err != nil {
			fmt.Println("Error in writing json")
		}
	} else {
		const failureMsg string = "Either email or password is not correct"
		err := json.NewEncoder(w).Encode(failureMsg)
		if err != nil {
			fmt.Println("Error in writing json")
		}
		fmt.Println("Redirecting to sign up")
		http.Redirect(w, r, "http://localhost:8000/signup", http.StatusSeeOther)
	}

}

func main() {

	http.HandleFunc("/admin", adminHandler)
	fmt.Println("Sign-Up service is running at port 3000")
	http.ListenAndServe(":3000", nil)

}
