mod handlers;
mod messages;
mod routes;
mod state;

use dotenv::dotenv;
use routes::setup_routes;
use sqlx::postgres::PgPoolOptions;
use std::env;
use tracing::info;


#[tokio::main]
async fn main() {
    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await
        .expect("Failed to create pool.");

    sqlx::migrate!().run(&pool).await.expect("Migration failed");

    let app = setup_routes(pool).await;

    info!("Starting server");

    Ok(app.into())
}
