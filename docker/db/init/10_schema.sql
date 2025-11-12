-- clinics
CREATE TABLE IF NOT EXISTS clinics (
                                       id BIGSERIAL PRIMARY KEY,
                                       name VARCHAR(255) NOT NULL,
    address TEXT,
    phone_number VARCHAR(20),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
    );

-- users (soft delete поддерживается GORM через deleted_at)
CREATE TABLE IF NOT EXISTS users (
                                     id BIGSERIAL PRIMARY KEY,
                                     clinic_id BIGINT REFERENCES clinics(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(20),
    type VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITHOUT TIME ZONE
    );

-- частичный уникальный индекс на email только для активных
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_indexes WHERE indexname = 'users_email_unique_idx'
  ) THEN
CREATE UNIQUE INDEX users_email_unique_idx
    ON users (lower(email))
    WHERE deleted_at IS NULL;
END IF;
END$$;

-- medical_data (минимум — чтобы не мешало)
CREATE TABLE IF NOT EXISTS medical_data (
                                            id BIGSERIAL PRIMARY KEY,
                                            user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    clinic_id BIGINT NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    diagnosis TEXT,
    complaint TEXT,
    treatment TEXT,
    medications TEXT,
    allergies TEXT,
    doctor_notes TEXT,
    lab_results TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
    );

-- data_transfer
CREATE TABLE IF NOT EXISTS data_transfer (
                                             id BIGSERIAL PRIMARY KEY,
                                             medical_data_id BIGINT NOT NULL REFERENCES medical_data(id) ON DELETE CASCADE,
    from_clinic_id BIGINT NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    to_clinic_id BIGINT NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
    );