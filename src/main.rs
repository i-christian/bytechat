mod state;

use axum::{http::StatusCode, routing::get};
use socketioxide::{
    extract::{Data, SocketRef, State},
    SocketIo,
};
use tower::ServiceBuilder;
use tower_http::cors::CorsLayer;
use tracing::info;

#[derive(Debug, serde::Deserialize)]
struct MessageIn {
    room: String,
    text: String,
}

#[derive(serde::Serialize)]
struct Messages {
    messages: Vec<state::Message>,
}

/// Handles socket connection events
///
/// This function is called when a new socket connection is established. It sets up
/// listeners for "join" and "message" events on the socket.
///
/// # Arguments
///
/// * `socket` - A reference to the connected socket
async fn on_connect(socket: SocketRef) {
    info!("socket connected: {}", socket.id);

    socket.on(
        "join",
        |socket: SocketRef, Data::<String>(room), store: State<state::MessageStore>| async move {
            info!("Received join: {:?}", room);
            let _ = socket.leave_all();
            let _ = socket.join(room.clone());
            let messages = store.get(&room).await;
            let _ = socket.emit("messages", Messages { messages });
        },
    );

    socket.on(
        "message",
        |socket: SocketRef, Data::<MessageIn>(data), store: State<state::MessageStore>| async move {
            info!("Received message: {:?}", data);

            let response = state::Message {
                text: data.text,
                user: format!("anon-{}", socket.id),
                date: chrono::Utc::now(),
            };

            store.insert(&data.room, response.clone()).await;

            let _ = socket.within(data.room).emit("message", response);
        },
    )
}

/// Fallback handler for incorrect URLs
///
/// Returns a 404 Not Found status code and a "Not Found" message.
///
/// # Returns
///
/// A tuple containing the status code and a static string slice
async fn fallback() -> (StatusCode, &'static str) {
    (StatusCode::NOT_FOUND, "Not Found")
}

/// Main function to set up and run the server
///
/// This function sets up the message store, socket.io layer, and the axum application
/// with routing and middleware. It starts the server using Shuttle runtime.
#[shuttle_runtime::main]
async fn main() -> shuttle_axum::ShuttleAxum {
    let messages = state::MessageStore::default();

    let (layer, io) = SocketIo::builder().with_state(messages).build_layer();

    io.ns("/", on_connect);

    let app = axum::Router::new()
        .fallback(fallback)
        .route("/", get(|| async { "Hello, World!" }))
        .with_state(io)
        .layer(
            ServiceBuilder::new()
                .layer(CorsLayer::permissive())
                .layer(layer),
        );

    info!("Starting server");

    Ok(app.into())
}
