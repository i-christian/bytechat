use chrono::Utc;
use serde::Serialize;
use sqlx::PgPool;

mod tests;

#[derive(Serialize, Clone, Debug)]
pub struct Message {
    pub text: String,
    pub user: String,
    pub date: chrono::DateTime<Utc>,
}

pub struct MessageStore {
    pub pool: PgPool,
}

impl MessageStore {
    /// Create a new MessageStore with a given PostgreSQL connection pool
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    /// Inserts a new message into the specified room
    ///
    /// # Arguments
    ///
    /// * `room` - The room to insert the message into
    /// * `message` - The message to be inserted
    pub async fn insert(&self, room: &str, message: Message) -> Result<(), sqlx::Error> {
        sqlx::query!(
            "INSERT INTO messages (room, text, username, created_at) VALUES ($1, $2, $3, $4)",
            room,
            message.text,
            message.user,
            message.date
        )
        .execute(&self.pool)
        .await?;
        Ok(())
    }

    /// Retrieves messages from the specified room
    ///
    /// # Arguments
    ///
    /// * `room` - The room to retrieve messages from
    ///
    /// # Returns
    ///
    /// A vector of messages in the room
    pub async fn get(&self, room: &str) -> Result<Vec<Message>, sqlx::Error> {
        let rows = sqlx::query!(
            "SELECT text, username, created_at FROM messages WHERE room = $1 ORDER BY created_at DESC LIMIT 20",
            room
        )
        .fetch_all(&self.pool)
        .await?;

        let messages = rows
            .into_iter()
            .map(|row| Message {
                text: row.text,
                user: row.username,
                date: row.created_at.unwrap_or_else(|| Utc::now()), // Provide a default value if `created_at` is `None`
            })
            .collect();

        Ok(messages)
    }
}
