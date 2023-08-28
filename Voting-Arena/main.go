package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	supa "github.com/nedpals/supabase-go"
)

// putting getTable in a function to update it incase of changes
func getTable(w http.ResponseWriter, r *http.Request) {
	// Fetching and showing Ballet Box Table
	supabaseUrl := os.Getenv("DB_URL")
	supabaseKey := os.Getenv("DB_KEY")
	supabase := supa.CreateClient(supabaseUrl, supabaseKey)

	var results []map[string]interface{}
	err := supabase.DB.From("ballettable").Select("*").Execute(&results)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		log.Println(err)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(&results)
		if err != nil {
			fmt.Println(err)

		}

	}

	fmt.Println(results) // Selected rows
}

func votingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}

	e := godotenv.Load("../.env")
	if e != nil {
		log.Fatalf("Error loading .env file: %v", e)
	}

	getTable(w, r)

}

func main() {

	http.HandleFunc("/votingArena", votingHandler)
	fmt.Println("Sign-Up service is running at port 5000")
	http.ListenAndServe(":5000", nil)

}
