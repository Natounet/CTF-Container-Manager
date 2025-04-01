package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// handleConnection manages the interaction with a connected client over a network connection.
// It authenticates the client using a secret key, logs connection attempts, and provides a menu
// for managing Docker containers.
//
// Parameters:
//   - conn: The network connection to the client.
//   - validChallenges: A slice of Challenge structs representing the Docker containers.
//   - secretKey: A byte slice containing the secret key for authentication.
//   - logFile: A file pointer to the log file for recording connection and container management events.
//
// The function performs the following steps:
//  1. Prompts the client to enter the secret key and validates it.
//  2. Logs the connection attempt (successful or unsuccessful).
//  3. Displays a menu for the client to start, stop, or restart Docker containers.
//  4. Logs each container management action performed by the client.
//  5. Handles client disconnection and logs the event.
func handleConnection(conn net.Conn, validChallenges []Challenge, secretKeyBytes []byte, logFile *os.File) {
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	logWriter := bufio.NewWriter(logFile)

	fmt.Fprint(writer, "Enter the secret key: ")
	writer.Flush()
	inputKey, _ := reader.ReadString('\n')
	inputKey = strings.TrimSpace(inputKey)
	secretKeyBytes = []byte(strings.TrimSpace(string(secretKeyBytes)))
	var inputKeyBytes []byte = []byte(inputKey)

	// Timing independant comparison to prevent timing attacks
	if !isSecretValid(inputKeyBytes, secretKeyBytes) {
		fmt.Fprintln(writer, "ERROR: Invalid secret key.")
		logWriter.WriteString(fmt.Sprintf("Invalid key attempt from %s\n", conn.RemoteAddr().String()))
		logWriter.Flush()
		writer.Flush()
		return
	}

	logWriter.WriteString(fmt.Sprintf("Successful connection from %s\n", conn.RemoteAddr().String()))
	logWriter.Flush()

	fmt.Fprintln(writer, "Press Enter to access the console")
	writer.Flush()
	reader.ReadString('\n')

	for {
		status := statusChallenges(validChallenges, writer)
		fmt.Fprint(writer, "\033[H\033[2J") // Clear screen
		listContainers(validChallenges, status, writer)

		fmt.Fprintln(writer, "Menu:")
		fmt.Fprintln(writer, "1. Start Container")
		fmt.Fprintln(writer, "2. Stop Container")
		fmt.Fprintln(writer, "3. Restart Container")
		fmt.Fprintln(writer, "4. Start all containers")
		fmt.Fprintln(writer, "5. Stop all containers")
		fmt.Fprintln(writer, "6. Restart all containers")
		fmt.Fprintln(writer, "0. Exit")
		fmt.Fprint(writer, "> ")
		writer.Flush()

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Fprint(writer, "Enter the index of the container to start: ")
			writer.Flush()
			index, err := checkInput(reader, writer)
			if err != nil {
				break
			}
			startContainer(index, validChallenges, status, writer)
			logWriter.WriteString(fmt.Sprintf("Started container %d by %s\n", index, conn.RemoteAddr().String()))

		case "2":
			fmt.Fprint(writer, "Enter the index of the container to stop: ")
			writer.Flush()
			index, err := checkInput(reader, writer)
			if err != nil {
				break
			}
			stopContainer(index, validChallenges, status, writer)
			logWriter.WriteString(fmt.Sprintf("Stopped container %d by %s\n", index, conn.RemoteAddr().String()))

		case "3":
			fmt.Fprint(writer, "Enter the index of the container to restart: ")
			writer.Flush()
			index, err := checkInput(reader, writer)
			if err != nil {
				break
			}
			restartContainer(index, validChallenges, status, writer)
			logWriter.WriteString(fmt.Sprintf("Restarted container %d by %s\n", index, conn.RemoteAddr().String()))

		case "4":
			fmt.Fprintln(writer, "Starting all containers...")
			startAllContainers(validChallenges, status, writer)
			logWriter.WriteString(fmt.Sprintf("Started all containers by %s\n", conn.RemoteAddr().String()))

		case "5":
			fmt.Fprintln(writer, "Stopping all containers...")
			stopAllContainers(validChallenges, status, writer)
			logWriter.WriteString(fmt.Sprintf("Stopped all containers by %s\n", conn.RemoteAddr().String()))

		case "6":
			fmt.Fprintln(writer, "Restarting all running containers...")
			restartAllContainers(validChallenges, status, writer)
			logWriter.WriteString(fmt.Sprintf("Stopped all containers by %s\n", conn.RemoteAddr().String()))

		case "0":
			fmt.Fprintln(writer, "Goodbye!")
			logWriter.WriteString(fmt.Sprintf("Connection closed by %s\n", conn.RemoteAddr().String()))

			return
		default:
			fmt.Fprintln(writer, "\033[31mInvalid choice. Please try again.\033[0m")
			writer.Flush()
		}

		fmt.Fprintln(writer, "Press Enter to continue...")
		writer.Flush()
		logWriter.Flush()
		reader.ReadString('\n')
	}
}
