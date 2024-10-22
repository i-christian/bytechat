-- Create Users Table
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  
    username VARCHAR(50) UNIQUE NOT NULL,                
    email VARCHAR(100) UNIQUE NOT NULL,                  
    password TEXT NOT NULL,                              
    admin_id UUID REFERENCES users(user_id),             
    status VARCHAR(20) DEFAULT 'active',                 
    description TEXT,                                    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create Rooms Table
CREATE TABLE IF NOT EXISTS rooms (
    room_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  
    name VARCHAR(100) NOT NULL,                         
    description TEXT,                                    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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
    delivered_at TIMESTAMP WITH TIME ZONE,                
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for optimization
CREATE INDEX IF NOT EXISTS idx_user_room ON chat_rooms(user_id, room_id);
CREATE INDEX IF NOT EXISTS idx_room_created_at ON messages(room_id, created_at);
