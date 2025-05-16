DO $$ BEGIN
IF NOT EXISTS (
    SELECT 1
    FROM pg_type
    WHERE typname = 'user_type'
) THEN CREATE TYPE user_type AS ENUM ('admin', 'customer');
END IF;
END $$;

DO $$ BEGIN
IF NOT EXISTS (
    SELECT 1
    FROM pg_type
    WHERE typname = 'reservation_type'
) THEN CREATE TYPE reservation_type AS ENUM ('confirmed', 'checked_in', 'checked_out');
END IF;
END $$;
