use axum::{routing::get, Router};
use sqlx::{postgres::PgPoolOptions, PgPool};
use dotenv::dotenv;
use std::env;
use tracing::info;


#[derive(Clone)]
struct MyState {
    pool: PgPool,
}

#[tokio::main]
async fn main() {
    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");

    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(&database_url)
        .await
        .expect("Failed to create a database connection.");

    sqlx::migrate!().run(&pool).await.expect("Migration failed");
    let state = MyState{ pool };

    let app = Router::new()
        .route("/hello", get(|| async {"Hello, World!"}))
        .with_state(state);

    info!("Starting server");
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.expect("Failed to bind port");
    axum::serve(listener, app).await.expect("Failed to start application");
}
