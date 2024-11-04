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
