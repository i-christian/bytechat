use axum::{
    extract::{Request, State},
    http::StatusCode,
    middleware::Next,
    response::{IntoResponse, Response},
    Json,
};
use axum_extra::extract::cookie::{Cookie, PrivateCookieJar, SameSite};
use serde::{Deserialize, Serialize};
use sqlx::prelude::FromRow;
use sqlx::Row;
use time::Duration;
use uuid::Uuid;

use crate::AppState;

#[derive(Deserialize)]
pub struct RegisterDetails {
    name: String,
    email: String,
    password: String,
}

#[derive(Deserialize)]
pub struct LoginDetails {
    email: String,
    password: String,
}

#[derive(Deserialize)]
pub struct UpdateUserDetails {
    pub name: Option<String>,
    pub password: Option<String>,
}

#[derive(Serialize, Deserialize, FromRow)]
pub struct UserInfo {
    pub id: Uuid,
    pub name: String,
    pub email: String,
}

pub async fn register(
    State(state): State<AppState>,
    Json(newuser): Json<RegisterDetails>,
) -> impl IntoResponse {
    let hashed_password = bcrypt::hash(newuser.password, 10).unwrap();

    let query = sqlx::query("INSERT INTO users (name, email, password) values ($1, $2, $3)")
        .bind(newuser.name)
        .bind(newuser.email)
        .bind(hashed_password)
        .execute(&state.postgres);

    match query.await {
        Ok(_) => (StatusCode::CREATED, "Account created!".to_string()).into_response(),
        Err(e) => (
            StatusCode::BAD_REQUEST,
            format!("Something went wrong: {e}"),
        )
            .into_response(),
    }
}

pub async fn login(
    State(state): State<AppState>,
    jar: PrivateCookieJar,
    Json(login): Json<LoginDetails>,
) -> Result<(PrivateCookieJar, StatusCode), StatusCode> {
    let query = sqlx::query("SELECT * FROM users WHERE email = $1")
        .bind(&login.email)
        .fetch_one(&state.postgres);

    match query.await {
        Ok(res) => {
            match bcrypt::verify(login.password, res.get("password")) {
                Ok(true) => {}
                Ok(false) => return Err(StatusCode::BAD_REQUEST),
                Err(_) => return Err(StatusCode::BAD_REQUEST),
            }

            let session_id = Uuid::new_v4();

            sqlx::query("INSERT INTO sessions (session_id, user_id) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET session_id = EXCLUDED.session_id")
                .bind(session_id)
                .bind(res.get::<Uuid, _>("id"))
                .execute(&state.postgres)
                .await
                .expect("Couldn't insert session :(");

            let cookie = Cookie::build(("sessionid", session_id.to_string()))
                .secure(!cfg!(debug_assertions)) // Only send the cookie over HTTPS in production
                .same_site(SameSite::Strict)
                .http_only(true)
                .path("/")
                .max_age(Duration::WEEK)
                .build();

            Ok((jar.add(cookie), StatusCode::OK))
        }

        Err(_) => Err(StatusCode::BAD_REQUEST),
    }
}

pub async fn get_all_users(State(state): State<AppState>) -> impl IntoResponse {
    let query = sqlx::query_as::<_, UserInfo>("SELECT id, name, email FROM users")
        .fetch_all(&state.postgres)
        .await;

    match query {
        Ok(users) => Json(users).into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("Failed to retrieve users: {}", e),
        )
            .into_response(),
    }
}

pub async fn get_user_id(jar: PrivateCookieJar, State(state): State<AppState>) -> Option<Uuid> {
    if let Some(cookie) = jar.get("sessionid") {
        if let Ok(session_id) = Uuid::parse_str(cookie.value()) {
            let query = sqlx::query("SELECT user_id FROM sessions WHERE session_id = $1")
                .bind(session_id)
                .fetch_one(&state.postgres)
                .await;

            match query {
                Ok(record) => {
                    let user_id = record.get("user_id");
                    Some(user_id)
                }
                Err(_) => None,
            }
        } else {
            None
        }
    } else {
        None
    }
}

pub async fn logout(
    State(state): State<AppState>,
    jar: PrivateCookieJar,
) -> Result<PrivateCookieJar, StatusCode> {
    let Some(cookie) = jar.get("sessionid").map(|cookie| cookie.value().to_owned()) else {
        return Ok(jar);
    };

    let query = sqlx::query("DELETE FROM sessions WHERE session_id = $1")
        .bind(cookie.parse::<Uuid>().unwrap())
        .execute(&state.postgres);

    match query.await {
        Ok(_) => Ok(jar.remove(Cookie::from("sessionid"))),
        Err(_) => Err(StatusCode::INTERNAL_SERVER_ERROR),
    }
}

pub async fn validate_session(
    jar: PrivateCookieJar,
    State(state): State<AppState>,
    request: Request,
    next: Next,
) -> (PrivateCookieJar, Response) {
    let Some(cookie) = jar.get("sessionid").map(|cookie| cookie.value().to_owned()) else {
        println!("Couldn't find a cookie in the jar");
        return (
            jar,
            (StatusCode::FORBIDDEN, "Forbidden!".to_string()).into_response(),
        );
    };

    let find_session =
        sqlx::query("SELECT * FROM sessions WHERE session_id = $1 AND expires > CURRENT_TIMESTAMP")
            .bind(cookie.parse::<Uuid>().unwrap())
            .execute(&state.postgres)
            .await;

    match find_session {
        Ok(_) => (jar, next.run(request).await),
        Err(_) => (
            jar,
            (StatusCode::FORBIDDEN, "Forbidden!".to_string()).into_response(),
        ),
    }
}

pub async fn edit_user(
    State(state): State<AppState>,
    jar: PrivateCookieJar,
    Json(details): Json<UpdateUserDetails>,
) -> impl IntoResponse {
    let Some(user_id) = get_user_id(jar, State(state.clone())).await else {
        return (StatusCode::FORBIDDEN, "Unauthorized").into_response();
    };

    let hashed_password = details
        .password
        .as_ref()
        .map(|password| bcrypt::hash(password, 10).unwrap());

    let query = sqlx::query(
        "UPDATE users 
         SET name = COALESCE($1, name), password = COALESCE($2, password)
         WHERE id = $3",
    )
    .bind(&details.name)
    .bind(&hashed_password)
    .bind(user_id)
    .execute(&state.postgres);

    match query.await {
        Ok(result) if result.rows_affected() > 0 => {
            (StatusCode::OK, "User updated successfully!").into_response()
        }
        Ok(_) => (StatusCode::NOT_FOUND, "User not found!").into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("Failed to update user: {e}"),
        )
            .into_response(),
    }
}

pub async fn delete_user(
    State(state): State<AppState>,
    jar: PrivateCookieJar,
) -> impl IntoResponse {
    let Some(user_id) = get_user_id(jar, State(state.clone())).await else {
        return (StatusCode::FORBIDDEN, "Unauthorized").into_response();
    };

    let query = sqlx::query("DELETE FROM users WHERE id = $1")
        .bind(user_id)
        .execute(&state.postgres);

    match query.await {
        Ok(result) if result.rows_affected() > 0 => {
            (StatusCode::OK, "User deleted successfully!").into_response()
        }
        Ok(_) => (StatusCode::NOT_FOUND, "User not found!").into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("Failed to delete user: {e}"),
        )
            .into_response(),
    }
}
