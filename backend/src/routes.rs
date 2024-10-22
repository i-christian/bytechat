use axum::{routing::get, Router};
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
    Router::new()
        .route("/", get(|| async { "Hello, World!" }));
}
