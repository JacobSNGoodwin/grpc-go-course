package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/maxbrain0/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Blog Service Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting Server....")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit - Go specific
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop() // stopping grpc server
	fmt.Println("Closing listener")
	lis.Close()
	fmt.Println("End of program")
}
