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

func getEnviromentVar(key string) string {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading env file")
	}

	return os.Getenv(key)
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {

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

	Email := requestData.Email
	Password := requestData.Password

	if Email == "" || Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {

		supabaseURL := getEnviromentVar("DB_URL")
		supabaseKEY := getEnviromentVar("DB_KEY")
		supabase := supa.CreateClient(supabaseURL, supabaseKEY)

		ctx := context.Background()
		user, err := supabase.Auth.SignUp(ctx, supa.UserCredentials{
			Email:    Email,
			Password: Password,
		})
		if err != nil {
			panic(err)
		}

		fmt.Println(user)
	}

}

func main() {
	http.HandleFunc("/signup", signUpHandler)
	fmt.Println("Sign-Up service is running at port 8001")
	http.ListenAndServe(":8001", nil)
}
