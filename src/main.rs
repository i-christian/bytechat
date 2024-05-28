mod handlers;
mod messages;
mod state;

use axum::routing::get;
use handlers::{fallback, on_connect};
use socketioxide::SocketIo;
use sqlx::postgres::PgPoolOptions;
use tower::ServiceBuilder;
use tower_http::cors::CorsLayer;
use tracing::info;

/// Main function to set up and run the server
///
/// This function sets up the message store, socket.io layer, and the axum application
/// with routing and middleware. It starts the server using Shuttle runtime.
#[shuttle_runtime::main]
async fn main() -> shuttle_axum::ShuttleAxum {
    let database_url = std::env::var("DATABASE_URL").expect("DATABASE_URL must be set");

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await
        .expect("Failed to create pool.");

    //sqlx::migrate!().run(&pool).await.expect("Migration failed");

    let messages = state::MessageStore::new(pool);

    let (layer, io) = SocketIo::builder().with_state(messages).build_layer();

    io.ns("/", on_connect);

    let app = axum::Router::new()
        .fallback(fallback)
        .route("/", get(|| async { "Hello, World!" }))
        .with_state(io)
        .layer(
            ServiceBuilder::new()
                .layer(CorsLayer::permissive())
                .layer(layer),
        );

    info!("Starting server");

    Ok(app.into())
}
