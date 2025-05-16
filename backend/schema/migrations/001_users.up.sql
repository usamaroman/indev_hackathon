DO $$ BEGIN
IF NOT EXISTS (
    SELECT 1
    FROM pg_type
    WHERE typname = 'user_type'
) THEN CREATE TYPE user_type AS ENUM ('admin', 'customer');
END IF;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255),
    user_type user_type NOT NULL,
    hotel_id INT REFERENCES hotels(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO users (login, password, user_type) VALUES
('customer', 'customer', 'customer');

INSERT INTO users (login, password, user_type, hotel_id) VALUES
('admin', 'admin', 'admin', 1);
