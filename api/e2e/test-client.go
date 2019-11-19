package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

const (
	address = "localhost:50051"
)

func logResponse(resp interface{}, err error) {
	if err != nil {
		log.Printf("Error making call: %v", err)
	}

	if respStr, err := json.MarshalIndent(resp, "", "\t"); err != nil {
		log.Printf("Error formatting response as JSON: %v", err)
	} else {
		log.Printf("Response:\n%s", respStr)
	}
}

func GetSchedule(client pg.APIClient) {
	log.Println("Making call to GetSchedule")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := client.GetSchedule(ctx, &pg.GetScheduleRequest{
		EventKey: "nyc-2019",
	})
	logResponse(r, err)
}

func RegisterPlayer(client pg.APIClient) {
	log.Println("Making call to RegisterPlayer")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := client.RegisterPlayer(ctx, &pg.RegisterPlayerRequest{
		EventKey:    "nyc-2019",
		PhoneNumber: "+5551231234",
		Name:        "Eric Morris",
		// League:      pg.League_LPGA,
	})
	logResponse(r, err)
}

func RequestPlayerLogin(client pg.APIClient) {
	log.Println("Making call to RequestPlayerLogin")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := client.RequestPlayerLogin(ctx, &pg.RequestPlayerLoginRequest{
		EventKey:    "nyc-2019",
		PhoneNumber: "+5551231234",
	})
	logResponse(r, err)
}

func PlayerLogin(client pg.APIClient) {
	buf := bufio.NewReader(os.Stdin)
	fmt.Println("Enter login code:")
	fmt.Print("> ")
	code, err := buf.ReadBytes('\n')
	if err != nil {
		log.Printf("Could not read input: %s", err)
	}

	codeAsNum, err := strconv.Atoi(string(code[:len(code)-1]))
	if err != nil {
		log.Printf("Could not parse input as a number: %s", err)
	}

	log.Println("Making call to PlayerLogin")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := client.PlayerLogin(ctx, &pg.PlayerLoginRequest{
		EventKey:    "nyc-2019",
		PhoneNumber: "+5551231234",
		AuthCode:    uint32(codeAsNum),
	})
	logResponse(r, err)
}

func get(r *bufio.Reader) string {
	t, _ := r.ReadString('\n')
	return strings.TrimSpace(t)
}

func printHelp() {
	fmt.Println("")
	fmt.Println("Enter number to select RPC method to test")
	fmt.Println("(or type 'q' to quit)")
	fmt.Println("-----------------------------------------")
	fmt.Println("1: GetSchedule")
	fmt.Println("2: RegisterPlayer")
	fmt.Println("3: RequestPlayerLogin")
	fmt.Println("4: PlayerLogin")
	fmt.Print("> ")
}

func main() {
	log.SetPrefix("[Test Client] ")

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Error: could not connect: %v", err)
	}
	defer conn.Close()
	client := pg.NewAPIClient(conn)

	buf := bufio.NewReader(os.Stdin)
	printHelp()
	cmd := get(buf)
	for ; cmd != "q"; cmd = get(buf) {
		fmt.Println("")
		switch cmd {
		case "1":
			GetSchedule(client)
		case "2":
			RegisterPlayer(client)
		case "3":
			RequestPlayerLogin(client)
		case "4":
			PlayerLogin(client)
		}
		printHelp()
	}
}
