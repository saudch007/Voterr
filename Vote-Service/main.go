package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	supa "github.com/nedpals/supabase-go"
	"github.com/streadway/amqp"
)

type CandidateData struct { // storing values in custom type
	CandidateName  string `json:"candidate_name"`
	CandidateVotes int    `json:"candidate_votes"`
}

func updateVotes(passedName string, passedVote int) {

	e := godotenv.Load("../.env")
	if e != nil {
		log.Fatalf("Error loading .env file: %v", e)
	}

	supabaseURL := os.Getenv("DB_URL")
	supabaseKEY := os.Getenv("DB_KEY")
	supabase := supa.CreateClient(supabaseURL, supabaseKEY)

	var results []map[string]interface{}
	query := supabase.DB.From("ballettable").Select("*").Eq("name", passedName)
	if err := query.Execute(&results); err != nil {
		fmt.Println("Error in querying row := ", err)
		return
	} else {
		fmt.Println("Success!")
		fmt.Println(results)
	}
	////////// Till now we have found the specific candidate
	//// Now update its votes based on passedVote
	passedVote += 1
	type newCandidateVal struct {
		Name string `json:"name"`
		Vote int    `json:"votes"`
	}
	updateRow := newCandidateVal{
		Name: passedName,
		Vote: passedVote,
	}

	var newResults []map[string]interface{}

	updateVotes := supabase.DB.From("ballettable").Update(updateRow).Eq("name", passedName)
	if data := updateVotes.Execute(&newResults); data != nil {
		fmt.Println("Error in update query := ", data)
	} else {
		fmt.Println("Success in update query!")
		fmt.Println(newResults)
	}
}

func consumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Update with your RabbitMQ connection details
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"candidate_queue", // Queue name
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	fmt.Println("Vote-Service is waiting for messages...")

	// Start consuming messages
	for msg := range msgs {
		fmt.Println("Received message:", string(msg.Body))

		var data CandidateData // data variable type of CandidateData
		err := json.Unmarshal(msg.Body, &data)
		if err != nil {
			log.Println("Error unmarshaling message:", err)
			continue
		}

		// Process the received candidate data
		fmt.Println("Candidate Name:", data.CandidateName)
		fmt.Println("Candidate Votes:", data.CandidateVotes)

		// Perform further actions based on candidate data

		candidateName := data.CandidateName
		candidateVotes := data.CandidateVotes

		updateVotes(candidateName, candidateVotes) // updating votes

	}
}

func voteServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}

	e := godotenv.Load("../.env")
	if e != nil {
		log.Fatalf("Error loading .env file: %v", e)
		return
	}

	// consuming
	consumer()
}

func main() {
	http.HandleFunc("/vote-service", voteServiceHandler)
	fmt.Println("Sign-Up service is running at port 6000")
	http.ListenAndServe(":6000", nil)
}
