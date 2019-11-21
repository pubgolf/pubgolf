package main

import (
	"context"
	"log"
	"time"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
	"google.golang.org/grpc/metadata"
)

var funcMap map[string]func(pg.APIClient) = map[string]func(pg.APIClient){
	"1": GetSchedule,
	"2": RegisterPlayer,
	"3": RequestPlayerLogin,
	"4": PlayerLogin,
}

func GetSchedule(client pg.APIClient) {
	authToken := getInput("authToken", "4de3e04c-eff3-4631-a493-c5156cb94153")
	eventKey := getInput("eventKey", "nyc-2019")

	log.Println("Making call to GetSchedule")
	header := metadata.New(map[string]string{"authorization": authToken})
	ctx := metadata.NewOutgoingContext(context.Background(), header)
	r, err := client.GetSchedule(ctx, &pg.GetScheduleRequest{
		EventKey: eventKey,
	})
	logResponse(r, err)
}

func RegisterPlayer(client pg.APIClient) {
	eventKey := getInput("eventKey", "nyc-2019")
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
	eventKey := getInput("eventKey", "nyc-2019")
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
	eventKey := getInput("eventKey", "nyc-2019")
	phoneNumber := getInput("phoneNumber", "+5551231234")
	authCode, err := getInputAsUInt32("authCode", 000000)
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
