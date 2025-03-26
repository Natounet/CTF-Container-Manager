package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: ./ctf-container-manager <secretKeyFile> <challengesFile> <IP> <Port>")
		os.Exit(1)
	}

	secretKeyFile := os.Args[1]
	challengesFile := os.Args[2]
	ip := os.Args[3]
	port := os.Args[4]

	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, "Failed to open application log file", true, nil)
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	writer := bufio.NewWriter(multiWriter)

	secretKey, err := os.ReadFile(secretKeyFile)
	check(err, "Failed to read secret key file", true, writer)

	challenges := loadChallenges(challengesFile, writer)
	verifyChallenges(challenges, writer)
	validChallenges := checkDocker(challenges, writer)

	fmt.Fprintln(writer, "Available challenges:")

	listContainers(validChallenges, statusChallenges(validChallenges, writer), writer)
	writer.Flush()

	// Générer les certificats TLS
	cert, err := generateCertificates()
	check(err, "Failed to generate TLS certificates", true, writer)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Écouter avec TLS
	listener, err := tls.Listen("tcp", ip+":"+port, tlsConfig)
	check(err, "Failed to start TLS server.", true, writer)
	defer listener.Close()

	check(err, "Failed to open log file", true, writer)
	defer logFile.Close()

	fmt.Printf("Server listening securely on %s:%s\n", ip, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, validChallenges, secretKey, logFile)
	}
}
