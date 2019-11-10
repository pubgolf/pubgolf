package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"google.golang.org/grpc"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

const (
	address = "localhost:50051"
)

func main() {
	log.SetPrefix("[Test Client] ")

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Error: could not connect: %v", err)
	}
	defer conn.Close()
	c := pg.NewAPIClient(conn)

	log.Println("Making call to GetSchedule")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetSchedule(ctx, &pg.GetScheduleRequest{EventKey: "nyc-2019"})
	if err != nil {
		log.Fatalf("Error getting schedule: %v", err)
	}

	if respStr, err := json.MarshalIndent(r.GetVenueList(), "", "\t"); err != nil {
		log.Fatalf("Error formatting schedule as JSON: %v", err)
	} else {
		log.Printf("GetSchedule Response:\n%s", respStr)
	}
}
