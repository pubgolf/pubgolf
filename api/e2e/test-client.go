package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"

	"google.golang.org/grpc"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

const (
	address = "localhost:50051"
)

func printHelp() {
	fmt.Println("")
	fmt.Println("Enter number to select RPC method to test")
	fmt.Println("(or type 'q' to quit)")
	fmt.Println("-----------------------------------------")
	for label, fn := range funcMap {
		fmt.Printf("%s: %s\n", label,
			strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")[1])
	}
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
	cmd := getCmd(buf)
	for ; cmd != "q"; cmd = getCmd(buf) {
		fmt.Println("")
		if fn, exists := funcMap[cmd]; exists {
			fn(client)
		}
		printHelp()
	}
}
