use axum::{routing::get, Router};
use axum_test::TestServer;

async fn get_hello() -> &'static str {
    "Hello, World!"
}

#[tokio::test]
async fn test_root_route() {
    // Build an application with a route.
    let app = Router::new().route("/", get(get_hello));

    // Run the application for testing.
    let server = TestServer::new(app).unwrap();

    // Get the request for the root route.
    let response = server.get("/").await;

    // Assert the response text is "Hello, World!".
    assert_eq!(response.text(), "Hello, World!");
}
