use crate::{router::create_api_router, AppState};
use axum::Router;
use axum_extra::extract::cookie::Key;
use axum_test::TestServer;
use dotenv::dotenv;
use sqlx::postgres::PgPoolOptions;
use std::net::{IpAddr, Ipv4Addr};

#[tokio::test]
async fn check_database_connectivity() {
    dotenv().ok();
    let durl = std::env::var("DATABASE_URL_TEST").expect("set DATABASE_URL_TEST");

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&durl)
        .await
        .expect("unable to make connection");

    assert!(!pool.is_closed());
}

async fn new_test_app() -> TestServer {
    dotenv().ok();
    let durl = std::env::var("DATABASE_URL_TEST").expect("set DATABASE_URL_TEST");
    let domain = std::env::var("DOMAIN_URL").expect("failed to load domain");

    let postgres = PgPoolOptions::new()
        .max_connections(5)
        .connect(&durl)
        .await
        .expect("unable to make connection");

    sqlx::migrate!()
        .run(&postgres)
        .await
        .expect("Failed to run migrations");

    let state = AppState {
        postgres,
        domain,
        key: Key::generate(),
    };

    let api_router = create_api_router(state);
    let app = Router::new().nest("/api", api_router);
    let localhost = IpAddr::V4(Ipv4Addr::new(0, 0, 0, 0));

    TestServer::builder()
        .http_transport_with_ip_port(Some(localhost), Some(8000))
        .save_cookies()
        .expect_success_by_default()
        .build(app)
        .expect("Failed to create a connection")
}

mod test_app_health;
mod test_register_user;
mod test_user_login;
