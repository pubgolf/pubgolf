package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getInput(label string, defaultValue string) string {
	buf := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter value for argument '%s' (default: '%s') >>> ", label, defaultValue)
	input, err := buf.ReadBytes('\n')
	if err != nil {
		log.Printf("Could not read input: %s", err)
	}
	if len(input) == 1 {
		return defaultValue
	}
	return string(input[:len(input)-1])
}

func getInputAsUInt32(label string, defaultValue uint32) (uint32, error) {
	num, err := strconv.Atoi(getInput(label, string(defaultValue)))
	return uint32(num), err
}

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

func getCmd(r *bufio.Reader) string {
	t, _ := r.ReadString('\n')
	return strings.TrimSpace(t)
}
