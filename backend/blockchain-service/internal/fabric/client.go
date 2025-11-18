package fabric

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	fabricConfig "github.com/medical-data-exchange/blockchain-service/internal/config"
)

type FabricClient struct {
	gateway  *gateway.Gateway
	network  *gateway.Network
	contract *gateway.Contract
	config   fabricConfig.FabricConfig
}

func NewFabricClient(cfg fabricConfig.FabricConfig) (*FabricClient, error) {
	// Создаем wallet
	wallet, err := gateway.NewFileSystemWallet(cfg.WalletPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Проверяем наличие identity в wallet
	if !wallet.Exists(cfg.UserName) {
		err = populateWallet(wallet, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to populate wallet: %w", err)
		}
	}

	// Создаем gateway
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(cfg.NetworkConfig)),
		gateway.WithIdentity(wallet, cfg.UserName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gateway: %w", err)
	}

	// Получаем network
	network, err := gw.GetNetwork(cfg.ChannelName)
	if err != nil {
		gw.Close()
		return nil, fmt.Errorf("failed to get network: %w", err)
	}

	// Получаем contract
	contract := network.GetContract(cfg.ChaincodeName)

	log.Printf("Connected to Fabric network: channel=%s, chaincode=%s", cfg.ChannelName, cfg.ChaincodeName)

	return &FabricClient{
		gateway:  gw,
		network:  network,
		contract: contract,
		config:   cfg,
	}, nil
}

func (fc *FabricClient) Close() {
	if fc.gateway != nil {
		fc.gateway.Close()
	}
}

func (fc *FabricClient) GetContract() *gateway.Contract {
	return fc.contract
}

func populateWallet(wallet *gateway.Wallet, cfg fabricConfig.FabricConfig) error {

	log.Println("Wallet identity not found, creating new identity...")

	return fmt.Errorf("wallet identity not found, please configure certificates")
}
