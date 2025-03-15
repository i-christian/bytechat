# ByteChat 🚀💬

![GitHub repo size](https://img.shields.io/github/repo-size/i-christian/bytechat?style=flat-square)
![GitHub license](https://img.shields.io/github/license/i-christian/bytechat?style=flat-square)
[![Build Status](https://github.com/i-christian/bytechat/actions/workflows/tests.yml/badge.svg)](https://github.com/i-christian/bytechat/actions/workflows/tests.yml)
[![GitHub Issues](https://img.shields.io/github/issues/i-christian/bytechat)](https://github.com/i-christian/bytechat/issues)
[![GitHub Contributors](https://img.shields.io/github/contributors/i-christian/bytechat)](https://github.com/i-christian/bytechat/graphs/contributors)
[![Last Commit](https://img.shields.io/github/last-commit/i-christian/bytechat)](https://github.com/i-christian/bytechat/commits/main)
[![Latest Release](https://img.shields.io/github/v/release/i-christian/bytechat?include_prereleases)](https://github.com/i-christian/bytechat/releases)


## Description
ByteChat is a real-time chat application with a backend built in Rust (using Axum and SocketIO) and a frontend powered by SolidJS. The backend leverages Golang’s high-performance capabilities, it also utilises WebSockets to offer a realtime communication among users of the application.

ByteChat is designed for scalability and concurrency, ideal for chat applications where multiple users require instant messaging. The backend handles numerous simultaneous connections, ensuring speed and reliability with Rust’s memory safety and performance.

### Key Features
- **Real-Time Communication**: Utilizing WebSockets for instant realtime, bidirectional communication.
- **Concurrency Handling**: Golang's concurrency features.
- **User Authentication**: The application utilises http only privately signed cookies to offer robust a session management and user authentication.
- **Message Persistence**: The application uses postgresql to save messages, and user information.
- **Robust and Scalable**: Built with Golang, ensuring a fast, secure, and scalable backend.
- **Dockerized Deployment**: Easily deploy ByteChat with Docker for a consistent environment.

## Prerequisites 📋
- [The WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)

## Getting Started 🛠️

### Clone the Repository
clone the `bytechat` repository to your local machine:
```
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

### Run the application (without Docker)
```
./buid.sh
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
