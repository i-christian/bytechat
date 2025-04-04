-- +goose Up

-- ROLES TABLE
CREATE TABLE roles (
    role_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

-- USERS TABLE
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    last_name VARCHAR(50) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    gender CHAR(1) NOT NULL CHECK (gender IN ('M', 'F')),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'offline' CHECK (status IN ('online', 'offline')),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    role_id UUID NOT NULL,
    CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE RESTRICT
);

-- SESSIONS TABLE
CREATE TABLE IF NOT EXISTS sessions (
    session_id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    expires TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP + INTERVAL '2 week',
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);


-- Create Rooms Table
CREATE TABLE IF NOT EXISTS rooms (
    room_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    room_type VARCHAR(20) NOT NULL CHECK (room_type IN ('public', 'private'))
);

-- Create ChatRooms Linking Table (for Many-to-Many relationship between Users and Rooms)
CREATE TABLE IF NOT EXISTS chat_rooms (
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,   
    room_id UUID REFERENCES rooms(room_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, room_id)
);

-- Create Messages Table
CREATE TABLE IF NOT EXISTS messages (
    message_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    room_id UUID REFERENCES rooms(room_id) ON DELETE CASCADE,
    text TEXT,
    multimedia BYTEA,                                     
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for Optimization
CREATE INDEX IF NOT EXISTS idx_user_room ON chat_rooms(user_id, room_id);
CREATE INDEX IF NOT EXISTS idx_room_created_at ON messages(room_id, created_at);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);

-- +goose Down
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS chat_rooms;;
DROP TABLE IF EXISTS messages;
