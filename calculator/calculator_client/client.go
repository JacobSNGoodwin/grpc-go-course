package main

import (
	"context"
	"fmt"
	"io"
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

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)
	doServerStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
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

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("In doServerStreaming...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 725,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			// end of stream
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes: %v", msg.GetPrimeFactor())
	}
}
