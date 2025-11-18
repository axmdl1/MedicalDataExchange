-- patient_access_requests - запросы пациентов на просмотр своих данных
CREATE TABLE IF NOT EXISTS patient_access_requests (
    id BIGSERIAL PRIMARY KEY,
    patient_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    clinic_id BIGINT NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    medical_data_id BIGINT NOT NULL REFERENCES medical_data(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, approved, rejected, expired
    blockchain_tx_id VARCHAR(255), -- ID транзакции в блокчейне
    requested_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    approved_at TIMESTAMP WITHOUT TIME ZONE,
    expires_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_patient_access_requests_patient_id ON patient_access_requests(patient_id);
CREATE INDEX IF NOT EXISTS idx_patient_access_requests_clinic_id ON patient_access_requests(clinic_id);
CREATE INDEX IF NOT EXISTS idx_patient_access_requests_status ON patient_access_requests(status);
CREATE INDEX IF NOT EXISTS idx_patient_access_requests_expires_at ON patient_access_requests(expires_at);

-- temporary_patient_data - временные копии данных для пациентов (удаляются через 15 минут)
CREATE TABLE IF NOT EXISTS temporary_patient_data (
    id BIGSERIAL PRIMARY KEY,
    access_token VARCHAR(255) UNIQUE NOT NULL,
    patient_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    clinic_id BIGINT NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    medical_data_id BIGINT NOT NULL REFERENCES medical_data(id) ON DELETE CASCADE,
    encrypted_data TEXT NOT NULL, -- Зашифрованная копия данных
    granted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_temporary_patient_data_access_token ON temporary_patient_data(access_token);
CREATE INDEX IF NOT EXISTS idx_temporary_patient_data_patient_id ON temporary_patient_data(patient_id);
CREATE INDEX IF NOT EXISTS idx_temporary_patient_data_expires_at ON temporary_patient_data(expires_at);
CREATE INDEX IF NOT EXISTS idx_temporary_patient_data_is_revoked ON temporary_patient_data(is_revoked);

-- Функция для автоматического удаления истекших временных данных
CREATE OR REPLACE FUNCTION delete_expired_temporary_data()
RETURNS void AS $$
BEGIN
    DELETE FROM temporary_patient_data
    WHERE expires_at < NOW() OR is_revoked = TRUE;

    UPDATE patient_access_requests
    SET status = 'expired'
    WHERE expires_at < NOW() AND status = 'approved';
END;
$$ LANGUAGE plpgsql;

-- Можно настроить pg_cron для периодического запуска, но для упрощения
-- будем вызывать из приложения
