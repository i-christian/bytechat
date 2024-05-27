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

/// fallback - handler for wrong URL
/// *StatusCode: param
/// *str: param
/// return - a tuple of status code and a string slice
async fn fallback() -> (StatusCode, &'static str) {
    (StatusCode::NOT_FOUND, "Not Found")
}

async fn handler(axum::extract::State(io): axum::extract::State<SocketIo>) {
    info!("handler called");
    let _ = io.emit("hello", "world");
}

#[shuttle_runtime::main]
async fn main() -> shuttle_axum::ShuttleAxum {
    let messages = state::MessageStore::default();

    let (layer, io) = SocketIo::builder().with_state(messages).build_layer();

    io.ns("/", on_connect);

    let app = axum::Router::new()
        .fallback(fallback)
        .route("/", get(|| async { "Hello, World!" }))
        .route("/hello", get(handler))
        .with_state(io)
        .layer(
            ServiceBuilder::new()
                .layer(CorsLayer::permissive())
                .layer(layer),
        );

    info!("Starting server");

    Ok(app.into())
}
