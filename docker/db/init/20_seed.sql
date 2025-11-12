-- базовая клиника
INSERT INTO clinics (name, address, phone_number)
VALUES ('City Clinic', 'Main Street 12', '+998901112233')
    ON CONFLICT DO NOTHING;

-- админ (пароль admin123)
DO $$
DECLARE
_hash TEXT := crypt('admin123', gen_salt('bf'));
BEGIN
  IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@system.com') THEN
    INSERT INTO users (first_name, last_name, email, phone_number, type, password)
    VALUES ('Admin', 'Root', 'admin@system.com', '+998900000000', 'admin', _hash);
END IF;
END $$;