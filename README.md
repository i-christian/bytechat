# ByteChat 🚀💬

![GitHub repo size](https://img.shields.io/github/repo-size/i-christian/bytechat?style=flat-square)
![GitHub license](https://img.shields.io/github/license/i-christian/bytechat?style=flat-square)
[![Build Status](https://github.com/i-christian/bytechat/actions/workflows/tests.yml/badge.svg)](https://github.com/i-christian/bytechat/actions/workflows/tests.yml)
[![GitHub Issues](https://img.shields.io/github/issues/i-christian/bytechat)](https://github.com/i-christian/bytechat/issues)
[![GitHub Contributors](https://img.shields.io/github/contributors/i-christian/bytechat)](https://github.com/i-christian/bytechat/graphs/contributors)
[![Last Commit](https://img.shields.io/github/last-commit/i-christian/bytechat)](https://github.com/i-christian/bytechat/commits/main)
[![Latest Release](https://img.shields.io/github/v/release/i-christian/bytechat?include_prereleases)](https://github.com/i-christian/bytechat/releases)


## Description
ByteChat is a real-time chat application demonstrating the use of WebSockets in a Go backend. The application features a server-side rendered frontend using `templ` and HTMX for dynamic updates.

The primary motivation for this project was to explore and implement a clean WebSocket solution in Go for real-time communication, supporting both public chat rooms and private(peer to peer) interactions.

## Key Features ✨

* **Real-Time Communication**: Built using webSockets for instant, bidirectional messaging between users.
* **Go Backend**: Leverages Go's performance, concurrency features, and strong typing for a reliable backend.
* **Server-Side HTML Rendering**: Uses `templ` for type-safe Go components rendered to HTML on the server.
* **Dynamic Frontend Updates**: Utilizes HTMX for WebSocket communication and seamless partial page updates without full page reloads.
* **Secure User Authentication**: Implements session management using cryptographically signed, HTTP-only cookies.
* **Message Persistence**: Stores user data, room information, and chat messages in a PostgreSQL database.
* **Database Migrations**: Uses `goose` for managing database schema changes.
* **Dockerized Deployment**: Comes with `docker-compose` for easy setup and deployment.

## Tech Stack 🛠️

* **Backend**: Go (Golang)
* **Web Framework/Routing**: `https://go-chi.io/#/`
* **WebSockets**: `https://github.com/coder/websocket` (websocket)
* **Templating**: `templ` (`https://templ.guide/`)
* **Frontend Interaction**: HTMX (`https://htmx.org/docs/#introduction`)
* **Database**: PostgreSQL
* **Migrations**: `goose` (`https://github.com/pressly/goose`)
* **Containerization**: Docker & Docker Compose

## Prerequisites 📋

* Go (version 1.21 or later recommended)
* PostgreSQL Server
* Docker & Docker Compose
* `make` (for using Makefile commands)
* `goose` CLI (The application has manages migration automatically)

## Getting Started 🚀

### 1. Clone the Repository

```bash
git clone https://github.com/i-christian/bytechat.git
cd bytechat
```

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

### Running with Docker 🐳
- **Build and Start**: Ensure Docker is installed, then run the following command to build and start the services.
```
  docker compose up -d
```

- **Access the Application**: Once the container is running, access ByteChat at:
```
  http://localhost:3000
```
- **Stop the Service**: Run the following command to stop the service:
```
  docker compose down --remove-orphans
```

## Contributing 🤝

I welcome contributions to improve ByteChat. Here’s how you can get started:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.

## Licensing 📄
ByteChat is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
