package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strconv"
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
	labelInts := make([]int, 0, len(funcMap))
	for label := range funcMap {
		labelInt, _ := strconv.Atoi(label)
		labelInts = append(labelInts, labelInt)
	}
	sort.Ints(labelInts)
	for _, labelInt := range labelInts {
		label := strconv.Itoa(labelInt)
		fn := funcMap[label]
		pointer := reflect.ValueOf(fn).Pointer()
		fnFullName := runtime.FuncForPC(pointer).Name()
		fnName := strings.Split(fnFullName, ".")[1]
		fmt.Printf("%s: %s\n", label, fnName)
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
