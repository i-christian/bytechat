use crate::auth::get_user_id;
use crate::AppState;
use axum::extract::ws::{Message, WebSocket};
use axum::{
    extract::{Path, State, WebSocketUpgrade},
    response::IntoResponse,
};

use axum_extra::extract::PrivateCookieJar;
use futures::stream::StreamExt;
use http::StatusCode;
use sqlx::types::Uuid;
use std::collections::HashMap;
use tokio::sync::{mpsc, RwLock};

pub type UserSockets = RwLock<HashMap<Uuid, mpsc::UnboundedSender<Message>>>;

pub async fn websocket_handler(
    ws: WebSocketUpgrade,
    Path(room_id): Path<Uuid>,
    State(state): State<AppState>,
    jar: PrivateCookieJar,
) -> impl IntoResponse {
    let user_id = match get_user_id(jar, State(state.clone())).await {
        Some(id) => id,
        None => {
            return (StatusCode::FORBIDDEN, "Unauthorized").into_response();
        }
    };

    ws.on_upgrade(move |socket| handle_socket(socket, room_id, user_id, state))
}

async fn handle_socket(mut socket: WebSocket, room_id: Uuid, user_id: Uuid, state: AppState) {
    let (sender, _receiver) = mpsc::unbounded_channel();

    state
        .user_sockets
        .write()
        .await
        .insert(user_id, sender.clone());

    let messages: Vec<(Uuid, Uuid, String, time::PrimitiveDateTime)> = sqlx::query_as(
        r#"
        SELECT message_id, user_id, text, created_at 
        FROM messages 
        WHERE room_id = $1 
        ORDER BY created_at ASC
        "#,
    )
    .bind(room_id)
    .fetch_all(&state.postgres)
    .await
    .unwrap_or_default();

    for (_, user_id, text, created_at) in messages {
        let msg = format!("{}: {} (at: {})", user_id, text, created_at);
        if let Err(e) = socket.send(Message::Text(msg)).await {
            eprintln!("Failed to send message: {:?}", e);
        }
    }

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

            let formatted_message = format!("{}: {}", user_id, text);
            for (&uid, tx) in state.user_sockets.read().await.iter() {
                if uid != user_id {
                    let _ = tx.send(Message::Text(formatted_message.clone()));
                }
            }
        }
    }

    state.user_sockets.write().await.remove(&user_id);
}
