use std::collections::{HashMap, VecDeque};
use tokio::sync::RwLock;

#[derive(serde::Serialize, Clone, Debug)]
pub struct Message {
    pub text: String,
    pub user: String,
    pub date: chrono::DateTime<chrono::Utc>,
}

pub type RoomStore = HashMap<String, VecDeque<Message>>;

#[derive(Default)]
pub struct MessageStore {
    pub messages: RwLock<RoomStore>,
}

impl MessageStore {
    /// Inserts a new message into the specified room
    ///
    /// This function writes the new message to the room's message list, ensuring
    /// the list does not exceed 20 messages by removing the oldest messages if necessary.
    ///
    /// # Arguments
    ///
    /// * `room` - The room to insert the message into
    /// * `message` - The message to be inserted
    pub async fn insert(&self, room: &String, message: Message) {
        let mut binding = self.messages.write().await;
        let messages = binding.entry(room.clone()).or_default();
        messages.push_front(message);
        messages.truncate(20);
    }

    /// Retrieves messages from the specified room
    ///
    /// This function returns a list of messages for the specified room in
    /// reverse chronological order (newest first).
    ///
    /// # Arguments
    ///
    /// * `room` - The room to retrieve messages from
    ///
    /// # Returns
    ///
    /// A vector of messages in the room
    pub async fn get(&self, room: &String) -> Vec<Message> {
        let messages = self.messages.read().await.get(room).cloned();
        messages.unwrap_or_default().into_iter().rev().collect()
    }
}
