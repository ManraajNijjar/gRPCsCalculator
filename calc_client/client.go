package main

import (
	"context"
	"grpcCourse/calculator/calcpb"
	"io"
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

	//getCalc(c)
	getDecomp(c)
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

func getDecomp(c calcpb.CalculationServiceClient) {
	req := &calcpb.PrimeNumberDecompRequest{
		NumberToDecomp: 120,
	}
	resStream, err := c.PrimeNumberDecomp(context.Background(), req)

	if err != nil {
		log.Fatalf("error calling stream greet: %v", err)
	}

	for {
		factor, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading stream: %v", err)
		}
		log.Printf("factor: %v", factor.GetResult())
	}
}
