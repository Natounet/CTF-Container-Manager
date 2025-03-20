package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: ./dockerLister <secretKeyFile> <challengesFile> <IP> <Port>")
		os.Exit(1)
	}

	secretKeyFile := os.Args[1]
	challengesFile := os.Args[2]
	ip := os.Args[3]
	port := os.Args[4]

	writer := bufio.NewWriter(os.Stdout)

	secretKey, err := os.ReadFile(secretKeyFile)
	check(err, "Failed to read secret key file", true, writer)

	challenges := loadChallenges(challengesFile, writer)
	verifyChallenges(challenges, writer)
	validChallenges := checkDocker(challenges, writer)

	fmt.Fprintln(writer, "Available challenges:")

	listContainers(validChallenges, statusChallenges(validChallenges, writer), writer)
	writer.Flush()

	listener, err := net.Listen("tcp", ip+":"+port)
	check(err, "Failed to start TCP server.", true, bufio.NewWriter(os.Stdout))
	defer listener.Close()

	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, "Failed to open log file", true, writer)
	defer logFile.Close()

	fmt.Printf("Server listening on %s:%s\n", ip, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, validChallenges, secretKey, logFile)
	}
}
