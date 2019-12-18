package main

import (
	"context"
	"flag"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/opendevstack/mockbucket/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {

	var portPtr = flag.Int("port", 8081, "set the port in which the service will listen. Default 8081")
	flag.Parse()

	var wait time.Duration

	router := mux.NewRouter()
	subRouter := router.PathPrefix("/rest/api/1.0/projects").Subrouter()
	subRouter.HandleFunc("/", api.CreateProject).Methods("POST")
	subRouter.HandleFunc("/{projectKey}/repos", api.CreateRepository).Methods("POST")
	subRouter.HandleFunc("/{projectKey}/repos/{repositorySlug}/webhooks", api.CreateWebhook).Methods("POST")

	http.Handle("/", router)

	log.Println("Prepering Server...")
	srv := &http.Server{
		Addr: "0.0.0.0:" + strconv.Itoa(*portPtr),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handlers.CORS()(router), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Println("Server running ...")
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

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	<-ctx.Done()
	// to finalize based on context cancellation.
	log.Println("Shutting service down")
	os.Exit(0)

}
