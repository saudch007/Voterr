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
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}

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

			emitSignUpEvent(Email, Password)

		}

	}

}

func emitSignUpEvent(email string, password string) {
	// RabbitMQ Connection //

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create a channel: %s", err)
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(
		"user_events",
		"fanout", // Exchange type
		true,     // Durable
		false,    // Auto-deleted
		false,    // Internal
		false,    // No-wait
		nil,
	)

	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"user_events", // Exchange
		"",            // Routing key (not used for fanout)
		false,         // Mandatory
		false,         // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(email + " " + password),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}

	log.Printf("Sign-up event emitted for user: %s", email)

}

func main() {
	http.HandleFunc("/signup", signUpHandler)
	fmt.Println("Sign-Up service is running at port 8000")
	http.ListenAndServe(":8000", nil)
}
