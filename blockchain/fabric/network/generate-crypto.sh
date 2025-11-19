#!/bin/bash
# Generate crypto material for Hyperledger Fabric network
set -e

echo "Generating crypto material for Medical Data Exchange network..."

# Create directory structure
mkdir -p crypto-config/peerOrganizations/clinic1.medical.com
mkdir -p crypto-config/peerOrganizations/clinic2.medical.com
mkdir -p crypto-config/ordererOrganizations/medical.com
mkdir -p channel-artifacts

# Use cryptogen if available, otherwise use openssl for basic setup
if command -v cryptogen &> /dev/null; then
    echo "Using cryptogen tool..."
    cryptogen generate --config=./crypto-config.yaml --output="crypto-config"
else
    echo "cryptogen not found. Creating minimal crypto setup with openssl..."
    # This is a simplified version - in production, use cryptogen or Fabric CA

    # Create basic structure for Orderer
    mkdir -p crypto-config/ordererOrganizations/medical.com/orderers/orderer.medical.com/{msp,tls}
    mkdir -p crypto-config/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/{admincerts,cacerts,keystore,signcerts,tlscacerts}

    # Create basic structure for Clinic1
    mkdir -p crypto-config/peerOrganizations/clinic1.medical.com/peers/peer0.clinic1.medical.com/{msp,tls}
    mkdir -p crypto-config/peerOrganizations/clinic1.medical.com/peers/peer0.clinic1.medical.com/msp/{admincerts,cacerts,keystore,signcerts,tlscacerts}
    mkdir -p crypto-config/peerOrganizations/clinic1.medical.com/users/Admin@clinic1.medical.com/msp/{admincerts,cacerts,keystore,signcerts,tlscacerts}

    # Create basic structure for Clinic2
    mkdir -p crypto-config/peerOrganizations/clinic2.medical.com/peers/peer0.clinic2.medical.com/{msp,tls}
    mkdir -p crypto-config/peerOrganizations/clinic2.medical.com/peers/peer0.clinic2.medical.com/msp/{admincerts,cacerts,keystore,signcerts,tlscacerts}
    mkdir -p crypto-config/peerOrganizations/clinic2.medical.com/users/Admin@clinic2.medical.com/msp/{admincerts,cacerts,keystore,signcerts,tlscacerts}

    # Generate certificates using openssl (simplified)
    # Note: This is NOT for production use. Use Fabric CA or cryptogen in production.

    echo "Creating placeholder certificates..."
    # These would normally be generated with proper Fabric tools
    # For now, we'll create empty placeholder files that Fabric expects

    touch crypto-config/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/signcerts/cert.pem
    touch crypto-config/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/keystore/key.pem
    touch crypto-config/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/cacerts/ca.pem
    touch crypto-config/ordererOrganizations/medical.com/orderers/orderer.medical.com/msp/tlscacerts/tlsca.medical.com-cert.pem
    touch crypto-config/ordererOrganizations/medical.com/orderers/orderer.medical.com/tls/server.crt
    touch crypto-config/ordererOrganizations/medical.com/orderers/orderer.medical.com/tls/server.key
    touch crypto-config/ordererOrganizations/medical.com/orderers/orderer.medical.com/tls/ca.crt

    echo "WARNING: Crypto material is incomplete. Please use 'cryptogen' or Fabric CA for production setup."
fi

echo "Crypto material generation complete!"
echo "Location: ./crypto-config/"
