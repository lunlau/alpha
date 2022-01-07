package main

import (
"context"
"log"

"google.golang.org/grpc"
	pb "github.com/lunlau/alpha/pb"
)

const PORT = "9001"

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewAlphaRuleEngineClient(conn)
	resp, err := client.AddRule(context.Background(), &pb.AddRuleRequest{
		Rules: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetRules())
}