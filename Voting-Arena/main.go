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

// message sender
func sendToVoteService(candidateName string, candidateVote int) error {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Update with your RabbitMQ connection details
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
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
		return err
	}

	message := map[string]interface{}{
		"candidate_name":  candidateName,
		"candidate_votes": candidateVote,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // Exchange
		q.Name, // Routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageJSON,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println("Candidate data sent to RabbitMQ")
	return nil
}

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
	var candidate_name string
	var candidate_votes int

	fmt.Printf("Candidate name?:\n")
	_, err := fmt.Scan(&candidate_name)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	fmt.Printf("Candidate votes?:\n")
	_, errVal := fmt.Scan(&candidate_votes)
	if errVal != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	fmt.Println("Candidate name is :", candidate_name)
	fmt.Println("Candidate vote is :", candidate_votes)

	type Candidate struct {
		Name  string `json:"name"`
		Votes int    `json:"votes"`
	}
	candidates := []Candidate{
		{candidate_name, candidate_votes},
	}

	supabaseURL := os.Getenv("DB_URL")
	supabaseKEY := os.Getenv("DB_KEY")
	supabase := supa.CreateClient(supabaseURL, supabaseKEY)

	var results []Candidate

	err = supabase.DB.From("ballettable").Insert(candidates).Execute(&results)
	if err != nil {
		log.Fatalf("Error is %s", err)
	} else {
		fmt.Println(results)
	}

	err = sendToVoteService(candidate_name, candidate_votes) // Sending to Vote-Service
	if err != nil {
		fmt.Println("Error sending data to Vote-Service", err)
		return
	}

}

func main() {

	http.HandleFunc("/votingArena", votingHandler)
	fmt.Println("Sign-Up service is running at port 5000")
	http.ListenAndServe(":5000", nil)

}
