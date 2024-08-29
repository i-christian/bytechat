use crate::messages::{MessageIn, Messages};
use crate::state;
use axum::http::StatusCode;
use socketioxide::extract::{Data, SocketRef, State};
use tracing::info;

/// Handles socket connection events
///
/// This function is called when a new socket connection is established. It sets up
/// listeners for "join", "message" and "typing" events on the socket.
///
/// # Arguments
///
/// * `socket` - A reference to the connected socket
pub async fn on_connect(socket: SocketRef) {
    info!("socket connected: {}", socket.id);

    socket.on(
        "join",
        |socket: SocketRef, Data::<String>(room), store: State<state::MessageStore>| async move {
            info!("Received join: {:?}", room);
            let _ = socket.leave_all();
            let _ = socket.join(room.clone());
            let messages = store.get(&room).await.unwrap_or_default();
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

            store.insert(&data.room, response.clone()).await.unwrap();

            let _ = socket.within(data.room).emit("message", response);
        },
    );

    socket.on(
        "typing",
        |socket: SocketRef, Data::<String>(room)| async move {
            info!("User {} is typing in room: {:?}", socket.id, room);

            let _ = socket
                .within(room)
                .except(socket.id.clone())
                .emit("typing", socket.id.clone());
        },
    );
}

/// Fallback handler for incorrect URLs
///
/// Returns a 404 Not Found status code and a "Not Found" message.
///
/// # Returns
///
/// A tuple containing the status code and a static string slice
pub async fn fallback() -> (StatusCode, &'static str) {
    (StatusCode::NOT_FOUND, "Not Found")
}
