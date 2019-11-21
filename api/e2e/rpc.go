package main

import (
	"context"
	"log"
	"time"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
	"google.golang.org/grpc/metadata"
)

var funcMap map[string]func(pg.APIClient) = map[string]func(pg.APIClient){
	"1": RegisterPlayer,
	"2": RequestPlayerLogin,
	"3": PlayerLogin,
	"4": GetSchedule,
	"5": GetScores,
	"6": GetScoresForPlayer,
}

func RegisterPlayer(client pg.APIClient) {
	eventKey := getInput("eventKey", "starts-in-30m")
	phoneNumber := getInput("phoneNumber", "+5551231234")
	name := getInput("name", "Eric Morris")

	league := pg.League_NONE
	leagueStr := getInput("league", "PGA")
	leagueInt, exists := pg.League_value[leagueStr]
	if exists {
		league = pg.League(leagueInt)
	}

	log.Println("Making call to RegisterPlayer")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := client.RegisterPlayer(ctx, &pg.RegisterPlayerRequest{
		EventKey:    eventKey,
		PhoneNumber: phoneNumber,
		Name:        name,
		League:      league,
	})
	logResponse(r, err)
}

func RequestPlayerLogin(client pg.APIClient) {
	eventKey := getInput("eventKey", "starts-in-30m")
	phoneNumber := getInput("phoneNumber", "+5551231234")

	log.Println("Making call to RequestPlayerLogin")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := client.RequestPlayerLogin(ctx, &pg.RequestPlayerLoginRequest{
		EventKey:    eventKey,
		PhoneNumber: phoneNumber,
	})
	logResponse(r, err)
}

func PlayerLogin(client pg.APIClient) {
	eventKey := getInput("eventKey", "starts-in-30m")
	phoneNumber := getInput("phoneNumber", "+5551231234")
	authCode, err := getInputAsUInt32("authCode", 111111)
	if err != nil {
		log.Printf("Could not parse input as a number: %s", err)
		return
	}

	log.Println("Making call to PlayerLogin")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := client.PlayerLogin(ctx, &pg.PlayerLoginRequest{
		EventKey:    eventKey,
		PhoneNumber: phoneNumber,
		AuthCode:    authCode,
	})
	logResponse(r, err)
}

func GetSchedule(client pg.APIClient) {
	authToken := getInput("authToken", "00000000-0000-4000-a000-300000000000")
	eventKey := getInput("eventKey", "current")

	log.Println("Making call to GetSchedule")
	header := metadata.New(map[string]string{"authorization": authToken})
	ctx := metadata.NewOutgoingContext(context.Background(), header)
	r, err := client.GetSchedule(ctx, &pg.GetScheduleRequest{
		EventKey: eventKey,
	})
	logResponse(r, err)
}

func GetScores(client pg.APIClient) {
	authToken := getInput("authToken", "00000000-0000-4000-a000-300000000000")
	eventKey := getInput("eventKey", "current")

	log.Println("Making call to GetScores")
	header := metadata.New(map[string]string{"authorization": authToken})
	ctx := metadata.NewOutgoingContext(context.Background(), header)
	r, err := client.GetScores(ctx, &pg.GetScoresRequest{
		EventKey: eventKey,
	})
	logResponse(r, err)
}

func GetScoresForPlayer(client pg.APIClient) {
	authToken := getInput("authToken", "00000000-0000-4000-a000-300000000000")
	eventKey := getInput("eventKey", "current")
	playerID := getInput("eventKey", "00000000-0000-4000-a000-200000000000")

	log.Println("Making call to GetScoresForPlayer")
	header := metadata.New(map[string]string{"authorization": authToken})
	ctx := metadata.NewOutgoingContext(context.Background(), header)
	r, err := client.GetScoresForPlayer(ctx, &pg.GetScoresForPlayerRequest{
		EventKey: eventKey,
		PlayerID: playerID,
	})
	logResponse(r, err)
}
