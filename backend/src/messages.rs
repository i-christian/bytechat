use crate::auth::get_user_id;
use crate::AppState;
use axum::extract::ws::{Message, WebSocket};
use axum::{
    extract::{Path, State, WebSocketUpgrade},
    response::IntoResponse,
};
use axum_extra::extract::PrivateCookieJar;
use futures::stream::StreamExt;
use futures::SinkExt;
use http::StatusCode;
use sqlx::types::Uuid;
use std::collections::HashMap;
use tokio::sync::{mpsc, RwLock};
use tokio::time::Duration;
use tracing::{error, info, warn};

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

    {
        state
            .user_sockets
            .write()
            .await
            .insert(user_id, sender.clone());
        state.user_status.write().await.insert(user_id, true);
    }

    let online_message = Message::Text(format!("User {} is online", user_id));
    broadcast_status_change(&state, user_id, online_message).await;

    if let Err(e) = send_previous_messages(&state, &mut socket, room_id).await {
        error!("Failed to send previous messages: {:?}", e);
        return;
    }

    let (mut socket_tx, mut socket_rx) = socket.split();
    let mut ping_interval = tokio::time::interval(Duration::from_secs(30));

    let read_handle = tokio::spawn({
        let state = state.clone();
        async move {
            while let Some(Ok(msg)) = socket_rx.next().await {
                match msg {
                    Message::Text(text) => {
                        if let Err(e) = handle_incoming_text(&state, user_id, room_id, &text).await
                        {
                            error!("Failed to handle message: {:?}", e);
                        }
                    }
                    Message::Close(_) => break,
                    _ => {}
                }
            }
        }
    });

    let write_handle = tokio::spawn(async move {
        loop {
            ping_interval.tick().await;
            if socket_tx.send(Message::Ping(vec![])).await.is_err() {
                warn!("Connection lost with user {}", user_id);
                break;
            }
        }
    });

    let _ = tokio::try_join!(read_handle, write_handle);

    {
        state.user_sockets.write().await.remove(&user_id);
        state.user_status.write().await.remove(&user_id);
    }

    let offline_message = Message::Text(format!("User {} is offline", user_id));
    broadcast_status_change(&state, user_id, offline_message).await;

    info!("User {} disconnected from room {}", user_id, room_id);
}

async fn send_previous_messages(
    state: &AppState,
    socket: &mut WebSocket,
    room_id: Uuid,
) -> Result<(), sqlx::Error> {
    let messages: Vec<(Uuid, Uuid, String, time::PrimitiveDateTime)> = sqlx::query_as(
        "SELECT message_id, user_id, text, created_at FROM messages WHERE room_id = $1 ORDER BY created_at ASC",
    )
    .bind(room_id)
    .fetch_all(&state.postgres)
    .await?;

    for (_, user_id, text, created_at) in messages {
        let msg = format!("{}: {} (at: {})", user_id, text, created_at);
        if let Err(e) = socket.send(Message::Text(msg)).await {
            error!("Failed to send message to user {}: {:?}", user_id, e);
        }
    }

    Ok(())
}

async fn handle_incoming_text(
    state: &AppState,
    user_id: Uuid,
    room_id: Uuid,
    text: &str,
) -> Result<(), sqlx::Error> {
    let message_id = Uuid::new_v4();

    sqlx::query(
        "INSERT INTO messages (message_id, user_id, room_id, text) VALUES ($1, $2, $3, $4)",
    )
    .bind(message_id)
    .bind(user_id)
    .bind(room_id)
    .bind(text)
    .execute(&state.postgres)
    .await?;

    let formatted_message = format!("{}: {}", user_id, text);
    broadcast_to_other_users(state, user_id, Message::Text(formatted_message)).await;

    Ok(())
}

async fn broadcast_status_change(state: &AppState, user_id: Uuid, message: Message) {
    let sockets = state.user_sockets.read().await;
    for (&uid, tx) in sockets.iter() {
        if uid != user_id {
            let _ = tx.send(message.clone());
        }
    }
}

async fn broadcast_to_other_users(state: &AppState, sender_id: Uuid, message: Message) {
    let sockets = state.user_sockets.read().await;
    for (&uid, tx) in sockets.iter() {
        if uid != sender_id {
            let _ = tx.send(message.clone());
        }
    }
}
