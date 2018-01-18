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

const (
	hostPort = ":3000"
	fredPort = ":50051"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v %v", r.Proto, r.Method, r.URL.EscapedPath())
	name := "Friend"

	conn, err := grpc.Dial(fredPort, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failled to dial Fred: %v", err)
	}
	defer conn.Close()

	fredClient := pb.NewFredClient(conn)
	req := &pb.Request{
		Name: name,
	}
	res, err := fredClient.GetIndex(context.Background(), req)
	if err != nil {
		log.Printf("Fred GetIndex error: %v", err)
	}
	fmt.Fprintf(w, res.GetMessage())
}

func main() {
	// Setup
	log.SetPrefix("GABE :: ")
	log.Print("Starting up...")

	// Handlers
	http.HandleFunc("/", indexHandler)

	// Server Startup
	server := http.Server{
		Addr:         hostPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		log.Print("Listening on " + hostPort)
		log.Fatal(server.ListenAndServe())
	}()

	// Graceful Shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan // Blocks until SIGINT or SIGTERM received
	log.Print("Shutdown signal received, exiting...")
	server.Shutdown(context.Background())
}
