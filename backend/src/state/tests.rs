#[cfg(test)]
mod tests {

    use crate::state::Message;
    use crate::state::MessageStore;

    use chrono::Utc;
    use sqlx::{Executor, PgPool};

    // Utility function to set up the test database connection pool
    async fn setup_test_db() -> PgPool {
        let database_url = std::env::var("DATABASE_URL").expect("DATABASE_URL must be set");
        let pool = PgPool::connect(&database_url)
            .await
            .expect("Failed to connect to database");

        // Create the messages table for testing
        pool.execute(
            "CREATE TABLE IF NOT EXISTS messages (
                id SERIAL PRIMARY KEY,
                room TEXT NOT NULL,
                text TEXT NOT NULL,
                username TEXT NOT NULL,
                created_at TIMESTAMPTZ NOT NULL
            )",
        )
        .await
        .expect("Failed to create messages table");

        // Clear the messages table before each test
        pool.execute("TRUNCATE TABLE messages RESTART IDENTITY CASCADE")
            .await
            .expect("Failed to truncate messages table");

        pool
    }

    #[tokio::test]
    async fn test_insert_message() {
        let pool = setup_test_db().await;
        let store = MessageStore::new(pool);

        let message = Message {
            text: String::from("Hello, world!"),
            user: String::from("test_user"),
            date: Utc::now(),
        };

        let result = store.insert("test_room", message).await;
        assert!(result.is_ok());
    }
}
