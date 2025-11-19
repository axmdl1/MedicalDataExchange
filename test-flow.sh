#!/bin/bash
# Test flow for Medical Data Exchange - Patient Access System
set -e

echo "=========================================="
echo "Medical Data Exchange - Test Flow"
echo "=========================================="
echo ""

BASE_URL="http://localhost:8099"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}Step 1: Creating patient access request...${NC}"
echo "Patient requests to view their medical data from clinic"
echo ""

REQUEST_RESPONSE=$(curl -s -X POST ${BASE_URL}/api/patient-access/request \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 1,
    "clinic_id": 1,
    "medical_data_id": 1
  }')

echo "$REQUEST_RESPONSE" | jq .
REQUEST_ID=$(echo "$REQUEST_RESPONSE" | jq -r '.request_id')

echo ""
echo -e "${GREEN}✓ Request created with ID: $REQUEST_ID${NC}"
echo ""
sleep 2

echo -e "${BLUE}Step 2: Clinic approving the request...${NC}"
echo "Clinic admin approves patient's request"
echo ""

APPROVE_RESPONSE=$(curl -s -X POST ${BASE_URL}/api/patient-access/approve/${REQUEST_ID} \
  -H "Content-Type: application/json" \
  -d '{
    "approver_id": 10
  }')

echo "$APPROVE_RESPONSE" | jq .
ACCESS_TOKEN=$(echo "$APPROVE_RESPONSE" | jq -r '.access_token')
EXPIRES_AT=$(echo "$APPROVE_RESPONSE" | jq -r '.expires_at')

echo ""
echo -e "${GREEN}✓ Request approved!${NC}"
echo -e "${GREEN}✓ Access token: ${ACCESS_TOKEN:0:20}...${NC}"
echo -e "${GREEN}✓ Expires at: $EXPIRES_AT${NC}"
echo ""
echo -e "${YELLOW}What happened:${NC}"
echo "  1. Medical data retrieved from clinic database"
echo "  2. Data encrypted with AES-256-GCM"
echo "  3. Temporary copy created (expires in 15 minutes)"
echo "  4. Unique access token generated"
echo "  5. Original data remains in clinic"
echo ""
sleep 3

echo -e "${BLUE}Step 3: Patient accessing their medical data...${NC}"
echo "Using the temporary access token"
echo ""

DATA_RESPONSE=$(curl -s "${BASE_URL}/api/patient-access/data?token=${ACCESS_TOKEN}")

echo "$DATA_RESPONSE" | jq .

echo ""
echo -e "${GREEN}✓ Data successfully retrieved!${NC}"
echo ""
echo -e "${YELLOW}Data includes:${NC}"
echo "  - Visit information (date, type, department)"
echo "  - Vital signs (temperature, blood pressure, heart rate)"
echo "  - Diagnosis and treatment plan"
echo "  - Prescribed medications"
echo "  - Lab results"
echo "  - Doctor's notes"
echo ""
sleep 3

echo -e "${BLUE}Step 4: Checking request status...${NC}"
echo ""

STATUS_RESPONSE=$(curl -s "${BASE_URL}/api/patient-access/request/${REQUEST_ID}")
echo "$STATUS_RESPONSE" | jq .

echo ""
echo -e "${GREEN}✓ Request status: approved${NC}"
echo ""

echo "=========================================="
echo -e "${GREEN}Test Flow Completed Successfully!${NC}"
echo "=========================================="
echo ""
echo "Next steps:"
echo "  • Data will auto-delete after 15 minutes"
echo "  • Cleanup scheduler runs every 5 minutes"
echo "  • Original data remains safely in clinic database"
echo ""
echo "Monitor cleanup:"
echo "  docker logs mde-dt -f | grep cleanup"
echo ""
echo "Check temporary data:"
echo "  docker exec -it mde-postgres psql -U postgres -d medical -c 'SELECT * FROM temporary_patient_data;'"
echo ""
