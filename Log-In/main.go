package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	supa "github.com/nedpals/supabase-go"
	"github.com/streadway/amqp"
)

var Email string
var Pass string

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}

func emitSignInEvent() {
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

	q, err := ch.QueueDeclare(
		"sign_up_queue", // Queue name
		true,            // Durable
		false,           // Delete when unused
		false,           // Exclusive
		false,           // No-wait
		nil,             // Arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,        // Queue name
		"",            // Routing key
		"user_events", // Exchange
		false,         // No-wait
		nil,
	)
	failOnError(err, "Failed to bind to the queue")

	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer
		true,   // Auto-ack
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		data := string(d.Body)
		parts := strings.Split(data, " ")
		if len(parts) != 2 {
			log.Printf("Invalid message format: %s", data)
			continue
		}

		email := parts[0]
		password := parts[1]

		////////

		Email = email
		Pass = password

		fmt.Printf("Received sign-up event for user: Email: %s, Password: %s\n", email, password)
	}

}

func signInHandler(w http.ResponseWriter, r *http.Request) {

	go emitSignInEvent()
	time.Sleep(time.Second * 3)

	fmt.Println("Welcome to Log In")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}

	e := godotenv.Load("../.env")
	if e != nil {
		log.Fatalf("Error loading .env file: %v", e)
	}

	if Email == "" || Pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		fmt.Printf("Email is : %s and password is %s", Email, Pass)
		supabaseUrl := os.Getenv("DB_URL")
		supabaseKey := os.Getenv("DB_KEY")
		supabase := supa.CreateClient(supabaseUrl, supabaseKey)

		ctx := context.Background()
		user, err := supabase.Auth.SignIn(ctx, supa.UserCredentials{
			Email:    Email,
			Password: Pass,
		})
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(user)

		}
	}
}

func main() {

	http.HandleFunc("/signin", signInHandler)
	fmt.Println("Log-In service is running at port 9000")
	http.ListenAndServe(":9000", nil)
}
