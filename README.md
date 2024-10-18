# ByteChat 🚀💬

![GitHub repo size](https://img.shields.io/github/repo-size/i-christian/bytechat?style=flat-square)
![GitHub license](https://img.shields.io/github/license/i-christian/bytechat?style=flat-square)

## Description
ByteChat is a real-time chat application backend built with Rust, leveraging the power of Axum and SocketIO. SocketIO provides a convenient abstraction over WebSockets, enabling seamless and efficient real-time communication between clients and the server. This setup allows for instant messaging capabilities, making it ideal for chat applications where low-latency interaction is crucial.

The backend is designed to handle multiple concurrent connections, manage chat messages, and broadcast messages to all connected clients in real-time. By using Rust, known for its performance and safety, the backend ensures robust and reliable handling of numerous simultaneous connections without compromising on speed.

### Key Features
- **Real-Time Communication**: Utilizing SocketIO over WebSockets for instant, bidirectional communication.
- **Concurrency Handling**: Rust’s async capabilities, powered by Tokio, manage numerous concurrent connections efficiently.
- **Message Broadcasting**: Broadcasts messages to all connected clients, perfect for group chats.
- **Robust and Scalable**: Built with Rust and Axum, ensuring a fast, secure, and scalable backend.
- **Extensible Architecture**: Modular design allows for easy integration of additional features such as authentication, message persistence, and user management.

## Prerequisites 📋
- [Rust](https://www.rust-lang.org/)
- [Axum](https://docs.rs/axum/latest/axum/)
- [socketioxide](https://docs.rs/socketioxide/latest/socketioxide/)

## How to Use 🛠️

### Clone the Repository
clone the `bytechat` repository to your local machine:
```
git clone https://github.com/i-christian/bytechat.git

cd bytechat
```

### Run the application
```
cargo run
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
