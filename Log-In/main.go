package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	supa "github.com/nedpals/supabase-go"
)

func signInHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome to Log In")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	e := godotenv.Load("../.env")
	if e != nil {
		log.Fatalf("Error loading .env file: %v", e)
	}

	Email := requestData.Email
	Password := requestData.Password

	if Email == "" || Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		fmt.Printf("Email is : %s and password is %s", Email, Password)
		supabaseUrl := os.Getenv("DB_URL")
		supabaseKey := os.Getenv("DB_KEY")
		supabase := supa.CreateClient(supabaseUrl, supabaseKey)

		ctx := context.Background()
		user, err := supabase.Auth.SignIn(ctx, supa.UserCredentials{
			Email:    Email,
			Password: Password,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(user)
	}

}

func main() {
	http.HandleFunc("/signin", signInHandler)
	fmt.Println("Log-In service is running at port 9000")
	http.ListenAndServe(":9000", nil)
}
