package main

import (
	"context"
	"fmt"
	"grpcCourse/calculator/calcpb"
	"io"
	"log"
	"time"

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
	//getDecomp(c)
	getAverage(c)
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
		log.Fatalf("Error calling decomp: %v", err)
	}
	log.Printf("Calculation: %v", res.Result)
}

func getDecomp(c calcpb.CalculationServiceClient) {
	req := &calcpb.PrimeNumberDecompRequest{
		NumberToDecomp: 120,
	}
	resStream, err := c.PrimeNumberDecomp(context.Background(), req)

	if err != nil {
		log.Fatalf("error calling stream decomp %v", err)
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

func getAverage(c calcpb.CalculationServiceClient) {
	stream, err := c.ComputeAverage(context.Background())

	requests := []*calcpb.ComputeAverageRequest{
		&calcpb.ComputeAverageRequest{
			Number: 1,
		},
		&calcpb.ComputeAverageRequest{
			Number: 2,
		},
		&calcpb.ComputeAverageRequest{
			Number: 3,
		},
		&calcpb.ComputeAverageRequest{
			Number: 4,
		},
	}

	if err != nil {
		log.Fatalf("error calling average: %v", err)
	}

	for _, req := range requests {
		fmt.Println("Sending req")
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error recieving average: %v", err)
	}

	fmt.Printf("Average: %v", res)
}
