package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// promptString prompts the user for a string input.
// If required is true, it loops until a non-empty string is provided.
func promptString(label string, required bool) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		if required {
			fmt.Printf("%s (*required): ", label)
		} else {
			fmt.Printf("%s (optional): ", label)
		}
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			return input
		}
		if !required {
			return ""
		}
	}
}

// promptInt prompts the user for an integer input.
// If required is true, it loops until a valid integer is provided.
// If not required and input is empty, it returns 0.
func promptInt(label string, required bool) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		if required {
			fmt.Printf("%s (*required): ", label)
		} else {
			fmt.Printf("%s (optional): ", label)
		}
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" && !required {
			return 0
		}
		val, err := strconv.Atoi(input)
		if err == nil {
			return val
		}
		fmt.Println("Invalid number, please try again.")
	}
}
