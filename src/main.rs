use axum::{http::StatusCode, routing::get, Router};
use sqlx::PgPool;

async fn fallback() -> (StatusCode, &'static str) {
    (StatusCode::NOT_FOUND, "Not Found")
}

async fn hello_world() -> &'static str {
    "Hello, world!"
}

#[shuttle_runtime::main]
async fn main(#[shuttle_shared_db::Postgres] db: PgPool) -> shuttle_axum::ShuttleAxum {
    sqlx::migrate!()
        .run(&db)
        .await
        .expect("Looks like something went wrong with migrations: (");

    let router = Router::new()
        .route("/", get(hello_world))
        .fallback(fallback);

    Ok(router.into())
}
