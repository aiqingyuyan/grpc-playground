package client

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
	"yanyu/go-grpc/generated/greet"
)

func GreetExample() {
	var options []grpc.DialOption

	options = append(options, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:50051", options...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	greetClient := greet.NewGreetServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	greetReq := &greet.GreetRequest{
		Greeting: &greet.Greeting{
			FirstName: "test",
			LastName: "shit",
		},
	}

	resp, err := greetClient.Greet(ctx, greetReq)
	if err != nil {
		log.Fatalf("%v.Greet(_) = _, %v", greetClient, err)
	}

	log.Println(resp)
}
