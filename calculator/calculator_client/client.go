package main

import (
	"context"
	"fmt"
	"log"

	"github.com/maxbrain0/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("In the client!")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewSumServiceClient(cc)

	doUnary(c)
}

func doUnary(c calculatorpb.SumServiceClient) {
	fmt.Println("Unary client call...")

	req := &calculatorpb.SumRequest{
		Sum: &calculatorpb.Sum{
			Num1: 513,
			Num2: -55,
		},
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Printf("Response from Sum: %v", res.Result)
}
