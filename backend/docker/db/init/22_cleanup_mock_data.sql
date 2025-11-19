-- Cleanup old mock data and ensure only real medical records exist
-- This script removes test/mock data that doesn't match our comprehensive schema

-- Delete any medical_data records that have generic "test" diagnosis or missing critical fields
DELETE FROM medical_data
WHERE diagnosis = 'test'
   OR diagnosis IS NULL
   OR (diagnosis_code IS NULL AND visit_date IS NULL);

-- Delete any medical_data that doesn't have proper visit information
-- (These are likely old mock records created before schema enhancement)
DELETE FROM medical_data
WHERE visit_date IS NULL
  AND department IS NULL
  AND attending_doctor IS NULL;

-- Log what we have now
DO $$
DECLARE
    record_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO record_count FROM medical_data;
    RAISE NOTICE 'Remaining medical_data records: %', record_count;
END $$;
