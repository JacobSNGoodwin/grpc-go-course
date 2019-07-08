package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

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

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a streaming request")

	max := int64(math.MinInt64)

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			// finished reading client stream
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		current := req.GetNumber()

		if current > max {
			max = current

			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Result: max,
			})

			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", err)
				return err
			}
		}
	}

}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot function invoked with %v\n", req)

	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received a negative input: %v", number))
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("In server main...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
