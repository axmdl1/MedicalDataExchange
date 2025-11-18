package config

import (
	"os"
)

type Config struct {
	GRPCAddr string
	Fabric   FabricConfig
}

type FabricConfig struct {
	NetworkConfig  string
	ChannelName    string
	ChaincodeName  string
	OrgName        string
	UserName       string
	ConnectionPath string
	WalletPath     string
	CertPath       string
	KeyPath        string
	TLSCertPath    string
	PeerEndpoint   string
	GatewayPeer    string
	MSPID          string
}

func MustLoad() *Config {
	return &Config{
		GRPCAddr: getEnv("GRPC_ADDR", ":50054"),
		Fabric: FabricConfig{
			NetworkConfig:  getEnv("FABRIC_NETWORK_CONFIG", "./config/connection-profile.yaml"),
			ChannelName:    getEnv("FABRIC_CHANNEL", "medical-channel"),
			ChaincodeName:  getEnv("FABRIC_CHAINCODE", "medical-access"),
			OrgName:        getEnv("FABRIC_ORG", "Clinic1MSP"),
			UserName:       getEnv("FABRIC_USER", "Admin"),
			ConnectionPath: getEnv("FABRIC_CONNECTION_PATH", "./fabric-network"),
			WalletPath:     getEnv("FABRIC_WALLET_PATH", "./wallet"),
			CertPath:       getEnv("FABRIC_CERT_PATH", "./crypto/users/Admin@clinic1.medical.com/msp/signcerts/cert.pem"),
			KeyPath:        getEnv("FABRIC_KEY_PATH", "./crypto/users/Admin@clinic1.medical.com/msp/keystore/key.pem"),
			TLSCertPath:    getEnv("FABRIC_TLS_CERT", "./crypto/peers/peer0.clinic1.medical.com/tls/ca.crt"),
			PeerEndpoint:   getEnv("FABRIC_PEER_ENDPOINT", "localhost:7051"),
			GatewayPeer:    getEnv("FABRIC_GATEWAY_PEER", "peer0.clinic1.medical.com"),
			MSPID:          getEnv("FABRIC_MSP_ID", "Clinic1MSP"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
