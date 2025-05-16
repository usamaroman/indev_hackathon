CREATE TABLE IF NOT EXISTS hotels (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO hotels (name) VALUES ('Hackathon hotel');

CREATE TABLE IF NOT EXISTS room_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,      
    capacity INT NOT NULL 
);

INSERT INTO room_types (name, capacity) VALUES 
('На 1 человека', 1), ('На 2 человека', 2), ('На 3 человека', 3), ('На 4 человека', 4);

CREATE TABLE IF NOT EXISTS rooms (
    room_number TEXT PRIMARY KEY,  
    hotel_id INT REFERENCES hotels(id),
    room_type_id INT REFERENCES room_types(id),
    floor INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS reservations (
    id SERIAL PRIMARY KEY,
    room_id TEXT REFERENCES rooms(room_number),
    user_id INT REFERENCES users(id),
    check_in DATE NOT NULL,
    check_out DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'confirmed', -- 'checked_in'
    created_at TIMESTAMP DEFAULT NOW() 
);

CREATE INDEX IF NOT EXISTS idx_reservations_room_dates ON reservations(room_id, check_in, check_out);

INSERT INTO rooms (hotel_id, room_type_id, floor, room_number) VALUES
(1, 1, 1, '105'),
(1, 1, 1, '106'),
(1, 2, 1, '104'),
(1, 2, 1, '108'),
(1, 3, 1, '107'),
(1, 3, 1, '103'),
(1, 4, 1, '101'),
(1, 4, 1, '102');
