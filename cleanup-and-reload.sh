#!/bin/bash
# Quick cleanup and reload of medical data without full rebuild

echo "========================================="
echo "Cleanup Mock Data and Reload Real Data"
echo "========================================="
echo ""

# Check if postgres is running
if ! docker ps | grep -q mde-postgres; then
    echo "Error: PostgreSQL container is not running"
    echo "Please run: docker-compose up -d"
    exit 1
fi

echo "Step 1: Cleaning up mock data..."
docker exec -i mde-postgres psql -U postgres -d medical <<EOF
-- Delete mock data
DELETE FROM medical_data WHERE diagnosis = 'test' OR diagnosis IS NULL;
DELETE FROM medical_data WHERE visit_date IS NULL AND department IS NULL;
EOF

echo ""
echo "Step 2: Checking current data..."
docker exec -i mde-postgres psql -U postgres -d medical -c "SELECT COUNT(*) as total_records FROM medical_data;"

echo ""
echo "Step 3: Loading real medical data..."
docker exec -i mde-postgres psql -U postgres -d medical < backend/docker/db/init/21_real_medical_data.sql

echo ""
echo "Step 4: Verification - Current medical records:"
docker exec -i mde-postgres psql -U postgres -d medical <<EOF
SELECT
    id,
    LEFT(diagnosis, 50) as diagnosis,
    diagnosis_code,
    department,
    visit_date::date
FROM medical_data
ORDER BY visit_date DESC;
EOF

echo ""
echo "========================================="
echo "✅ Cleanup complete!"
echo "========================================="
echo ""
echo "You should now see real medical records like:"
echo "  - Acute upper respiratory tract infection (J06.9)"
echo "  - Chronic gastritis with reflux esophagitis (K29.5)"
echo "  - Essential (primary) hypertension (I10)"
echo "  - Allergic rhinitis, seasonal (J30.1)"
echo ""
echo "Refresh your browser to see updated data!"
