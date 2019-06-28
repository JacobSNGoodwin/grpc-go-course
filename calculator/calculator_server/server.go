package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/maxbrain0/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function invoked with %v\n", req)
	term1 := req.GetSum().GetNum1()
	term2 := req.GetSum().GetNum2()

	res := &calculatorpb.SumResponse{
		Result: term1 + term2,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function invoked with %v\n", req)

	number := req.GetNumber()
	var k int64
	k = 2

	for number > 1 {
		if number%k == 0 {
			// send a response
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: k,
			}

			stream.Send(res)

			number = number / k
		} else {
			k = k + 1
		}
	}

	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with a streaming request")

	var sum int64
	var count int64

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			average := float64(sum) / float64(count)
			// finished reading client stream
			// SendAndClose returns an error, so we can return who thing-a-ma-bob
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: average,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		count++
		sum += req.GetNumber()
	}
}

func main() {
	fmt.Println("In server main...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
