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

func signUpHandler(w http.ResponseWriter, r *http.Request) {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	e := json.NewDecoder(r.Body).Decode(&requestData)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Email := requestData.Email
	Password := requestData.Password

	if Email == "" || Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {

		supabaseURL := os.Getenv("DB_URL")
		supabaseKEY := os.Getenv("DB_KEY")
		supabase := supa.CreateClient(supabaseURL, supabaseKEY)

		ctx := context.Background()
		user, err := supabase.Auth.SignUp(ctx, supa.UserCredentials{
			Email:    Email,
			Password: Password,
		})
		if err != nil {
			panic(err)
		} else {

			fmt.Println(user)
			redirectURL := fmt.Sprintf("http://localhost:9000/signin?email=%s&password=%s", Email, Password)

			http.Redirect(w, r, redirectURL, http.StatusSeeOther)

		}

	}

}

func main() {
	http.HandleFunc("/signup", signUpHandler)
	fmt.Println("Sign-Up service is running at port 8000")
	http.ListenAndServe(":8000", nil)
}
