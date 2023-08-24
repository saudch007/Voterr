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

	var ValidateAdmin struct {
		AdEmail string `json:"email"`
		AdPass  string `json:"pass"`
		// voting creds
		CandidateNameReq string `json:"cnamereq"`
		VotesReq         int    `json:"votesreq"`
	}

	err := json.NewDecoder(r.Body).Decode(&ValidateAdmin)
	if err != nil {
		log.Printf("Error is %s", err)
	}

	// checking email and pass

	if adminEmail == ValidateAdmin.AdEmail && adminPass == ValidateAdmin.AdPass {

		var BalletBox struct {
			CandidateName string `json:"cname"`
			Votes         int    `json:"votes"`
		}

		// Passed Candidate name and vote from user and extracted from request
		BalletBox.CandidateName = ValidateAdmin.CandidateNameReq
		BalletBox.Votes = ValidateAdmin.VotesReq

		Name := ValidateAdmin.CandidateNameReq
		NVotes := ValidateAdmin.VotesReq

		fmt.Printf("Candidate name: %s and Candidate votes %d", Name, NVotes)

		// creating a table and pushing the name of candidate and votes in supabase

		type Candidate struct {
			Name  string `json:"name"`
			Votes int    `json:"votes"`
		}
		candidates := []Candidate{
			{Name, NVotes},
		}

		// all set now candidates object is to be pushed in table in supabase
		supabaseURL := os.Getenv("DB_URL")
		supabaseKEY := os.Getenv("DB_KEY")
		supabase := supa.CreateClient(supabaseURL, supabaseKEY)

		var results []Candidate

		err := supabase.DB.From("ballettable").Insert(candidates).Execute(&results)
		if err != nil {
			log.Fatalf("Error is %s", err)
		} else {
			fmt.Println(results)
		}

	} else {

		var FailureMsg string = "Email or password is not correct"
		json.NewEncoder(w).Encode(FailureMsg)
	}
}

func main() {

	http.HandleFunc("/admin", adminHandler)
	fmt.Println("Sign-Up service is running at port 3000")
	http.ListenAndServe(":3000", nil)

}
