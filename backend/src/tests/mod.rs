use std::net::{IpAddr, Ipv4Addr};
use axum_extra::extract::cookie::Key;
use axum_test::TestServer;
use axum::Router;
use sqlx::postgres::PgPoolOptions;
use dotenv::dotenv;
use crate::{router::create_api_router, AppState};

#[tokio::test]
async fn check_database_connectivity() {
    dotenv().ok();
    let durl = std::env::var("DATABASE_URL").expect("set DATABASE_URL");

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&durl)
        .await
        .expect("unable to make connection");

    assert_eq!(pool.is_closed(), false);
}


async fn new_test_app() -> TestServer {
    dotenv().ok();
    let durl = std::env::var("DATABASE_URL").expect("set DATABASE_URL");
    let domain = std::env::var("DOMAIN_URL").expect("failed to load domain");

    let postgres = PgPoolOptions::new()
        .max_connections(5)
        .connect(&durl)
        .await
        .expect("unable to make connection");

    let state = AppState {
        postgres,
        domain,
        key: Key::generate(),
    };

    let api_router = create_api_router(state);
    let app = Router::new().nest("/api", api_router);
    let localhost = IpAddr::V4(Ipv4Addr::new(0,0,0,0));

    TestServer::builder()
        .http_transport_with_ip_port(Some(localhost), Some(8000))
        .save_cookies()
        .expect_success_by_default()
        .build(app)
        .unwrap()

}

#[cfg(test)]
mod test_app_health {
    use super::*;

    #[tokio::test]
    async fn it_should_respond() {
        let server = new_test_app().await;
        let response = server.get(&"/api/health").await;
        response.assert_text("Hello world!");
    }
}

// #[cfg(test)]
// mod test_register_user {
//     use super::*;

//     use axum::Json;

//     #[tokio::test]
//     async fn it_should_create_user() {
//         let server = new_test_app().await;

//         let response = server
//             .post(&"/login")        
//     }
// }
