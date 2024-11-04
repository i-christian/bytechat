# ByteChat 🚀💬

![GitHub repo size](https://img.shields.io/github/repo-size/i-christian/bytechat?style=flat-square)
![GitHub license](https://img.shields.io/github/license/i-christian/bytechat?style=flat-square)
[![Build Status](https://github.com/i-christian/bytechat/actions/workflows/tests.yml/badge.svg)](https://github.com/i-christian/bytechat/actions/workflows/tests.yml)
[![GitHub Issues](https://img.shields.io/github/issues/i-christian/bytechat)](https://github.com/i-christian/bytechat/issues)
[![GitHub Contributors](https://img.shields.io/github/contributors/i-christian/bytechat)](https://github.com/i-christian/bytechat/graphs/contributors)
[![Last Commit](https://img.shields.io/github/last-commit/i-christian/bytechat)](https://github.com/i-christian/bytechat/commits/main)
[![Latest Release](https://img.shields.io/github/v/release/i-christian/bytechat?include_prereleases)](https://github.com/i-christian/bytechat/releases)


## Description
ByteChat is a real-time chat application with a backend built in Rust (using Axum and SocketIO) and a frontend powered by SolidJS. The backend leverages Rust’s high-performance capabilities, while SocketIO provides an efficient WebSocket abstraction for low-latency, bidirectional communication.

ByteChat is designed for scalability and concurrency, ideal for chat applications where multiple users require instant messaging. The backend handles numerous simultaneous connections, ensuring speed and reliability with Rust’s memory safety and performance.

### Key Features
- **Real-Time Communication**: Utilizing SocketIO over WebSockets for instant, bidirectional communication.
- **Concurrency Handling**: Rust’s async capabilities, powered by Tokio, manage numerous concurrent connections efficiently.
- **Message Broadcasting**: Broadcasts messages to all connected clients, perfect for group chats.
- **Robust and Scalable**: Built with Rust and Axum, ensuring a fast, secure, and scalable backend.
- **Extensible Architecture**: Modular design allows for easy integration of additional features such as authentication, message persistence, and user management
- **Dockerized Deployment**: Easily deploy ByteChat with Docker for a consistent environment.

## Prerequisites 📋
- [Rust](https://www.rust-lang.org/)
- [Axum](https://docs.rs/axum/latest/axum/)
- [socketioxide](https://docs.rs/socketioxide/latest/socketioxide/)
- [Node.js & npm](https://nodejs.org/)
- [Docker](https://www.docker.com/)

## How to Use 🛠️

### Clone the Repository
clone the `bytechat` repository to your local machine:
```
git clone https://github.com/i-christian/bytechat.git

cd bytechat
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
./buid.sh or cargo run
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

## Acknowledgements 🙌
Special thanks to the developers of Rust, Axum, and socketioxide for their excellent tools and libraries.

Enjoy chatting with ByteChat! 🚀💬
