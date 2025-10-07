-- Users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) DEFAULT 'user', -- 'user' atau 'admin'
    created_at TIMESTAMP DEFAULT NOW()
);

-- Schedules
CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    movie_title VARCHAR(255) NOT NULL,
    cinema_branch VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    show_time TIMESTAMP NOT NULL,
    total_seats INT NOT NULL,
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'cancelled'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Seats (1 schedule = N seats)
CREATE TABLE seats (
    id SERIAL PRIMARY KEY,
    schedule_id INT NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    seat_number VARCHAR(10) NOT NULL, -- e.g., "A1", "B5"
    status VARCHAR(20) DEFAULT 'available', -- 'available', 'locked', 'sold'
    locked_until TIMESTAMP,
    sold_at TIMESTAMP,
    UNIQUE(schedule_id, seat_number)
);

-- Refunds (optional for audit)
CREATE TABLE refunds (
    id SERIAL PRIMARY KEY,
    schedule_id INT NOT NULL,
    seat_id INT NOT NULL,
    reason VARCHAR(100), -- 'user_refund', 'cinema_cancel'
    refunded_at TIMESTAMP DEFAULT NOW()
);

-- Admin user // email : admin@example.com , password : passwordadmin
INSERT INTO users (email, password_hash, role) VALUES
('admin@example.com', '$2y$10$rWtnD2zZkuIw/EbrEG1ma.TyHCHjcWn6LgIIMr0qclwbNyiFRJBJi', 'admin');

-- Regular user // email : user@example.com , password : passworduser
INSERT INTO users (email, password_hash, role) VALUES
('user@example.com', '$2y$10$5m08xsUKsP535ALYosXL1e4.sN7m.6rDSw5Sj2YyVb34EOXMtjS.q', 'user');