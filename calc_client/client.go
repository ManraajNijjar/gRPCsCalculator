package main

import (
	"context"
	"grpcCourse/calculator/calcpb"
	"log"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := calcpb.NewCalculationServiceClient(conn)

	getCalc(c)
}

func getCalc(c calcpb.CalculationServiceClient) {
	req := &calcpb.CalculationRequest{
		Calculation: &calcpb.Calculation{
			FirstInt:  10,
			SecondInt: 15,
		},
	}
	res, err := c.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling greet: %v", err)
	}
	log.Printf("Calculation: %v", res.Result)
}
