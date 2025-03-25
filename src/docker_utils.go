package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// checkDocker checks if Docker is installed on the system and verifies the presence of Docker images
// corresponding to the provided challenges. It writes status messages to the provided writer.
//
// Parameters:
//   - challenges: A slice of Challenge structs representing the challenges to verify.
//   - writer: A bufio.Writer to write status messages to.
//
// Returns:
//   - A slice of Challenge structs representing the challenges that have corresponding Docker images.
//
// The function performs the following steps:
//  1. Checks if Docker is installed by running the "which docker" command.
//  2. If Docker is not installed, it writes an error message to the writer and exits the program.
//  3. If Docker is installed, it retrieves the path to the Docker executable.
//  4. Lists the Docker images available on the system using the "docker image list" command.
//  5. Iterates over the provided challenges and checks if each challenge's Docker image is present.
//  6. Writes a message to the writer for each missing Docker image.
//  7. Returns a slice of challenges that have corresponding Docker images.
func checkDocker(challenges []Challenge, writer *bufio.Writer) []Challenge {
	out, err := exec.Command("/usr/bin/which", "docker").Output()
	check(err, "Failed to check if Docker is installed", true, writer)

	if string(out) == "" {
		fmt.Fprintln(writer, "\033[31mDocker is not installed on the system.\033[0m")
		os.Exit(1)
	} else {
		fmt.Fprintln(writer, "Docker is installed on the system. \033[32m✔\033[0m")
	}

	dockerPath := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
	fmt.Fprintln(writer, "Docker is installed at:", dockerPath)

	outImages, err := exec.Command(dockerPath, "image", "list", "--format", "{{.Repository}}").Output()
	check(err, "Failed to list Docker images", true, writer)

	var validChallenges []Challenge
	for _, challenge := range challenges {
		if strings.Contains(string(outImages), challenge.Shortname) {
			validChallenges = append(validChallenges, challenge)
		} else {
			fmt.Fprintf(writer, "\033[31mThe docker image: %s is missing.\033[0m\n", challenge.Shortname)
		}
	}

	return validChallenges
}

// startContainer starts a Docker container for the given challenge index.
// It takes the following parameters:
// - index: The index of the challenge in the challenges slice.
// - challenges: A slice of Challenge structs containing information about each challenge.
// - status: A slice of strings representing the status of each challenge (not used in this function).
// - writer: A bufio.Writer to write output messages.
//
// If the index is out of range, an error message is written to the writer.
// The function constructs a Docker run command to start the container with the appropriate port mapping
// and container name based on the challenge information. If the command fails, an error message is written
// to the writer.
func startContainer(index uint, challenges []Challenge, status []string, writer *bufio.Writer) {
	if index >= uint(len(challenges)) {
		fmt.Fprintln(writer, "\033[31mERROR\033[0m: Index out of range.")
		return
	}

	cmd := exec.Command("docker", "run", "-p",
		fmt.Sprintf("%d:%d", challenges[index].Exposed_port, challenges[index].Exposed_port),
		, "--rm", "--name", challenges[index].Shortname,
		challenges[index].Shortname,
	)
	err := cmd.Run()
	
	check(err, fmt.Sprintf("Failed to start container %s.", challenges[index].Shortname), false, writer)
}

// stopContainer stops a Docker container based on the provided index in the challenges slice.
// It writes an error message to the writer if the index is out of range or if the container fails to stop.
//
// Parameters:
//   - index: The index of the container in the challenges slice to stop.
//   - challenges: A slice of Challenge structs containing information about the containers.
//   - status: A slice of strings representing the status of each container (not used in this function).
//   - writer: A bufio.Writer to write error messages to.
//
// The function uses the "docker container restart" command to restart the specified container.
// If the command fails, it calls the check function to handle the error.
func stopContainer(index uint, challenges []Challenge, status []string, writer *bufio.Writer) {
	if index >= uint(len(challenges)) {
		fmt.Fprintln(writer, "\033[31mERROR\033[0m: Index out of range.")
		return
	}

	cmd := exec.Command("docker", "container", "stop", challenges[index].Shortname)
	err := cmd.Run()
	check(err, fmt.Sprintf("Failed to stop container %s.", challenges[index].Shortname), false, writer)
}

// restartContainer restarts a Docker container based on the provided index in the challenges slice.
// If the index is out of range, it writes an error message to the provided writer.
//
// Parameters:
//   - index: The index of the container to restart in the challenges slice.
//   - challenges: A slice of Challenge structs containing information about the containers.
//   - status: A slice of strings representing the status of each container (not used in the function).
//   - writer: A bufio.Writer to write output messages.
//
// The function uses the "docker container restart" command to restart the specified container.
// If the command fails, it calls the check function to handle the error.
func restartContainer(index uint, challenges []Challenge, status []string, writer *bufio.Writer) {
	if index >= uint(len(challenges)) {
		fmt.Fprintln(writer, "\033[31mERROR\033[0m: Index out of range.")
		return
	}

	cmd := exec.Command("docker", "container", "restart", challenges[index].Shortname)
	err := cmd.Run()
	check(err, fmt.Sprintf("Failed to restart container %s.", challenges[index].Shortname), false, writer)
}
