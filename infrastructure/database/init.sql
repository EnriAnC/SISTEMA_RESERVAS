-- Sistema de Reservas Database Schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(50) DEFAULT 'user' CHECK (role IN ('admin', 'user', 'manager')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- Resources table
CREATE TABLE resources (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(100) NOT NULL,
    location VARCHAR(255),
    capacity INTEGER DEFAULT 1,
    price_per_hour DECIMAL(10,2) DEFAULT 0.00,
    amenities JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Resource availability slots
CREATE TABLE resource_availability (
    id SERIAL PRIMARY KEY,
    resource_id INTEGER REFERENCES resources(id) ON DELETE CASCADE,
    day_of_week INTEGER CHECK (day_of_week >= 0 AND day_of_week <= 6), -- 0=Sunday, 6=Saturday
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(resource_id, day_of_week, start_time, end_time)
);

-- Bookings table
CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    resource_id INTEGER REFERENCES resources(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed')),
    total_price DECIMAL(10,2) DEFAULT 0.00,
    notes TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancellation_reason TEXT
);

-- Notifications table
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL,
    title VARCHAR(255),
    message TEXT NOT NULL,
    channel VARCHAR(50) DEFAULT 'email' CHECK (channel IN ('email', 'sms', 'push', 'webhook')),
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('high', 'normal', 'low')),
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed')),
    is_read BOOLEAN DEFAULT false,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP,
    read_at TIMESTAMP
);

-- Audit log table for tracking changes
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    record_id INTEGER NOT NULL,
    action VARCHAR(20) NOT NULL CHECK (action IN ('INSERT', 'UPDATE', 'DELETE')),
    old_values JSONB,
    new_values JSONB,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_uuid ON users(uuid);
CREATE INDEX idx_users_active ON users(is_active);

CREATE INDEX idx_resources_type ON resources(type);
CREATE INDEX idx_resources_active ON resources(is_active);
CREATE INDEX idx_resources_uuid ON resources(uuid);

CREATE INDEX idx_resource_availability_resource_id ON resource_availability(resource_id);
CREATE INDEX idx_resource_availability_day ON resource_availability(day_of_week);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_resource_id ON bookings(resource_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_start_time ON bookings(start_time);
CREATE INDEX idx_bookings_end_time ON bookings(end_time);
CREATE INDEX idx_bookings_uuid ON bookings(uuid);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

CREATE INDEX idx_audit_logs_table_record ON audit_logs(table_name, record_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Functions for automatic timestamp updates
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for automatic timestamp updates
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_resources_updated_at BEFORE UPDATE ON resources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bookings_updated_at BEFORE UPDATE ON bookings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to prevent booking conflicts
CREATE OR REPLACE FUNCTION check_booking_conflict()
RETURNS TRIGGER AS $$
DECLARE
    conflict_count INTEGER;
BEGIN
    -- Check for overlapping bookings
    SELECT COUNT(*)
    INTO conflict_count
    FROM bookings
    WHERE resource_id = NEW.resource_id
      AND status IN ('pending', 'confirmed')
      AND id != COALESCE(NEW.id, 0)
      AND (
          (NEW.start_time >= start_time AND NEW.start_time < end_time) OR
          (NEW.end_time > start_time AND NEW.end_time <= end_time) OR
          (NEW.start_time <= start_time AND NEW.end_time >= end_time)
      );
    
    IF conflict_count > 0 THEN
        RAISE EXCEPTION 'Booking conflict detected for resource % between % and %', 
            NEW.resource_id, NEW.start_time, NEW.end_time;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to prevent booking conflicts
CREATE TRIGGER check_booking_conflict_trigger
    BEFORE INSERT OR UPDATE ON bookings
    FOR EACH ROW EXECUTE FUNCTION check_booking_conflict();

-- Audit trigger function
CREATE OR REPLACE FUNCTION audit_trigger()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        INSERT INTO audit_logs (table_name, record_id, action, old_values)
        VALUES (TG_TABLE_NAME, OLD.id, TG_OP, row_to_json(OLD));
        RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO audit_logs (table_name, record_id, action, old_values, new_values)
        VALUES (TG_TABLE_NAME, NEW.id, TG_OP, row_to_json(OLD), row_to_json(NEW));
        RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        INSERT INTO audit_logs (table_name, record_id, action, new_values)
        VALUES (TG_TABLE_NAME, NEW.id, TG_OP, row_to_json(NEW));
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create audit triggers for all main tables
CREATE TRIGGER audit_users_trigger
    AFTER INSERT OR UPDATE OR DELETE ON users
    FOR EACH ROW EXECUTE FUNCTION audit_trigger();

CREATE TRIGGER audit_resources_trigger
    AFTER INSERT OR UPDATE OR DELETE ON resources
    FOR EACH ROW EXECUTE FUNCTION audit_trigger();

CREATE TRIGGER audit_bookings_trigger
    AFTER INSERT OR UPDATE OR DELETE ON bookings
    FOR EACH ROW EXECUTE FUNCTION audit_trigger();

CREATE TRIGGER audit_notifications_trigger
    AFTER INSERT OR UPDATE OR DELETE ON notifications
    FOR EACH ROW EXECUTE FUNCTION audit_trigger();

-- Insert sample data
INSERT INTO users (email, password_hash, first_name, last_name, phone, role) VALUES
('admin@example.com', '$2y$10$example_hash', 'Admin', 'User', '+1234567890', 'admin'),
('john.doe@example.com', '$2y$10$example_hash', 'John', 'Doe', '+1234567891', 'user'),
('jane.smith@example.com', '$2y$10$example_hash', 'Jane', 'Smith', '+1234567892', 'manager');

INSERT INTO resources (name, description, type, location, capacity, price_per_hour, amenities) VALUES
('Conference Room A', 'Large conference room with projector', 'meeting_room', 'Floor 1', 10, 50.00, '{"projector": true, "whiteboard": true, "video_conference": true}'),
('Conference Room B', 'Small meeting room', 'meeting_room', 'Floor 2', 6, 30.00, '{"whiteboard": true, "phone": true}'),
('Auditorium', 'Main auditorium for events', 'auditorium', 'Ground Floor', 100, 200.00, '{"projector": true, "sound_system": true, "stage": true}'),
('Hot Desk 1', 'Flexible workspace', 'desk', 'Floor 3', 1, 15.00, '{"monitor": true, "ethernet": true}');

INSERT INTO resource_availability (resource_id, day_of_week, start_time, end_time) VALUES
-- Conference Room A (Monday to Friday, 8 AM to 6 PM)
(1, 1, '08:00', '18:00'),
(1, 2, '08:00', '18:00'),
(1, 3, '08:00', '18:00'),
(1, 4, '08:00', '18:00'),
(1, 5, '08:00', '18:00'),
-- Conference Room B (Monday to Friday, 9 AM to 5 PM)
(2, 1, '09:00', '17:00'),
(2, 2, '09:00', '17:00'),
(2, 3, '09:00', '17:00'),
(2, 4, '09:00', '17:00'),
(2, 5, '09:00', '17:00');

-- Sample bookings
INSERT INTO bookings (user_id, resource_id, start_time, end_time, status, total_price, notes) VALUES
(2, 1, '2024-06-10 10:00:00', '2024-06-10 12:00:00', 'confirmed', 100.00, 'Team meeting'),
(3, 2, '2024-06-11 14:00:00', '2024-06-11 15:00:00', 'pending', 30.00, 'Client call');
