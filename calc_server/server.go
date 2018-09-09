package main

import (
	"context"
	"grpcCourse/calculator/calcpb"
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
