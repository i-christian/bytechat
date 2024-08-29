use crate::state;
#[derive(Debug, serde::Deserialize)]
pub struct MessageIn {
    pub room: String,
    pub text: String,
}

#[derive(serde::Serialize)]
pub struct Messages {
    pub messages: Vec<state::Message>,
}
