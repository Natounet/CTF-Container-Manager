package main

import (
	"bufio"
	"crypto/subtle"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// check handles error checking and logging.
// If an error is present, it writes the provided message and error to the given writer in red color.
// If abort is true, the function will terminate the program.
// It returns false if an error is present, otherwise true.
//
// Parameters:
//   - e: The error to check.
//   - message: The message to log if an error is present.
//   - abort: Whether to terminate the program if an error is present.
//   - writer: The writer to which the message and error will be logged.
//
// Returns:
//   - bool: False if an error is present, otherwise true.
func check(e error, message string, abort bool, writer *bufio.Writer) bool {
	if e != nil {
		fmt.Fprintf(writer, "\033[31m%s\033[0m\n", message)
		fmt.Fprintf(writer, "\033[31m%s\033[0m\n", e.Error())
		writer.Flush()
		if abort {
			os.Exit(1)
		}
		return false
	}
	return true
}

// checkInput reads a line of input from the provided bufio.Reader, trims any
// surrounding whitespace, and attempts to parse it as an unsigned integer.
// If the input is successfully parsed, it returns the parsed value as a uint.
// If an error occurs during reading or parsing, it writes an error message to
// the provided bufio.Writer, flushes the writer, and returns an error.
//
// Parameters:
// - reader: A pointer to a bufio.Reader from which the input will be read.
// - writer: A pointer to a bufio.Writer to which error messages will be written.
//
// Returns:
// - uint: The parsed unsigned integer if the input is valid.
// - error: An error if reading or parsing the input fails.
func checkInput(reader *bufio.Reader, writer *bufio.Writer) (uint, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(writer, "\033[31mERROR\033[0m: Failed to read input.")
		writer.Flush()
		return 0, err
	}

	input = strings.TrimSpace(input)
	index, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		fmt.Fprintln(writer, "\033[31mERROR\033[0m: Invalid input. Please enter a valid number.")
		writer.Flush()
		return 0, err
	}

	return uint(index), nil
}

// listContainers prints a formatted list of Docker containers to the provided writer.
// It takes a slice of Challenge structs, a slice of status strings, and a bufio.Writer as arguments.
//
// Parameters:
//   - challenges: A slice of Challenge structs containing information about each Docker container.
//   - status: A slice of strings representing the status of each Docker container.
//   - writer: A bufio.Writer to which the formatted list will be written.
//
// The output format includes the following columns:
//   - N°: The index number of the container in the list.
//   - Fullname: The full name of the container.
//   - Shortname: The short name of the container.
//   - Exposed Port: The port exposed by the container.
//   - Status: The status of the container.
//
// Example output:
//
//	N°   | Fullname                      | Shortname            | Exposed Port | Status
//	--------------------------------------------------------------------------------
//	0    | example_fullname_1            | example_shortname_1  | 8080         | running
//	1    | example_fullname_2            | example_shortname_2  | 9090         | stopped
func listContainers(challenges []Challenge, status []string, writer *bufio.Writer) {
	fmt.Fprintf(writer,
		"%-4s | %-30s | %-20s | %-12s | %s\n",
		"N°", "Fullname", "Shortname", "Exposed Port", "Status",
	)
	fmt.Fprintln(writer,
		strings.Repeat("-", 80),
	)
	for i, challenge := range challenges {
		fmt.Fprintf(
			writer,
			"%-4d | %-30s | %-20s | %-12d | %s\n",
			i,
			challenge.Fullname,
			challenge.Shortname,
			challenge.Exposed_port,
			status[i],
		)
	}
	writer.Flush()
}

// statusChallenges retrieves the status of Docker containers for a given list of challenges.
// It uses the `docker container list` command to get the names and statuses of running containers,
// and matches them with the provided challenges to determine their status.
//
// Parameters:
//   - challenges: A slice of Challenge structs, each representing a challenge with a Shortname field.
//   - writer: A bufio.Writer to write error messages if any command execution fails.
//
// Returns:
//   - A slice of strings representing the status of each challenge. The status is color-coded:
//   - Green (\033[32m) if the container is running (status contains "up").
//   - Red (\033[31m) if the container is not running or not found.
func statusChallenges(challenges []Challenge, writer *bufio.Writer) []string {
	outNames, err := exec.Command("docker", "container", "list", "--format", "{{.Image}}").Output()
	check(err, "Failed to list Docker containers", true, writer)

	outStatus, err := exec.Command("docker", "container", "list", "--format", "{{.Status}}").Output()
	check(err, "Failed to get Docker container status", true, writer)

	containerNames := strings.Split(strings.TrimSpace(string(outNames)), "\n")
	containerStatuses := strings.Split(strings.TrimSpace(string(outStatus)), "\n")

	var statuses []string
	for _, challenge := range challenges {
		found := false
		for j, name := range containerNames {
			if strings.Contains(name, challenge.Shortname) {
				status := containerStatuses[j]
				if strings.Contains(strings.ToLower(status), "up") {
					status = "\033[32m" + strings.Trim(status, "'") + "\033[0m"
				} else {
					status = "\033[31m" + strings.Trim(status, "'") + "\033[0m"
				}
				statuses = append(statuses, status)
				found = true
				break
			}
		}
		if !found {
			statuses = append(statuses, "\033[31mNot Running\033[0m")
		}
	}

	return statuses
}

// isSecretValid compares two byte slices, `input_bytes` and `secret_bytes`,
// using a constant-time comparison to prevent timing attacks. It returns
// true if the two slices are equal, and false otherwise.
//
// Parameters:
//   - input_bytes: The first byte slice to compare.
//   - secret_bytes: The second byte slice to compare.
//
// Returns:
//   - bool: True if the byte slices are equal, false otherwise.
func isSecretValid(input_bytes []byte, secret_bytes []byte) bool {

	return subtle.ConstantTimeCompare(input_bytes, secret_bytes) == 1
}
