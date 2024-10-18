use chrono::Utc;
use serde::Serialize;
use sqlx::PgPool;
use std::collections::HashMap;
use std::sync::{Arc, Mutex};


/// Represents a message with text, user, and creation date.
#[derive(Serialize, Clone, Debug)]
pub struct Message {
    pub text: String,
    pub user: String,
    pub date: chrono::DateTime<Utc>,
}

/// Represents a store for managing messages.
pub struct MessageStore {
    pub pool: PgPool,
    cache: Arc<Mutex<HashMap<String, Vec<Message>>>>,
}

impl MessageStore {
    /// Creates a new MessageStore with a given PostgreSQL connection pool.
    ///
    /// # Arguments
    ///
    /// * `pool` - The PostgreSQL connection pool.
    ///
    /// # Returns
    ///
    /// A new instance of MessageStore.
    pub fn new(pool: PgPool) -> Self {
        Self {
            pool,
            cache: Arc::new(Mutex::new(HashMap::new())),
        }
    }

    /// Inserts a new message into the specified room.
    ///
    /// This method inserts a new message into the database and invalidates the cache
    /// for the specified room to ensure that subsequent calls to `get` retrieve fresh data.
    ///
    /// # Arguments
    ///
    /// * `room` - The room to insert the message into.
    /// * `message` - The message to be inserted.
    ///
    /// # Returns
    ///
    /// A Result indicating success or an SQLx error.
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

        // Remove cached messages for this room after inserting a new message
        self.cache.lock().unwrap().remove(room);

        Ok(())
    }

    /// Retrieves messages from the specified room, either from cache or database.
    ///
    /// This method first checks if messages for the specified room are available in the cache.
    /// If found, it returns the cached messages. If not found, it retrieves messages from the
    /// database, caches them, and then returns them.
    ///
    /// # Arguments
    ///
    /// * `room` - The room to retrieve messages from.
    ///
    /// # Returns
    ///
    /// A Result containing a vector of messages or an SQLx error.
    pub async fn get(&self, room: &str) -> Result<Vec<Message>, sqlx::Error> {
        // Check if messages for this room are already in the cache
        if let Some(messages) = self.cache.lock().unwrap().get(room) {
            return Ok(messages.clone()); // Return cached messages
        }

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
                date: row.created_at.unwrap_or_else(|| Utc::now()),
            })
            .collect::<Vec<Message>>();

        // Cache the retrieved messages for this room
        self.cache
            .lock()
            .unwrap()
            .insert(room.to_string(), messages.clone());

        Ok(messages)
    }
}
