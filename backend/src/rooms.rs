use axum::{
    extract::State,
    http::StatusCode,
    response::IntoResponse,
    Json,
};
use serde::{Deserialize, Serialize};
use crate::AppState;

#[derive(Deserialize, Serialize)]
pub struct Room {
    pub name: String,
    pub description: Option<String>,
}

pub async fn create_room(State(state): State<AppState>, Json(room): Json<Room>) -> impl IntoResponse {
    let query = sqlx::query("INSERT INTO rooms (name, description) VALUES ($1, $2)")
        .bind(&room.name)
        .bind(&room.description)
        .execute(&state.postgres);
    match query.await {
        Ok(_) => (StatusCode::CREATED, "Room created successfully!").into_response(),
        Err(_) => (StatusCode::BAD_REQUEST, "Room creation failed!").into_response(),
    }
}

pub async fn list_rooms(State(state): State<AppState>) -> impl IntoResponse {
    let rooms = sqlx::query_as!(Room, "SELECT name, description FROM rooms")
        .fetch_all(&state.postgres)
        .await;
    match rooms {
        Ok(rooms) => Json(rooms).into_response(),
        Err(_) => (StatusCode::INTERNAL_SERVER_ERROR, "Failed to retrieve rooms").into_response(),
    }
}
