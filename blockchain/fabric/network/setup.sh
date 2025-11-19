#!/bin/bash
# Blockchain network setup script for Medical Data Exchange
set -e

CHANNEL_NAME="medical-channel"
CHAINCODE_NAME="medical-access"
CHAINCODE_VERSION="1.0"
CHAINCODE_SEQUENCE="1"

echo "========================================="
echo "Medical Data Exchange - Blockchain Setup"
echo "========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

function printStep() {
    echo -e "${GREEN}>>> $1${NC}"
}

function printError() {
    echo -e "${RED}ERROR: $1${NC}"
}

function printWarning() {
    echo -e "${YELLOW}WARNING: $1${NC}"
}

# Step 1: Generate crypto material
printStep "Step 1: Generating crypto material..."
if [ ! -d "crypto-config/peerOrganizations" ]; then
    ./generate-crypto.sh
else
    printWarning "Crypto material already exists. Skipping generation."
fi

# Step 2: Start the network
printStep "Step 2: Starting Hyperledger Fabric network..."
docker-compose -f docker-compose-fabric.yaml up -d

# Wait for containers to start
printStep "Waiting for containers to be ready..."
sleep 10

# Check if all containers are running
printStep "Checking container status..."
docker ps --format "table {{.Names}}\t{{.Status}}" | grep -E "orderer|peer0|cli"

# Step 3: Create channel
printStep "Step 3: Creating channel '${CHANNEL_NAME}'..."
docker exec cli bash -c "
    cd /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts && \
    configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ${CHANNEL_NAME}.tx -channelID ${CHANNEL_NAME} 2>/dev/null || echo 'Channel tx already exists' && \
    cd .. && \
    peer channel create -o orderer.medical.com:7050 -c ${CHANNEL_NAME} -f ./channel-artifacts/${CHANNEL_NAME}.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/tlscacerts/tlsca.medical.com-cert.pem 2>/dev/null || echo 'Channel already created'
"

printStep "Channel created successfully!"

# Step 4: Join peers to channel
printStep "Step 4: Joining peers to channel..."

# Join Clinic1 peer
printStep "Joining peer0.clinic1 to channel..."
docker exec cli bash -c "
    peer channel join -b ${CHANNEL_NAME}.block
"

# Switch to Clinic2 peer and join
printStep "Joining peer0.clinic2 to channel..."
docker exec cli bash -c "
    export CORE_PEER_LOCALMSPID=Clinic2MSP && \
    export CORE_PEER_ADDRESS=peer0.clinic2.medical.com:9051 && \
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic2.medical.com/peers/peer0.clinic2.medical.com/tls/ca.crt && \
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic2.medical.com/users/Admin@clinic2.medical.com/msp && \
    peer channel join -b ${CHANNEL_NAME}.block
"

printStep "All peers joined the channel!"

# Step 5: Update anchor peers
printStep "Step 5: Updating anchor peers..."

# Update anchor peer for Clinic1
docker exec cli bash -c "
    cd channel-artifacts && \
    configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate Clinic1MSPanchors.tx -channelID ${CHANNEL_NAME} -asOrg Clinic1MSP 2>/dev/null && \
    cd .. && \
    peer channel update -o orderer.medical.com:7050 -c ${CHANNEL_NAME} -f ./channel-artifacts/Clinic1MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/tlscacerts/tlsca.medical.com-cert.pem
"

# Update anchor peer for Clinic2
docker exec cli bash -c "
    export CORE_PEER_LOCALMSPID=Clinic2MSP && \
    export CORE_PEER_ADDRESS=peer0.clinic2.medical.com:9051 && \
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic2.medical.com/peers/peer0.clinic2.medical.com/tls/ca.crt && \
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic2.medical.com/users/Admin@clinic2.medical.com/msp && \
    cd channel-artifacts && \
    configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate Clinic2MSPanchors.tx -channelID ${CHANNEL_NAME} -asOrg Clinic2MSP 2>/dev/null && \
    cd .. && \
    peer channel update -o orderer.medical.com:7050 -c ${CHANNEL_NAME} -f ./channel-artifacts/Clinic2MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/tlscacerts/tlsca.medical.com-cert.pem
"

printStep "Anchor peers updated!"

# Step 6: Package chaincode
printStep "Step 6: Packaging chaincode..."
docker exec cli bash -c "
    cd /opt/gopath/src/github.com/chaincode/fabric/chaincode/medical-access && \
    GO111MODULE=on go mod vendor 2>/dev/null || true && \
    cd /opt/gopath/src/github.com/hyperledger/fabric/peer && \
    peer lifecycle chaincode package ${CHAINCODE_NAME}.tar.gz --path /opt/gopath/src/github.com/chaincode/fabric/chaincode/medical-access/ --lang golang --label ${CHAINCODE_NAME}_${CHAINCODE_VERSION}
"

printStep "Chaincode packaged!"

# Step 7: Install chaincode on peers
printStep "Step 7: Installing chaincode on peers..."

# Install on Clinic1 peer
printStep "Installing on peer0.clinic1..."
docker exec cli bash -c "
    peer lifecycle chaincode install ${CHAINCODE_NAME}.tar.gz
"

# Install on Clinic2 peer
printStep "Installing on peer0.clinic2..."
docker exec cli bash -c "
    export CORE_PEER_LOCALMSPID=Clinic2MSP && \
    export CORE_PEER_ADDRESS=peer0.clinic2.medical.com:9051 && \
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic2.medical.com/peers/peer0.clinic2.medical.com/tls/ca.crt && \
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic2.medical.com/users/Admin@clinic2.medical.com/msp && \
    peer lifecycle chaincode install ${CHAINCODE_NAME}.tar.gz
"

printStep "Chaincode installed on all peers!"

# Step 8: Query installed chaincode to get package ID
printStep "Step 8: Querying installed chaincode..."
PACKAGE_ID=$(docker exec cli peer lifecycle chaincode queryinstalled | grep ${CHAINCODE_NAME}_${CHAINCODE_VERSION} | awk '{print $3}' | sed 's/,//')
echo "Package ID: ${PACKAGE_ID}"

# Step 9: Approve chaincode for both orgs
printStep "Step 9: Approving chaincode for organizations..."

# Approve for Clinic1
printStep "Approving for Clinic1MSP..."
docker exec cli bash -c "
    peer lifecycle chaincode approveformyorg -o orderer.medical.com:7050 --channelID ${CHANNEL_NAME} --name ${CHAINCODE_NAME} --version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/tlscacerts/tlsca.medical.com-cert.pem
"

# Approve for Clinic2
printStep "Approving for Clinic2MSP..."
docker exec cli bash -c "
    export CORE_PEER_LOCALMSPID=Clinic2MSP && \
    export CORE_PEER_ADDRESS=peer0.clinic2.medical.com:9051 && \
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic2.medical.com/peers/peer0.clinic2.medical.com/tls/ca.crt && \
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic2.medical.com/users/Admin@clinic2.medical.com/msp && \
    peer lifecycle chaincode approveformyorg -o orderer.medical.com:7050 --channelID ${CHANNEL_NAME} --name ${CHAINCODE_NAME} --version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/tlscacerts/tlsca.medical.com-cert.pem
"

printStep "Chaincode approved by all organizations!"

# Step 10: Check commit readiness
printStep "Step 10: Checking commit readiness..."
docker exec cli peer lifecycle chaincode checkcommitreadiness --channelID ${CHANNEL_NAME} --name ${CHAINCODE_NAME} --version ${CHAINCODE_VERSION} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/tlscacerts/tlsca.medical.com-cert.pem --output json

# Step 11: Commit chaincode
printStep "Step 11: Committing chaincode to channel..."
docker exec cli bash -c "
    peer lifecycle chaincode commit -o orderer.medical.com:7050 --channelID ${CHANNEL_NAME} --name ${CHAINCODE_NAME} --version ${CHAINCODE_VERSION} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/tlscacerts/tlsca.medical.com-cert.pem --peerAddresses peer0.clinic1.medical.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic1.medical.com/peers/peer0.clinic1.medical.com/tls/ca.crt --peerAddresses peer0.clinic2.medical.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic2.medical.com/peers/peer0.clinic2.medical.com/tls/ca.crt
"

printStep "Chaincode committed successfully!"

# Step 12: Query committed chaincode
printStep "Step 12: Querying committed chaincode..."
docker exec cli peer lifecycle chaincode querycommitted --channelID ${CHANNEL_NAME} --name ${CHAINCODE_NAME}

# Step 13: Initialize ledger
printStep "Step 13: Initializing ledger..."
docker exec cli peer chaincode invoke -o orderer.medical.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/tlscacerts/tlsca.medical.com-cert.pem -C ${CHANNEL_NAME} -n ${CHAINCODE_NAME} --peerAddresses peer0.clinic1.medical.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic1.medical.com/peers/peer0.clinic1.medical.com/tls/ca.crt --peerAddresses peer0.clinic2.medical.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/clinic2.medical.com/peers/peer0.clinic2.medical.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}Blockchain network is ready!${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo "Channel: ${CHANNEL_NAME}"
echo "Chaincode: ${CHAINCODE_NAME}"
echo "Version: ${CHAINCODE_VERSION}"
echo ""
echo "You can now use the blockchain service to interact with the network."
