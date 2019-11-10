package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

const (
	PORT = ":50051"
)

type server struct {
	pg.UnimplementedAPIServer
}

func makeImg() string {
	return fmt.Sprintf("https://via.placeholder.com/640x480/%02x%02x%02x"+
		"?text=Bottle%%20Open%%20NYC%%202019", rand.Intn(255), rand.Intn(255),
		rand.Intn(255))
}

func (s *server) GetSchedule(ctx context.Context,
	req *pg.GetScheduleRequest) (*pg.GetScheduleReply, error) {
	log.Printf("Returning schedule for Event Key: %s", req.GetEventKey())
	if req.GetEventKey() != "nyc-2019" {
		return nil, status.Error(codes.NotFound, "eventKey was not found")
	}

	var venueList pg.VenueList
	for i := uint32(1); i < 10; i++ {
		venue := pg.VenueStop{
			StopID: i,
			Venue: &pg.Venue{
				VenueID:   i,
				Name:      fmt.Sprintf("Foo Bar (Heh) %d", i),
				Address:   fmt.Sprintf("%d Address St, City ST, USA", rand.Intn(9899)+100),
				Image:     makeImg(),
				StartTime: string(time.Now().Format(time.RFC3339)),
			},
		}
		venueList.Venues = append(venueList.Venues, &venue)
	}

	return &pg.GetScheduleReply{VenueList: &venueList}, nil
}

func main() {
	log.SetPrefix("[API Server] ")
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pg.RegisterAPIServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
