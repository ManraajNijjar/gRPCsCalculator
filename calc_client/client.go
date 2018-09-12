package main

import (
	"context"
	"fmt"
	"grpcCourse/calculator/calcpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
	//getAverage(c)
	//getMax(c)
	getSquareRoot(c)
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

func getMax(c calcpb.CalculationServiceClient) {
	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("error calling findMax: %v", err)
	}
	requests := []*calcpb.FindMaxRequest{
		&calcpb.FindMaxRequest{
			Number: 1,
		},
		&calcpb.FindMaxRequest{
			Number: 5,
		},
		&calcpb.FindMaxRequest{
			Number: 3,
		},
		&calcpb.FindMaxRequest{
			Number: 6,
		},
		&calcpb.FindMaxRequest{
			Number: 2,
		},
		&calcpb.FindMaxRequest{
			Number: 20,
		},
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			fmt.Println("Sending value")
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("error closing stream: %v", err)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
			}
			if err != nil {
				log.Fatalf("Error recieving: %v", err)
				close(waitc)
			}
			fmt.Printf("Received New Max: %v", res.GetResult())
		}
	}()

	<-waitc
}

func getSquareRoot(c calcpb.CalculationServiceClient) {
	fmt.Println("Starting SquareRoot")

	res, err := c.SquareRoot(context.Background(), &calcpb.SquareRootRequest{
		Number: -10,
	})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("Negative Number")
			}
		} else {
			log.Fatalf("Big Error from SquareRoot: %v", err)
		}
	}
	fmt.Println("Square Root: %v", res.GetResult())
}
