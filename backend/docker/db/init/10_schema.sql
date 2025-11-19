-- clinics
CREATE TABLE IF NOT EXISTS clinics (
                                       id BIGSERIAL PRIMARY KEY,
                                       name VARCHAR(255) NOT NULL,
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
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

-- medical_data - comprehensive medical records
CREATE TABLE IF NOT EXISTS medical_data (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    clinic_id BIGINT NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,

    -- Visit information
    visit_date TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    visit_type VARCHAR(50), -- emergency, routine, follow-up, consultation
    department VARCHAR(100), -- cardiology, neurology, surgery, etc.
    attending_doctor VARCHAR(255),

    -- Patient complaints and symptoms
    chief_complaint TEXT,
    symptoms TEXT, -- JSON array or comma-separated
    symptom_duration VARCHAR(100),
    pain_level INTEGER CHECK (pain_level >= 0 AND pain_level <= 10),

    -- Vital signs
    temperature DECIMAL(4,1), -- Celsius
    blood_pressure_systolic INTEGER,
    blood_pressure_diastolic INTEGER,
    heart_rate INTEGER,
    respiratory_rate INTEGER,
    oxygen_saturation INTEGER,
    weight DECIMAL(5,2), -- kg
    height DECIMAL(5,2), -- cm
    bmi DECIMAL(4,2),

    -- Medical assessment
    diagnosis TEXT NOT NULL,
    diagnosis_code VARCHAR(20), -- ICD-10 code
    secondary_diagnoses TEXT, -- JSON array of additional diagnoses
    severity VARCHAR(20), -- mild, moderate, severe, critical

    -- Medical history
    medical_history TEXT,
    surgical_history TEXT,
    family_history TEXT,
    allergies TEXT, -- medications, food, environmental
    current_medications TEXT, -- JSON array of medications

    -- Treatment and care plan
    treatment_plan TEXT NOT NULL,
    prescribed_medications TEXT, -- JSON array with dosage, frequency
    procedures_performed TEXT,
    lab_tests_ordered TEXT,
    imaging_ordered TEXT,
    referrals TEXT,

    -- Lab results
    lab_results JSONB, -- structured lab data
    lab_results_summary TEXT,

    -- Doctor's notes
    doctor_notes TEXT,
    follow_up_instructions TEXT,
    follow_up_date TIMESTAMP WITHOUT TIME ZONE,
    restrictions TEXT, -- physical activity restrictions, dietary restrictions

    -- Administrative
    insurance_info TEXT,
    billing_code VARCHAR(50),
    estimated_cost DECIMAL(10,2),

    -- Metadata
    record_status VARCHAR(20) DEFAULT 'active', -- active, archived, amended
    confidentiality_level VARCHAR(20) DEFAULT 'normal', -- normal, restricted, highly_restricted
    data_hash VARCHAR(64), -- SHA-256 hash for blockchain verification

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_medical_data_user_id ON medical_data(user_id);
CREATE INDEX IF NOT EXISTS idx_medical_data_clinic_id ON medical_data(clinic_id);
CREATE INDEX IF NOT EXISTS idx_medical_data_visit_date ON medical_data(visit_date);
CREATE INDEX IF NOT EXISTS idx_medical_data_diagnosis_code ON medical_data(diagnosis_code);
CREATE INDEX IF NOT EXISTS idx_medical_data_department ON medical_data(department);

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