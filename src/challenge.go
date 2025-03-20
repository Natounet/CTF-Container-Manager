package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Challenge struct {
	Fullname     string
	Shortname    string
	Exposed_port uint16
}

// loadChallenges reads a JSON file containing challenges, unmarshals the data into a slice of Challenge structs,
// and returns the slice. It takes a filename as a string and a bufio.Writer for logging errors.
// If an error occurs during reading or unmarshalling, it logs the error and exits the program.
//
// Parameters:
//   - filename: The path to the JSON file containing the challenges.
//   - writer: A bufio.Writer used for logging errors.
//
// Returns:
//   - A slice of Challenge structs containing the unmarshalled data from the JSON file.
func loadChallenges(filename string, writer *bufio.Writer) []Challenge {
	challengesBytes, err := os.ReadFile(filename)
	check(err, "Failed to read challenges file", true, writer)

	var challenges []Challenge
	err = json.Unmarshal(challengesBytes, &challenges)
	check(err, "Failed to unmarshal challenges JSON", true, writer)

	return challenges
}

// verifyChallenges checks for duplicate exposed ports in the given list of challenges.
// If a duplicate port is found, it writes an error message to the provided writer and exits the program.
//
// Parameters:
//   - challenges: A slice of Challenge structs to be verified.
//   - writer: A bufio.Writer to write error messages to.
//
// The function maintains a map to track exposed ports and ensures no duplicates exist.
// If a duplicate port is detected, an error message is printed in red, and the program terminates.
func verifyChallenges(challenges []Challenge, writer *bufio.Writer) {
	portMap := make(map[uint16]bool)
	for _, challenge := range challenges {
		if portMap[challenge.Exposed_port] {
			fmt.Fprintf(writer, "\033[31mERROR\033[0m: Duplicate exposed port found: %d\n", challenge.Exposed_port)
			os.Exit(1)
		}
		portMap[challenge.Exposed_port] = true
	}
}
