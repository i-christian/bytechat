use super::*;

#[tokio::test]
async fn it_should_respond() {
    let server = new_test_app().await;
    let response = server.get("/api/health").await;
    response.assert_text("Hello world!");
}
