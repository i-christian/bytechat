use crate::AppState;
use axum::{
    http, middleware,
    routing::{get, post, put},
    Router,
};
use http::header::{ACCEPT, AUTHORIZATION, ORIGIN};
use http::HeaderValue;
use http::Method;
use tower_http::cors::CorsLayer;

use crate::auth::{
    delete_user, edit_user, get_all_users, login, logout, register, validate_session,
};
use crate::rooms::{
    create_private_room, create_public_room, delete_room, edit_room, join_room, list_rooms,
};

pub fn create_api_router(state: AppState) -> Router {
    let cors = CorsLayer::new()
        .allow_credentials(true)
        .allow_methods(vec![Method::GET, Method::POST, Method::PUT, Method::DELETE])
        .allow_headers(vec![ORIGIN, AUTHORIZATION, ACCEPT])
        .allow_origin(state.domain.parse::<HeaderValue>().unwrap());

    let rooms_router = Router::new()
        .route("/", get(list_rooms))
        .route("/", post(join_room))
        .route("/create/private", post(create_private_room))
        .route("/create/public", post(create_public_room))
        .route("/:room_id", put(edit_room).delete(delete_room));

    let auth_router = Router::new()
        .route("/register", post(register))
        .route("/login", post(login))
        .route("/logout", get(logout))
        .route("/users", put(edit_user).delete(delete_user))
        .route("/get_all_users", get(get_all_users));

    Router::new()
        // nest protected routes here
        .nest("/rooms", rooms_router)
        .layer(middleware::from_fn_with_state(
            state.clone(),
            validate_session,
        ))
        .nest("/auth", auth_router)
        .route("/health", get(hello_world))
        .with_state(state)
        .layer(cors)
}

pub async fn hello_world() -> &'static str {
    "Hello world!"
}
