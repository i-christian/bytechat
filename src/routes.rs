use crate::handlers::{fallback, on_connect};
use crate::state::MessageStore;
use axum::{routing::get, Router};
use socketioxide::SocketIo;
use sqlx::postgres::PgPool;
use tower::ServiceBuilder;
use tower_http::cors::CorsLayer;

/// Sets up and returns the Axum router with routes, middleware, and state.
///
/// # Arguments
///
/// * `pool` - A PostgreSQL connection pool.
///
/// # Returns
///
/// An Axum Router configured with routes, middleware, and state.
pub async fn setup_routes(pool: PgPool) -> Router {
    let messages = MessageStore::new(pool);

    let (layer, io) = SocketIo::builder().with_state(messages).build_layer();
    io.ns("/", on_connect);

    Router::new()
        .fallback(fallback)
        .route("/", get(|| async { "Hello, World!" }))
        .with_state(io)
        .layer(
            ServiceBuilder::new()
                .layer(CorsLayer::permissive())
                .layer(layer),
        )
}
