-- Real medical data with comprehensive patient records
DO $$
DECLARE
    _patient1_id BIGINT;
    _patient2_id BIGINT;
    _patient3_id BIGINT;
BEGIN
    -- Get patient IDs
    SELECT id INTO _patient1_id FROM users WHERE email = 'patient1@example.com';
    SELECT id INTO _patient2_id FROM users WHERE email = 'patient2@example.com';
    SELECT id INTO _patient3_id FROM users WHERE email = 'patient3@example.com';

    -- Delete old simple medical data
    DELETE FROM medical_data WHERE user_id IN (_patient1_id, _patient2_id, _patient3_id);

    -- Patient 1: ОРВИ (Acute Respiratory Viral Infection) at City Clinic
    INSERT INTO medical_data (
        user_id, clinic_id, visit_date, visit_type, department, attending_doctor,
        chief_complaint, symptoms, symptom_duration, pain_level,
        temperature, blood_pressure_systolic, blood_pressure_diastolic,
        heart_rate, respiratory_rate, oxygen_saturation, weight, height, bmi,
        diagnosis, diagnosis_code, severity,
        medical_history, allergies, current_medications,
        treatment_plan, prescribed_medications, lab_tests_ordered,
        lab_results, lab_results_summary,
        doctor_notes, follow_up_instructions, follow_up_date, restrictions,
        record_status, confidentiality_level
    ) VALUES (
        _patient1_id, 1,
        NOW() - INTERVAL '3 days', 'emergency', 'General Medicine', 'Dr. Ivan Petrov',
        'High fever and severe headache', 'Fever, headache, body aches, fatigue, dry cough', '2 days', 6,
        38.7, 125, 78, 92, 18, 98, 68.5, 165.0, 25.16,
        'Acute upper respiratory tract infection (URTI)', 'J06.9', 'moderate',
        'No chronic diseases. Vaccinated against COVID-19 (2 doses).', 'No known drug allergies', 'None',
        'Rest at home, increase fluid intake (2-3 liters daily), symptomatic treatment. Return if fever persists >3 days or breathing difficulties occur.',
        '[{"name": "Paracetamol", "dosage": "500mg", "frequency": "Every 6 hours as needed", "duration": "5 days"},
          {"name": "Vitamin C", "dosage": "1000mg", "frequency": "Once daily", "duration": "7 days"}]',
        'Complete blood count (CBC)',
        '{"wbc": 8.5, "rbc": 4.5, "hemoglobin": 13.2, "platelets": 245, "neutrophils": 65, "lymphocytes": 28}',
        'WBC: 8.5 x10^9/L (normal), Hemoglobin: 13.2 g/dL (normal), Platelets: 245 x10^9/L (normal)',
        'Patient presents with classic viral URTI symptoms. Vital signs stable. Lungs clear on auscultation. No signs of bacterial complications. Advised symptomatic management and isolation at home.',
        'Monitor temperature twice daily. Return immediately if: fever >39°C for >3 days, difficulty breathing, chest pain, or worsening symptoms. Wear mask around others. Rest for at least 5 days.',
        NOW() + INTERVAL '5 days',
        'Avoid strenuous physical activity. Stay home from work/school for 5 days minimum. Avoid contact with elderly or immunocompromised individuals.',
        'active', 'normal'
    );

    -- Patient 1: Preventive checkup at City Clinic (1 month ago)
    INSERT INTO medical_data (
        user_id, clinic_id, visit_date, visit_type, department, attending_doctor,
        chief_complaint, pain_level,
        temperature, blood_pressure_systolic, blood_pressure_diastolic,
        heart_rate, oxygen_saturation, weight, height, bmi,
        diagnosis, diagnosis_code, severity,
        medical_history, allergies, current_medications,
        treatment_plan, prescribed_medications,
        lab_results, lab_results_summary,
        doctor_notes, follow_up_instructions, follow_up_date,
        record_status, confidentiality_level
    ) VALUES (
        _patient1_id, 1,
        NOW() - INTERVAL '35 days', 'routine', 'General Medicine', 'Dr. Olga Smirnova',
        'Annual preventive health checkup', 0,
        36.6, 118, 76, 72, 99, 67.0, 165.0, 24.61,
        'Healthy individual - routine checkup', 'Z00.0', 'mild',
        'No chronic diseases. No previous hospitalizations.', 'No known allergies', 'Multivitamin daily',
        'Continue healthy lifestyle. Regular physical activity 150min/week. Balanced diet with vegetables and fruits.',
        '[{"name": "Vitamin D3", "dosage": "2000 IU", "frequency": "Once daily", "duration": "Ongoing"}]',
        '{"cholesterol_total": 4.8, "ldl": 2.9, "hdl": 1.5, "triglycerides": 1.1, "glucose_fasting": 5.1, "vitamin_d": 42}',
        'Total Cholesterol: 4.8 mmol/L (optimal), LDL: 2.9 mmol/L (optimal), HDL: 1.5 mmol/L (good), Glucose: 5.1 mmol/L (normal), Vitamin D: 42 ng/mL (sufficient)',
        'Overall excellent health. All parameters within normal range. Patient maintains active lifestyle and balanced diet. No concerns identified.',
        'Continue current healthy habits. Next routine checkup in 12 months. Regular dental checkup recommended.',
        NOW() + INTERVAL '11 months',
        'active', 'normal'
    );

    -- Patient 2: Gastritis at Regional Hospital
    INSERT INTO medical_data (
        user_id, clinic_id, visit_date, visit_type, department, attending_doctor,
        chief_complaint, symptoms, symptom_duration, pain_level,
        temperature, blood_pressure_systolic, blood_pressure_diastolic,
        heart_rate, weight, height, bmi,
        diagnosis, diagnosis_code, secondary_diagnoses, severity,
        medical_history, surgical_history, allergies, current_medications,
        treatment_plan, prescribed_medications, procedures_performed, imaging_ordered,
        lab_results, lab_results_summary,
        doctor_notes, follow_up_instructions, follow_up_date, restrictions,
        insurance_info, billing_code,
        record_status, confidentiality_level
    ) VALUES (
        _patient2_id, 2,
        NOW() - INTERVAL '5 days', 'routine', 'Gastroenterology', 'Dr. Maria Ivanova',
        'Persistent stomach pain and heartburn', 'Epigastric pain, heartburn after meals, nausea, belching', '3 weeks', 7,
        36.8, 138, 85, 78, 82.0, 175.0, 26.78,
        'Chronic gastritis with reflux esophagitis', 'K29.5', '["GERD - K21.9"]', 'moderate',
        'History of irregular meals and stress. Smoker (10 cigarettes/day for 5 years).', 'Appendectomy 2015', 'Aspirin (causes gastric upset)', 'None currently',
        'Strict dietary modifications: avoid spicy, fried, acidic foods. Small frequent meals. Elevate head of bed. Smoking cessation strongly advised. Proton pump inhibitor therapy. Follow-up endoscopy if symptoms persist.',
        '[{"name": "Omeprazole", "dosage": "20mg", "frequency": "Once daily before breakfast", "duration": "8 weeks"},
          {"name": "Sucralfate", "dosage": "1g", "frequency": "4 times daily before meals and bedtime", "duration": "4 weeks"},
          {"name": "Probiotics", "dosage": "1 capsule", "frequency": "Once daily", "duration": "4 weeks"}]',
        'Upper endoscopy (EGDS) performed',
        'Upper GI endoscopy scheduled in 6 weeks if no improvement',
        '{"h_pylori_test": "negative", "hemoglobin": 12.8, "ferritin": 45, "liver_enzymes_alt": 28, "liver_enzymes_ast": 24}',
        'H. pylori test: Negative. Hemoglobin: 12.8 g/dL (low-normal). Ferritin: 45 ng/mL (low-normal). Liver enzymes: Normal.',
        'Patient presents with classic symptoms of gastritis and GERD. Endoscopy shows moderate inflammation of gastric mucosa and lower esophagus. No ulcers or malignancy detected. H. pylori negative. Lifestyle modifications crucial for management. Patient counseled on smoking cessation and stress management.',
        'Return in 2 weeks for symptom review. Upper GI endoscopy in 6 weeks if symptoms persist despite treatment. Keep food diary. Avoid: coffee, alcohol, chocolate, mint, carbonated drinks. Eat dinner 3+ hours before bedtime. Consider smoking cessation program.',
        NOW() + INTERVAL '2 weeks',
        'Avoid heavy lifting or strenuous exercise for 1 week. No alcohol for duration of treatment. Smoking cessation recommended.',
        'State Insurance Policy #8765432', 'K29.5',
        'active', 'normal'
    );

    -- Patient 2: Hypertension follow-up at Regional Hospital
    INSERT INTO medical_data (
        user_id, clinic_id, visit_date, visit_type, department, attending_doctor,
        chief_complaint, symptoms, symptom_duration, pain_level,
        temperature, blood_pressure_systolic, blood_pressure_diastolic,
        heart_rate, weight, height, bmi,
        diagnosis, diagnosis_code, severity,
        medical_history, family_history, allergies, current_medications,
        treatment_plan, prescribed_medications, lab_tests_ordered,
        lab_results, lab_results_summary,
        doctor_notes, follow_up_instructions, follow_up_date, restrictions,
        record_status, confidentiality_level
    ) VALUES (
        _patient2_id, 2,
        NOW() - INTERVAL '20 days', 'follow-up', 'Cardiology', 'Dr. Petr Kuznetsov',
        'Follow-up for high blood pressure', 'Occasional headaches, fatigue', '6 months', 3,
        36.7, 145, 92, 76, 82.0, 175.0, 26.78,
        'Essential (primary) hypertension, Stage 2', 'I10', 'moderate',
        'Diagnosed with hypertension 6 months ago. Currently on medication. Smoker.', 'Father had MI at age 58. Mother has hypertension.', 'Aspirin',
        'Enalapril 10mg once daily',
        'Continue antihypertensive medication. Increase to 20mg if BP remains elevated. DASH diet. Reduce sodium intake (<2g/day). Regular aerobic exercise 30min 5x/week. Weight reduction target: 5kg. Smoking cessation critical. Home BP monitoring twice daily.',
        '[{"name": "Enalapril", "dosage": "20mg", "frequency": "Once daily in morning", "duration": "Ongoing"},
          {"name": "Amlodipine", "dosage": "5mg", "frequency": "Once daily", "duration": "Ongoing"},
          {"name": "Aspirin", "dosage": "75mg", "frequency": "Once daily with food", "duration": "Ongoing"}]',
        'Lipid panel, kidney function (creatinine, eGFR), electrolytes, HbA1c',
        '{"cholesterol_total": 6.2, "ldl": 4.1, "hdl": 1.1, "triglycerides": 2.3, "creatinine": 88, "egfr": 85, "sodium": 140, "potassium": 4.2, "hba1c": 5.6}',
        'Total Cholesterol: 6.2 mmol/L (high). LDL: 4.1 mmol/L (high). HDL: 1.1 mmol/L (low). Triglycerides: 2.3 mmol/L (borderline high). Kidney function: Normal. HbA1c: 5.6% (normal).',
        'BP remains elevated despite current medication (145/92). Increasing Enalapril dose and adding Amlodipine. Lipid profile shows elevated cholesterol - will monitor, may need statin therapy if lifestyle changes insufficient. Kidney function normal. ECG shows no acute abnormalities. Patient advised on cardiovascular risk factors: smoking, obesity, sedentary lifestyle. Strong emphasis on lifestyle modifications.',
        'Keep BP diary - measure morning and evening. Return in 4 weeks for BP check and medication adjustment. Target BP <140/90. Schedule: dietitian consultation for DASH diet, smoking cessation program enrollment. Repeat lipid panel in 3 months.',
        NOW() + INTERVAL '4 weeks',
        'Limit sodium to <2g/day. Limit alcohol. No heavy weight lifting until BP controlled. Regular gentle aerobic exercise encouraged.',
        'active', 'normal'
    );

    -- Patient 3: Allergic rhinitis at Medical Center Plus
    INSERT INTO medical_data (
        user_id, clinic_id, visit_date, visit_type, department, attending_doctor,
        chief_complaint, symptoms, symptom_duration, pain_level,
        temperature, blood_pressure_systolic, blood_pressure_diastolic,
        heart_rate, respiratory_rate, oxygen_saturation, weight, height, bmi,
        diagnosis, diagnosis_code, severity,
        medical_history, family_history, allergies, current_medications,
        treatment_plan, prescribed_medications, referrals,
        lab_results, lab_results_summary,
        doctor_notes, follow_up_instructions, follow_up_date, restrictions,
        record_status, confidentiality_level
    ) VALUES (
        _patient3_id, 3,
        NOW() - INTERVAL '7 days', 'routine', 'Allergy & Immunology', 'Dr. Alexey Sidorov',
        'Seasonal allergies - nasal congestion and itchy eyes', 'Nasal congestion, sneezing, runny nose, itchy watery eyes, throat irritation', '2 weeks (seasonal)', 2,
        36.5, 115, 72, 68, 16, 99, 58.0, 168.0, 20.55,
        'Allergic rhinitis, seasonal (hay fever)', 'J30.1', 'mild',
        'History of seasonal allergies since childhood. Worse in spring (March-May).', 'Mother has asthma. Sister has eczema.', 'Birch pollen, grass pollen', 'Cetirizine 10mg as needed',
        'Allergen avoidance strategies. Antihistamine therapy during allergy season. Keep windows closed during high pollen count days. Use air purifier with HEPA filter. Nasal saline rinses. Consider allergen immunotherapy if symptoms severe/prolonged.',
        '[{"name": "Cetirizine", "dosage": "10mg", "frequency": "Once daily at bedtime", "duration": "During allergy season"},
          {"name": "Fluticasone nasal spray", "dosage": "2 sprays each nostril", "frequency": "Once daily in morning", "duration": "During allergy season"},
          {"name": "Artificial tears", "dosage": "As needed", "frequency": "For eye symptoms", "duration": "As needed"}]',
        'Allergist for comprehensive allergy testing and possible immunotherapy consultation',
        '{"total_ige": 245, "specific_ige_birch": "4.2", "specific_ige_grass": "3.8", "eosinophils": 6.2}',
        'Total IgE: 245 kU/L (elevated). Specific IgE - Birch pollen: 4.2 kU/L (Class 3 - high). Specific IgE - Grass pollen: 3.8 kU/L (Class 3 - high). Eosinophils: 6.2% (elevated).',
        'Patient has classic allergic rhinitis triggered by tree and grass pollens. Blood tests confirm sensitization to birch and grass pollens. Symptoms are seasonal, predictable pattern. Currently mild-moderate severity. Responding well to antihistamines. Nasal corticosteroid added for better symptom control. Discussed allergen avoidance measures and potential for immunotherapy if symptoms worsen or become year-round.',
        'Start medications 1-2 weeks before typical allergy season. Monitor local pollen forecast. Keep allergy diary to track symptoms and triggers. Return if symptoms worsen or become year-round. Allergist consultation scheduled for next month to discuss immunotherapy options.',
        NOW() + INTERVAL '4 weeks',
        'Avoid outdoor activities during peak pollen times (early morning 5-10am). Shower after being outdoors. Keep car/home windows closed during high pollen days.',
        'active', 'normal'
    );

    -- Patient 3: Sports injury - ankle sprain at Medical Center Plus
    INSERT INTO medical_data (
        user_id, clinic_id, visit_date, visit_type, department, attending_doctor,
        chief_complaint, symptoms, symptom_duration, pain_level,
        temperature, blood_pressure_systolic, blood_pressure_diastolic,
        heart_rate, weight, height, bmi,
        diagnosis, diagnosis_code, severity,
        medical_history, allergies, current_medications,
        treatment_plan, prescribed_medications, procedures_performed, imaging_ordered,
        lab_results_summary,
        doctor_notes, follow_up_instructions, follow_up_date, restrictions,
        record_status, confidentiality_level
    ) VALUES (
        _patient3_id, 3,
        NOW() - INTERVAL '12 days', 'emergency', 'Orthopedics', 'Dr. Natalia Fedorova',
        'Right ankle pain and swelling after fall while playing basketball', 'Ankle pain, swelling, difficulty bearing weight, bruising', '6 hours', 8,
        36.6, 118, 75, 82, 58.0, 168.0, 20.55,
        'Sprain of lateral ligament of right ankle', 'S93.41', 'moderate',
        'Active lifestyle. Plays basketball regularly. No previous ankle injuries.', 'Birch pollen, grass pollen', 'Cetirizine (for allergies)',
        'RICE protocol: Rest, Ice, Compression, Elevation. Immobilization with ankle brace. Non-weight bearing for 48 hours, then gradual weight bearing as tolerated. Physical therapy after initial healing. NSAIDs for pain management.',
        '[{"name": "Ibuprofen", "dosage": "400mg", "frequency": "Three times daily with food", "duration": "7 days"},
          {"name": "Topical diclofenac gel", "dosage": "Apply to affected area", "frequency": "2-3 times daily", "duration": "10 days"}]',
        'Physical examination, ankle palpation, range of motion assessment',
        'X-ray right ankle - AP and lateral views',
        'X-ray: No fracture identified. Soft tissue swelling noted. Lateral ligament strain confirmed clinically.',
        'Patient sustained inversion injury to right ankle during basketball. Moderate swelling and tenderness over lateral malleolus and ATFL (anterior talofibular ligament). Grade 2 ankle sprain (partial ligament tear). X-ray negative for fracture. Patient provided with ankle brace and crutches. Instructed on RICE protocol. Pain well controlled with NSAIDs. Weight bearing to be progressed as tolerated.',
        'Apply ice 20min every 2-3 hours for first 48 hours. Keep ankle elevated above heart level. Wear ankle brace continuously for 2 weeks except when showering. Begin gentle range of motion exercises after 3 days. Physical therapy referral provided - start in 1 week. Return in 2 weeks for reassessment. Return sooner if: increased pain, numbness, skin color changes, or inability to bear weight.',
        NOW() + INTERVAL '2 weeks',
        'No weight bearing for first 48 hours (use crutches). Gradual return to weight bearing as tolerated. No sports or running for 4-6 weeks. May swim after 2 weeks if comfortable. Follow physical therapy program strictly.',
        'active', 'normal'
    );

END $$;
