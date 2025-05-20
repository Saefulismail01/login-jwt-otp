-- DDL untuk table OTPs
CREATE TABLE IF NOT EXISTS otps (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_email_code UNIQUE (email, code)
);

-- Index untuk mempercepat pencarian
CREATE INDEX IF NOT EXISTS idx_otps_email ON otps(email);
CREATE INDEX IF NOT EXISTS idx_otps_code ON otps(code);
CREATE INDEX IF NOT EXISTS idx_otps_expires_at ON otps(expires_at);

-- Trigger untuk memperbarui updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_otps_updated_at
    BEFORE UPDATE ON otps
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger untuk membersihkan OTP yang expired
CREATE OR REPLACE FUNCTION cleanup_expired_otps()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM otps WHERE expires_at < CURRENT_TIMESTAMP;
    RETURN NULL;
END;
$$ language 'plpgsql';

-- Trigger yang akan dijalankan setiap 10 menit
CREATE EVENT TRIGGER cleanup_expired_otps_trigger
    ON ddl_command_end
    WHEN TAG IN ('CREATE TABLE', 'ALTER TABLE')
    EXECUTE FUNCTION cleanup_expired_otps();
