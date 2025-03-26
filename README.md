# CTF Container Manager

CTF Container Manager is a lightweight tool designed to streamline the management of Docker-based challenges for Capture The Flag \(CTF\) events. It provides a remote interface for teams hosting challenges to start, stop, and restart containers via a TCP/TLS connection using `openssl`.

---

![image](https://github.com/user-attachments/assets/dc3b502d-7098-4d13-bc0c-9c322a1221c5)


## Features

- **Remote Container Management**: Start, stop, and restart Docker containers remotely.
- **Challenge Configuration**: Define challenges in a JSON file with details such as container name, exposed port, etc.
- **Authentication**: Access with a printable ASCII passphrase stored in `secret.key`.
- **Status Monitoring**: View the status of all configured containers (running or stopped).
- **Logging**: Logs all connections and container management actions for auditing purposes.
- **TLS Encryption**: All network communications are secured using Transport Layer Security (TLS)

---

## Installation

1. Clone the repository containing the project files:
```bash
git clone https://github.com/Natounet/CTF-Container-Manager.git
cd CTF-Container-Manager/src
```

2. Build the project:
```bash
go build -o manager
```

4. Ensure Docker is installed and accessible on the host machine.

---

## Usage

### Starting the Server
Run the server with the following command:

```bash
./ctf-container-manager <secret.key> <challenges.json> <IP> <Port>
```

- `<secret.key>`: Path to the file containing the secret key for authentication.
- `<challenges.json>`: Path to the JSON file describing challenges.
- `<IP>`: The IP address on which the server will listen.
- `<Port>`: The port number for client connections.

Example:

```bash
./ctf-container-manager example_secret.key example_challenges.json 127.0.0.1 9000
```

![image](https://github.com/user-attachments/assets/36810a20-6d1a-41b8-91df-6b78f3ec2372)


---

### Challenge Configuration

Challenges are described in a JSON file \(`example_challenges.json`\). Each challenge includes:
- `fullname`: A descriptive name for the challenge.
- `shortname`: The name of the Docker image \(must exist locally\).
- `exposed_port`: The port exposed by the container for players.

```javascript
[
{"fullname": "CyberPhoenix", "shortname": "cyberphoenix", "exposed_port": 9000},
{"fullname": "CryptoMaze", "shortname": "cryptomaze", "exposed_port": 5678}
]
```

---

### Client Interaction

Clients can connect to the server using `openssl` since the server use TLS:

```bash
openssl s_client -connect <IP>:<PORT> -quiet
```
![image](https://github.com/user-attachments/assets/b8b93b91-db11-4b9c-bfb5-5ada7e154a8b)


Upon connection:
1. Enter the secret key for authentication.
2. Access a menu to manage containers:
   - Start Container
   - Stop Container
   - Restart Container
   - Exit

The server will display available challenges and their statuses.

---

## Security Considerations

- Ensure `secret.key` is securely stored and accessible only by authorized users.
- Use strong passwords in `secret.key` to prevent unauthorized access.

---

## Logs

All connection attempts and container management actions are logged in `server.log`. This includes:
- Successful/failed authentication attempts.
- Actions performed (start, stop, restart) along with timestamps and client IPs.

---

## Requirements

- Go programming language installed (`>= v1.21`).
- Docker installed and running on the host machine.
- Access to TCP ports for client-server communication.

---

---

## Troubleshooting

### Common Errors
1. **Docker Not Installed**:
   Ensure Docker is installed on your system and accessible via CLI.
   
2. **Missing Docker Images**:
   Verify that all images specified in `example_challenges.json` exist locally using:

3. **Duplicate Ports**:
Ensure each challenge has a unique exposed port in `example_challenges.json`.

4. **Invalid Secret Key**:
Verify that clients are using the correct key stored in `secret.key`.

---


