use crate::AppState;
use axum::{
    extract::{Path, State},
    http::StatusCode,
    response::IntoResponse,
    Json,
};
use serde::{Deserialize, Serialize};
use sqlx::prelude::FromRow;
use uuid::Uuid;

#[derive(Deserialize, Serialize, FromRow)]
pub struct Room {
    pub name: String,
    pub description: Option<String>,
}

#[derive(Deserialize, Serialize)]
pub struct UpdateRoom {
    pub name: Option<String>,
    pub description: Option<String>,
}

#[derive(Deserialize)]
pub struct JoinRoomRequest {
    pub user_id: Uuid,
    pub room_id: Uuid,
}

pub async fn join_room(
    State(state): State<AppState>,
    Json(request): Json<JoinRoomRequest>,
) -> impl IntoResponse {
    let room_type: Option<String> =
        sqlx::query_scalar("SELECT room_type FROM rooms WHERE room_id = $1")
            .bind(request.room_id)
            .fetch_optional(&state.postgres)
            .await
            .expect("Failed to fetch room type");

    if room_type.as_deref() == Some("private") {
        let count = sqlx::query_scalar!(
            "SELECT COUNT(*) FROM chat_rooms WHERE room_id = $1",
            request.room_id
        )
        .fetch_one(&state.postgres)
        .await
        .unwrap_or_default();

        if count.unwrap_or(0) >= 2 {
            return (
                StatusCode::BAD_REQUEST,
                "Private room already has two users",
            )
                .into_response();
        }
    }

    let query = sqlx::query(
        "INSERT INTO chat_rooms (user_id, room_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
    )
    .bind(request.user_id)
    .bind(request.room_id)
    .execute(&state.postgres);

    match query.await {
        Ok(_) => (StatusCode::OK, "Joined room successfully").into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("Failed to join room: {}", e),
        )
            .into_response(),
    }
}

pub async fn create_room(
    State(state): State<AppState>,
    Json(room): Json<Room>,
) -> impl IntoResponse {
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
    let rooms_result = sqlx::query_as::<_, Room>("SELECT name, description FROM rooms")
        .fetch_all(&state.postgres)
        .await;

    match rooms_result {
        Ok(rooms) => Json(rooms).into_response(),
        Err(err) => {
            eprintln!("Database error: {}", err);
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                "Failed to retrieve rooms",
            )
                .into_response()
        }
    }
}

pub async fn delete_room(
    State(state): State<AppState>,
    Path(room_id): Path<Uuid>,
) -> impl IntoResponse {
    let query = sqlx::query("DELETE FROM rooms WHERE room_id = $1")
        .bind(room_id)
        .execute(&state.postgres);

    match query.await {
        Ok(result) if result.rows_affected() > 0 => {
            (StatusCode::OK, "Room deleted successfully!").into_response()
        }
        Ok(_) => (StatusCode::NOT_FOUND, "Room not found!").into_response(),
        Err(_) => (StatusCode::INTERNAL_SERVER_ERROR, "Failed to delete room").into_response(),
    }
}

pub async fn edit_room(
    State(state): State<AppState>,
    Path(room_id): Path<Uuid>,
    Json(update): Json<UpdateRoom>,
) -> impl IntoResponse {
    let query = sqlx::query(
        "UPDATE rooms SET name = COALESCE($1, name), description = COALESCE($2, description) WHERE room_id = $3"
    )
    .bind(&update.name)
    .bind(&update.description)
    .bind(room_id)
    .execute(&state.postgres);

    match query.await {
        Ok(result) if result.rows_affected() > 0 => {
            (StatusCode::OK, "Room updated successfully!").into_response()
        }
        Ok(_) => (StatusCode::NOT_FOUND, "Room not found!").into_response(),
        Err(_) => (StatusCode::INTERNAL_SERVER_ERROR, "Failed to update room").into_response(),
    }
}
