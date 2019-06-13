package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/maxbrain0/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function invoked with %v", req)
	term1 := req.GetSum().GetNum1()
	term2 := req.GetSum().GetNum2()

	res := &calculatorpb.SumResponse{
		Result: term1 + term2,
	}

	return res, nil
}

func main() {
	fmt.Println("In server main...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
