// gabe/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/dicurrio/protorepo/fred"
	"google.golang.org/grpc"
)

var (
	hostAddress = os.Getenv("HOST_ADDRESS")
	fredAddress = os.Getenv("FRED_ADDRESS")
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v %v", r.Proto, r.Method, r.URL.EscapedPath())
	name := "Friend"

	conn, err := grpc.Dial(fredAddress, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to dial Fred: %v", err)
	}
	defer conn.Close()

	fredClient := pb.NewFredClient(conn)
	req := &pb.Request{
		Name: name,
	}
	res, err := fredClient.GetIndex(context.Background(), req)
	if err != nil {
		log.Printf("Fred GetIndex error: %v", err)
		fmt.Fprintf(w, "Internal Server Error")
	} else {
		fmt.Fprintf(w, res.GetMessage())
	}
}

func main() {
	// Setup
	log.SetPrefix("GABE :: ")
	log.Print("Starting up...")

	// Handlers
	http.HandleFunc("/", indexHandler)

	// Server Startup
	server := http.Server{
		Addr:         hostAddress,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		log.Print("Listening on " + hostAddress)
		log.Fatal(server.ListenAndServe())
	}()

	// Graceful Shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan // Blocks until SIGINT or SIGTERM received
	log.Print("Shutdown signal received, exiting...")
	server.Shutdown(context.Background())
}
