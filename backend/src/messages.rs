use crate::AppState;
use axum::extract::ws::{Message, WebSocket};
use axum::{
    extract::{Path, State, WebSocketUpgrade},
    response::IntoResponse,
};

use futures::stream::StreamExt;
use sqlx::types::Uuid;
use std::collections::HashMap;
use tokio::sync::{mpsc, RwLock};

pub type UserSockets = RwLock<HashMap<Uuid, mpsc::UnboundedSender<Message>>>;

pub async fn websocket_handler(
    ws: WebSocketUpgrade,
    Path(room_id): Path<Uuid>,
    State(state): State<AppState>,
) -> impl IntoResponse {
    ws.on_upgrade(move |socket| handle_socket(socket, room_id, state))
}

async fn handle_socket(mut socket: WebSocket, room_id: Uuid, state: AppState) {
    let (sender, _receiver) = mpsc::unbounded_channel();
    let user_id = Uuid::new_v4();

    state.user_sockets.write().await.insert(user_id, sender);

    while let Some(Ok(msg)) = socket.next().await {
        if let Message::Text(text) = msg {
            let message_id = Uuid::new_v4();
            let query = sqlx::query(
                "INSERT INTO messages (message_id, user_id, room_id, text) VALUES ($1, $2, $3, $4)",
            )
            .bind(message_id)
            .bind(user_id)
            .bind(room_id)
            .bind(&text)
            .execute(&state.postgres);

            if let Err(e) = query.await {
                eprintln!("Failed to save message: {e}");
                continue;
            }

            for (&uid, tx) in state.user_sockets.read().await.iter() {
                if uid != user_id {
                    let _ = tx.send(Message::Text(text.clone()));
                }
            }
        }
    }

    state.user_sockets.write().await.remove(&user_id);
}
