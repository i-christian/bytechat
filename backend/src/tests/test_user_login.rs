use super::*;

use http::StatusCode;
use serde_json::json;

#[tokio::test]
async fn it_should_login_successfully() {
    let server = new_test_app().await;

    let response = server
        .post("/api/auth/login")
        .json(&json!({
                "email": "boruto@email.com",
                "password": "changethis"
        }))
        .await;

    let result = response.status_code();
    assert_eq!(result, StatusCode::OK);
}

#[tokio::test]
#[should_panic]
async fn it_should_fail_to_login_user_with_wrong_password() {
    let server = new_test_app().await;

    let response = server
        .post("/api/auth/login")
        .json(&json!({
                "email": "boruto@email.com",
                "password": "wrongpassword"
        }))
        .await;

    let result = response.status_code();
    assert_eq!(result, StatusCode::BAD_REQUEST);
}

#[tokio::test]
#[should_panic]
async fn it_should_not_login_without_email() {
    let server = new_test_app().await;

    let response = server
        .post("/api/auth/login")
        .json(&json!({
                "password": "wrongpassword"
        }))
        .await;

    let result = response.status_code();
    assert_eq!(result, StatusCode::BAD_REQUEST);
}
