mod handlers;
mod messages;
mod routes;
mod state;

use routes::setup_routes;
use sqlx::postgres::PgPoolOptions;
use tracing::info;

#[shuttle_runtime::main]
async fn main() -> shuttle_axum::ShuttleAxum {
    let database_url = std::env::var("DATABASE_URL").expect("DATABASE_URL must be set");

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await
        .expect("Failed to create pool.");

    //sqlx::migrate!().run(&pool).await.expect("Migration failed");

    let app = setup_routes(pool).await;

    info!("Starting server");

    Ok(app.into())
}
