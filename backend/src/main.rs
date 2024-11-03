use std::env;
use dotenv::dotenv;

use axum::extract::FromRef;
use axum::Router;
use axum_extra::extract::cookie::Key;
use sqlx::{postgres::PgPoolOptions, PgPool};
use tower_http::services::{ServeDir, ServeFile};

mod auth;
mod router;

use router::create_api_router;
use tracing::info;
use tracing_subscriber::FmtSubscriber;

#[derive(Clone)]
pub struct AppState {
    pub postgres: PgPool,
    pub domain: String,
    pub key: Key,
}

impl FromRef<AppState> for Key {
    fn from_ref(state: &AppState) -> Self {
        state.key.clone()
    }
}


#[tokio::main]
async fn main() {
    tracing::subscriber::set_global_default(FmtSubscriber::default())
        .expect("setting default subscriber failed");

    let (database_url, domain) = grab_secrets();
    let postgres = PgPoolOptions::new()
        .max_connections(10)
        .connect(&database_url)
        .await
        .expect("Failed to create a database connection.");
    
    sqlx::migrate!()
        .run(&postgres)
        .await
        .expect("Failed to run migrations");
    
    info!("Database connected successfully");

    let state = AppState {
        postgres,
        domain,
        key: Key::generate(),
    };

    let api_router = create_api_router(state.clone());

    let router = Router::new().nest("/api", api_router).nest_service(
        "/",
        ServeDir::new("dist").not_found_service(ServeFile::new("dist/index.html")),
    );

    info!("Started Application on: http://{}:3000", state.domain);
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.expect("Failed to bind port");
    axum::serve(listener, router).await.expect("Failed to start application");
}

fn grab_secrets() -> (String, String) {
    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");

    let domain: String = match env::var("DOMAIN_URL") {
        Ok(var) => var,
        Err(..) => "None".to_string()
        
    };

    (
        database_url,
        domain
    )
}
