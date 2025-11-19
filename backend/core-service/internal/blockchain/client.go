package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	medical "github.com/axmdl1/MedicalDataExchange/backend/blockchain-contracts/medical"
)

// AccessRequest удобная локальная структура
type AccessRequest struct {
	RecordId  string
	Patient   string
	Clinic    string
	TokenHash [32]byte
	ExpiresAt int64
	Status    string
}

type Blockchain struct {
	Client   *ethclient.Client
	Contract *medical.MedicalAccess
	chainID  *big.Int
}

// NewBlockchain подключается к RPC и бинит контракт.
// ОЖИДАЕТ: env CONTRACT_ADDR (0x...), RPC_URL (например http://ganache:8545)
func NewBlockchain() *Blockchain {
	rpc := os.Getenv("RPC_URL")
	if rpc == "" {
		rpc = "http://ganache:8545"
	}
	contractAddr := os.Getenv("CONTRACT_ADDR")
	if contractAddr == "" {
		log.Fatal("CONTRACT_ADDR not set")
	}

	chainID := big.NewInt(1337)
	if v := os.Getenv("CHAIN_ID"); v != "" {
		if i, ok := new(big.Int).SetString(v, 10); ok {
			chainID = i
		}
	}

	cli, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatalf("failed to dial rpc %s: %v", rpc, err)
	}

	contract, err := medical.NewMedicalAccess(common.HexToAddress(contractAddr), cli)
	if err != nil {
		log.Fatalf("failed to bind contract: %v", err)
	}

	return &Blockchain{
		Client:   cli,
		Contract: contract,
		chainID:  chainID,
	}
}

// VerifyToken проверяет токен (view call)
func (bc *Blockchain) VerifyToken(requestID string, tokenHash common.Hash) (bool, error) {
	ok, err := bc.Contract.VerifyToken(&bind.CallOpts{Context: context.Background()}, requestID, tokenHash)
	return ok, err
}

// GetRequest читает request с цепочки и возвращает локальную структуру
func (bc *Blockchain) GetRequest(requestID string) (*AccessRequest, error) {
	raw, err := bc.Contract.GetRequest(&bind.CallOpts{Context: context.Background()}, requestID)
	if err != nil {
		return nil, err
	}
	return &AccessRequest{
		RecordId:  raw.RecordId,
		Patient:   raw.Patient,
		Clinic:    raw.Clinic,
		TokenHash: raw.TokenHash,
		ExpiresAt: raw.ExpiresAt.Int64(),
		Status:    raw.Status,
	}, nil
}

// Approve вызывает approveAccess транзакцию (использует CLINIC_PRIVATE_KEY)
func (bc *Blockchain) Approve(requestID string, tokenHash common.Hash, expiresAt int64) (string, error) {
	sk := os.Getenv("CLINIC_PRIVATE_KEY")
	if sk == "" {
		return "", fmt.Errorf("CLINIC_PRIVATE_KEY not set")
	}
	priv, err := crypto.HexToECDSA(sk)
	if err != nil {
		return "", fmt.Errorf("invalid clinic private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(priv, bc.chainID)
	if err != nil {
		return "", err
	}
	auth.GasLimit = uint64(5_000_000)
	auth.Value = big.NewInt(0)

	tx, err := bc.Contract.ApproveAccess(auth, requestID, tokenHash, big.NewInt(expiresAt))
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// Request creates on-chain request (optionally, clinic could instead create on DB).
// We expose it in case you want to write request to chain directly.
func (bc *Blockchain) Request(requestID, recordID, patient, clinic string, privKeyHex string) (string, error) {
	// if privKeyHex == "" use CLINIC_PRIVATE_KEY
	if privKeyHex == "" {
		privKeyHex = os.Getenv("CLINIC_PRIVATE_KEY")
	}
	priv, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(priv, bc.chainID)
	if err != nil {
		return "", err
	}
	auth.GasLimit = uint64(5_000_000)
	auth.Value = big.NewInt(0)

	tx, err := bc.Contract.RequestAccess(auth, requestID, recordID, patient, clinic)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// Revoke (mark revoked on chain)
func (bc *Blockchain) Revoke(requestID string, privKeyHex string) (string, error) {
	// In your contract there's RequestRevoked event but no function — if you added function, call it here.
	// If not, we'll just return error placeholder.
	return "", fmt.Errorf("Revoke function is not implemented in the contract")
}
