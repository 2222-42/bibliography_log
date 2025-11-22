package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	reader     *bufio.Reader
	readerOnce sync.Once
)

func getReader() *bufio.Reader {
	readerOnce.Do(func() {
		reader = bufio.NewReader(os.Stdin)
	})
	return reader
}

// promptString prompts the user for a string input.
// If required is true, it loops until a non-empty string is provided.
func promptString(label string, required bool) string {
	r := getReader()
	for {
		if required {
			fmt.Printf("%s (*required): ", label)
		} else {
			fmt.Printf("%s (optional): ", label)
		}
		input, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// If we got some input before EOF, use it
				if strings.TrimSpace(input) != "" {
					fmt.Println() // Add newline for cleaner output
					return strings.TrimSpace(input)
				}
				// If EOF with no input and optional, return empty
				if !required {
					return ""
				}
				// If EOF with no input and required, exit with error
				fmt.Printf("\nError: required input for '%s' not provided before EOF\n", label)
				os.Exit(1)
			}
			fmt.Printf("\nError reading input: %v\n", err)
			os.Exit(1)
		}
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
	r := getReader()
	for {
		if required {
			fmt.Printf("%s (*required): ", label)
		} else {
			fmt.Printf("%s (optional): ", label)
		}
		input, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if strings.TrimSpace(input) != "" {
					fmt.Println()
					val, convErr := strconv.Atoi(strings.TrimSpace(input))
					if convErr == nil {
						return val
					}
				}
				if !required {
					return 0
				}
				// If required and EOF with no valid input, exit with error
				fmt.Printf("\nError: required integer input for '%s' but reached end of input.\n", label)
				os.Exit(1)
			}
			fmt.Printf("\nError reading input: %v\n", err)
			os.Exit(1)
		}
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
