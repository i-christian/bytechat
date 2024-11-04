use super::*;

use http::StatusCode;
use serde_json::json;

#[tokio::test]
async fn it_should_create_user() {
    let server = new_test_app().await;

    let response = server
        .post("/api/auth/register")
        .json(&json!({
                "name": "Boruto",
                "email": "boruto@email.com",
                "password": "changethis"
        }))
        .await;

    let result = response.status_code();
    assert_eq!(result, StatusCode::CREATED);
}

#[tokio::test]
#[should_panic]
async fn it_should_fail_to_create_user() {
    let server = new_test_app().await;

    let response = server
        .post("/api/auth/register")
        .json(&json!({
                "name": "Boruto",
                "email": "boruto@email.com",
                "password": "changethis"
        }))
        .await;

    let result = response.status_code();
    assert_eq!(result, StatusCode::BAD_REQUEST);
}
