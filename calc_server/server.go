package main

import (
	"context"
	"fmt"
	"grpcCourse/calculator/calcpb"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculate(ctx context.Context, req *calcpb.CalculationRequest) (*calcpb.CalculationResponse, error) {

	responseValue := req.GetCalculation().GetFirstInt() + req.GetCalculation().GetSecondInt()

	res := &calcpb.CalculationResponse{
		Result: responseValue,
	}
	return res, nil
}

func (*server) PrimeNumberDecomp(req *calcpb.PrimeNumberDecompRequest, stream calcpb.CalculationService_PrimeNumberDecompServer) error {
	intToDecompose := req.GetNumberToDecomp()
	var k int32
	k = 2
	for intToDecompose > 1 {
		fmt.Println(k)
		if intToDecompose%k == 0 {
			res := &calcpb.PrimeNumberDecompResponse{
				Result: k,
			}
			stream.Send(res)
			intToDecompose = intToDecompose / k
			k = 2
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calcpb.CalculationService_ComputeAverageServer) error {
	result := float32(0)
	n := float32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//Finished reading from stream
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Result: result / n,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading: %v", err)
		}

		result += req.GetNumber()
		n++
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterCalculationServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
