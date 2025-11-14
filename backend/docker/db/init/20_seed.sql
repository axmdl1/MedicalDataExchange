-- Клиники
INSERT INTO clinics (id, name, address, phone, email)
VALUES
    (1, 'City Clinic', 'Main Street 12, Tashkent', '+998901112233', 'info@cityclinic.uz'),
    (2, 'Regional Hospital', 'Hospital Ave 45, Tashkent', '+998902223344', 'contact@hospital.uz'),
    (3, 'Medical Center Plus', 'Health Street 78, Tashkent', '+998903334455', 'info@medcenter.uz'),
    (4, 'Family Health Center', 'Park Road 23, Samarkand', '+998904445566', 'family@health.uz'),
    (5, 'Children Hospital', 'Kids Street 15, Bukhara', '+998905556677', 'children@hospital.uz'),
    (6, 'Cardiology Center', 'Heart Avenue 89, Tashkent', '+998906667788', 'cardio@center.uz'),
    (7, 'Dental Clinic ProSmile', 'Smile Street 34, Tashkent', '+998907778899', 'prosmile@dental.uz')
ON CONFLICT (id) DO NOTHING;

-- Reset sequence to avoid duplicate key errors
SELECT setval('clinics_id_seq', (SELECT MAX(id) FROM clinics));

-- Админ (пароль: admin123)
DO $$
DECLARE
    _hash TEXT := crypt('admin123', gen_salt('bf'));
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@system.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Admin', 'Root', 'admin@system.com', '+998900000000', 'admin', _hash, NULL);
    END IF;
END $$;

-- Сотрудники (пароль: doctor123)
DO $$
DECLARE
    _hash TEXT := crypt('doctor123', gen_salt('bf'));
BEGIN
    -- City Clinic (3 сотрудника)
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'doctor1@clinic.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Иван', 'Петров', 'doctor1@clinic.com', '+998911111111', 'employee', _hash, 1);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'nurse1@clinic.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Ольга', 'Смирнова', 'nurse1@clinic.com', '+998911111112', 'employee', _hash, 1);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin1@clinic.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Сергей', 'Волков', 'admin1@clinic.com', '+998911111113', 'employee', _hash, 1);
    END IF;

    -- Regional Hospital (4 сотрудника)
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'doctor2@hospital.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Мария', 'Иванова', 'doctor2@hospital.com', '+998922222222', 'employee', _hash, 2);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'doctor2_2@hospital.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Петр', 'Кузнецов', 'doctor2_2@hospital.com', '+998922222223', 'employee', _hash, 2);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'nurse2@hospital.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Елена', 'Павлова', 'nurse2@hospital.com', '+998922222224', 'employee', _hash, 2);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'surgeon@hospital.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Андрей', 'Соколов', 'surgeon@hospital.com', '+998922222225', 'employee', _hash, 2);
    END IF;

    -- Medical Center Plus (2 сотрудника)
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'doctor3@medcenter.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Алексей', 'Сидоров', 'doctor3@medcenter.com', '+998933333333', 'employee', _hash, 3);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'therapist@medcenter.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Наталья', 'Федорова', 'therapist@medcenter.com', '+998933333334', 'employee', _hash, 3);
    END IF;

    -- Family Health Center (3 сотрудника)
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'doctor4@family.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Виктор', 'Морозов', 'doctor4@family.com', '+998944444441', 'employee', _hash, 4);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'pediatr@family.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Светлана', 'Зайцева', 'pediatr@family.com', '+998944444442', 'employee', _hash, 4);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'nurse4@family.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Татьяна', 'Белова', 'nurse4@family.com', '+998944444443', 'employee', _hash, 4);
    END IF;

    -- Children Hospital (2 сотрудника)
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'pediatr@children.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Анастасия', 'Лебедева', 'pediatr@children.com', '+998955555551', 'employee', _hash, 5);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'neonatolog@children.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Дмитрий', 'Орлов', 'neonatolog@children.com', '+998955555552', 'employee', _hash, 5);
    END IF;

    -- Cardiology Center (2 сотрудника)
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'cardiolog@cardio.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Владимир', 'Попов', 'cardiolog@cardio.com', '+998966666661', 'employee', _hash, 6);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'surgeon_cardio@cardio.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Игорь', 'Васильев', 'surgeon_cardio@cardio.com', '+998966666662', 'employee', _hash, 6);
    END IF;

    -- Dental Clinic (2 сотрудника)
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'dentist@dental.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Роман', 'Николаев', 'dentist@dental.com', '+998977777771', 'employee', _hash, 7);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'orthodont@dental.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Юлия', 'Семенова', 'orthodont@dental.com', '+998977777772', 'employee', _hash, 7);
    END IF;
END $$;

-- Пациенты (пароль: patient123)
DO $$
DECLARE
    _hash TEXT := crypt('patient123', gen_salt('bf'));
    _patient1_id BIGINT;
    _patient2_id BIGINT;
    _patient3_id BIGINT;
BEGIN
    -- Пациент 1
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient1@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Анна', 'Смирнова', 'patient1@example.com', '+998944444444', 'patient', _hash, NULL)
        RETURNING id INTO _patient1_id;
    ELSE
        SELECT id INTO _patient1_id FROM users WHERE email = 'patient1@example.com';
    END IF;

    -- Пациент 2
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient2@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Дмитрий', 'Козлов', 'patient2@example.com', '+998955555555', 'patient', _hash, NULL)
        RETURNING id INTO _patient2_id;
    ELSE
        SELECT id INTO _patient2_id FROM users WHERE email = 'patient2@example.com';
    END IF;

    -- Пациент 3
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient3@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Елена', 'Новикова', 'patient3@example.com', '+998966666666', 'patient', _hash, NULL)
        RETURNING id INTO _patient3_id;
    ELSE
        SELECT id INTO _patient3_id FROM users WHERE email = 'patient3@example.com';
    END IF;

    -- Медицинские записи для пациента 1 (City Clinic)
    IF NOT EXISTS (SELECT 1 FROM medical_data WHERE user_id = _patient1_id AND clinic_id = 1) THEN
        INSERT INTO medical_data (user_id, clinic_id, diagnosis, complaint, treatment, medications, allergies, doctor_notes, lab_results, created_at)
        VALUES
            (_patient1_id, 1, 'ОРВИ', 'Высокая температура 38.5, слабость, головная боль', 'Постельный режим, обильное питье', 'Парацетамол 500мг 3 раза в день', 'Нет', 'Повторный осмотр через 5 дней', 'Общий анализ крови: в норме', NOW() - INTERVAL '10 days'),
            (_patient1_id, 1, 'Профилактический осмотр', 'Плановый осмотр', 'Не требуется', 'Витамин D', 'Нет', 'Здорова, рекомендован активный образ жизни', 'Анализы в норме', NOW() - INTERVAL '30 days');
    END IF;

    -- Медицинские записи для пациента 2 (Regional Hospital)
    IF NOT EXISTS (SELECT 1 FROM medical_data WHERE user_id = _patient2_id AND clinic_id = 2) THEN
        INSERT INTO medical_data (user_id, clinic_id, diagnosis, complaint, treatment, medications, allergies, doctor_notes, lab_results, created_at)
        VALUES
            (_patient2_id, 2, 'Гастрит', 'Боли в желудке, изжога', 'Диета, исключить острое и жареное', 'Омепразол 20мг утром натощак', 'Аспирин', 'Контроль через 2 недели, возможна ФГДС', 'Биохимия крови: в норме', NOW() - INTERVAL '5 days'),
            (_patient2_id, 2, 'Артериальная гипертензия', 'Повышенное давление 150/95', 'Контроль давления, снижение соли', 'Эналаприл 10мг утром', 'Нет', 'Вести дневник давления', 'ЭКГ: без патологии', NOW() - INTERVAL '15 days');
    END IF;

    -- Медицинские записи для пациента 3 (Medical Center Plus)
    IF NOT EXISTS (SELECT 1 FROM medical_data WHERE user_id = _patient3_id AND clinic_id = 3) THEN
        INSERT INTO medical_data (user_id, clinic_id, diagnosis, complaint, treatment, medications, allergies, doctor_notes, lab_results, created_at)
        VALUES
            (_patient3_id, 3, 'Аллергический ринит', 'Заложенность носа, чихание, слезотечение', 'Избегать контакта с аллергенами', 'Цетиризин 10мг на ночь', 'Пыльца березы', 'Рекомендована консультация аллерголога', 'Общий IgE повышен', NOW() - INTERVAL '3 days');
    END IF;

    -- Запросы на передачу данных (Data Transfers)
    -- Пациент 1 переезжает из City Clinic в Regional Hospital
    IF NOT EXISTS (SELECT 1 FROM data_transfer WHERE user_id = _patient1_id) THEN
        INSERT INTO data_transfer (medical_data_id, from_clinic_id, to_clinic_id, user_id, status, created_at)
        SELECT id, 1, 2, _patient1_id, 'pending', NOW() - INTERVAL '2 days'
        FROM medical_data
        WHERE user_id = _patient1_id AND clinic_id = 1
        LIMIT 1;
    END IF;

    -- Пациент 2 хочет получить мнение врача из Medical Center Plus
    IF NOT EXISTS (SELECT 1 FROM data_transfer WHERE user_id = _patient2_id AND status = 'confirmed') THEN
        INSERT INTO data_transfer (medical_data_id, from_clinic_id, to_clinic_id, user_id, status, created_at)
        SELECT id, 2, 3, _patient2_id, 'confirmed', NOW() - INTERVAL '7 days'
        FROM medical_data
        WHERE user_id = _patient2_id AND clinic_id = 2 AND diagnosis = 'Гастрит'
        LIMIT 1;
    END IF;

END $$;