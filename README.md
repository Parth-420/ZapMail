# SMTP Server in Go

## Overview

This project is a lightweight and efficient SMTP server written from scratch in Go. It is designed to handle incoming emails and process SMTP commands efficiently, making it suitable for learning, experimentation, or lightweight mail-handling applications.

## Features

- Fully implemented SMTP protocol from scratch.
- Handles common SMTP commands such as `HELO`, `MAIL FROM`, `RCPT TO`, `DATA`, and `QUIT`.
- Supports multiple concurrent client connections.
- Logs incoming email transactions.
- Configurable settings for port and hostname.
- Written in pure Go without external dependencies.

## Prerequisites

- Go 1.20 or later installed.
- A machine running macOS, Linux, or Windows.

## Installation

1. Clone this repository:
   ```sh
   git clone https://github.com/yourusername/ZapMail.git
   cd ZapMail
   ```
2. ZBuild the server:
   ```sh
   go build -o  ZapMail main.go
   ```
3. Run the server:
   ```sh
   ./ZapMail
   ```

## Configuration

You can customize the SMTP server settings via environment variables:

- `SMTP_HOST`: The hostname or IP address to bind the server (default: `localhost`).
- `SMTP_PORT`: The port number for the server to listen on (default: `2525`).

Example:

```sh
SMTP_HOST=0.0.0.0 SMTP_PORT=2525 ./smtp-server
```

## Usage

Once the server is running, you can connect to it using an SMTP client like `telnet`:

```sh
telnet localhost 2525
```

Then, interact with the server using standard SMTP commands:

```sh
HELO example.com
MAIL FROM:<sender@example.com>
RCPT TO:<receiver@example.com>
DATA
Subject: Test Email

This is a test email.
.
QUIT
```

## Roadmap

- Add support for authentication (`AUTH LOGIN` and `AUTH PLAIN`).
- Implement email forwarding and relaying.
- Store received emails in a queue or database.
- Support for TLS encryption (`STARTTLS`).

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository.
2. Create a new feature branch.
3. Commit your changes and push to your fork.
4. Submit a pull request.

##

