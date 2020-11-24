package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kafka/producer/controllers"
)

func main() {
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	producer, _ := controllers.NewProducer([]string{"localhost:9091", "localhost:9092", "localhost:9093"})

	// TOPIC
	r.HandleFunc("/topic/list", controllers.ListTopic).Methods(http.MethodGet)
	r.HandleFunc("/topic/create", controllers.CreateTopic).Methods(http.MethodPost)

	// PRODUCER
	r.HandleFunc("/produce/{topic_name}", producer.StartProduce).Methods(http.MethodPost)

	log.Println("Server started listening on PORT 3000")

	srv := &http.Server{
		Addr: "0.0.0.0:3000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

}
