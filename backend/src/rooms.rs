use crate::{auth::get_user_id, AppState};
use axum::{
    extract::{Path, State},
    http::StatusCode,
    response::IntoResponse,
    Json,
};

use axum_extra::extract::cookie::PrivateCookieJar;

use serde::{Deserialize, Serialize};
use sqlx::prelude::FromRow;
use uuid::Uuid;

#[derive(Deserialize, Serialize, FromRow)]
pub struct Room {
    pub name: String,
    pub description: Option<String>,
}

#[derive(Deserialize, Serialize)]
pub struct CreateRoomRequest {
    pub room: Room,
}

#[derive(Deserialize, Serialize)]
pub struct UpdateRoom {
    pub name: Option<String>,
    pub description: Option<String>,
}

#[derive(Deserialize)]
pub struct JoinRoomRequest {
    pub room_id: Uuid,
}

#[derive(Deserialize, Serialize)]
pub struct CreatePrivateRoomRequest {
    pub room: Room,
    pub user1_id: Uuid,
    pub user2_id: Uuid,
}

pub async fn create_public_room(
    State(state): State<AppState>,
    Json(request): Json<CreateRoomRequest>,
) -> impl IntoResponse {
    let room_id = Uuid::new_v4();
    let room_type = "public".to_string();

    let query = sqlx::query(
        "INSERT INTO rooms (room_id, name, description, room_type) VALUES ($1, $2, $3, $4)",
    )
    .bind(room_id)
    .bind(&request.room.name)
    .bind(&request.room.description)
    .bind(&room_type)
    .execute(&state.postgres)
    .await;

    match query {
        Ok(_) => (StatusCode::CREATED, Json(room_id)).into_response(),
        Err(_) => (StatusCode::BAD_REQUEST, "Room creation failed!").into_response(),
    }
}

#[axum::debug_handler]
pub async fn create_private_room(
    State(state): State<AppState>,
    Json(request): Json<CreatePrivateRoomRequest>,
) -> impl IntoResponse {
    let room_id = Uuid::new_v4();
    let room_type = "private".to_string();

    let query = sqlx::query(
        "INSERT INTO rooms (room_id, name, description, room_type) VALUES ($1, $2, $3, $4)",
    )
    .bind(room_id)
    .bind(&request.room.name)
    .bind(&request.room.description)
    .bind(&room_type)
    .execute(&state.postgres)
    .await;

    match query {
        Ok(_) => {
            sqlx::query("INSERT INTO chat_rooms (user_id, room_id) VALUES ($1, $2)")
                .bind(request.user1_id)
                .bind(room_id)
                .execute(&state.postgres)
                .await
                .ok();

            sqlx::query("INSERT INTO chat_rooms (user_id, room_id) VALUES ($1, $2)")
                .bind(request.user2_id)
                .bind(room_id)
                .execute(&state.postgres)
                .await
                .ok();

            (StatusCode::CREATED, Json(room_id)).into_response()
        }
        Err(_) => (StatusCode::BAD_REQUEST, "Room creation failed!").into_response(),
    }
}

pub async fn join_room(
    State(state): State<AppState>,
    jar: PrivateCookieJar,
    Json(request): Json<JoinRoomRequest>,
) -> impl IntoResponse {
    let Some(user_id) = get_user_id(jar, State(state.clone())).await else {
        return (StatusCode::FORBIDDEN, "Unauthorized").into_response();
    };

    let room_type: Option<String> =
        sqlx::query_scalar("SELECT room_type FROM rooms WHERE room_id = $1")
            .bind(request.room_id)
            .fetch_optional(&state.postgres)
            .await
            .expect("Failed to fetch room type");

    if room_type.as_deref() != Some("public") {
        return (StatusCode::FORBIDDEN, "Cannot join a private room").into_response();
    }

    let query = sqlx::query(
        "INSERT INTO chat_rooms (user_id, room_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
    )
    .bind(user_id)
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

pub async fn list_rooms(State(state): State<AppState>) -> impl IntoResponse {
    let rooms_result = sqlx::query_as::<_, Room>("SELECT room_id, name, description FROM rooms")
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
