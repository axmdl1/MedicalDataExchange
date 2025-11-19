package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	medical "github.com/axmdl1/MedicalDataExchange/backend/blockchain-contracts/medical"
)

// connectWithRetry пытается подключиться к RPC с повторами и таймаутом.
// rpcURL — "http://ganache:8545" (по умолчанию).
func connectWithRetry(rpcURL string, attempts int, delay time.Duration) *ethclient.Client {
	for i := 1; i <= attempts; i++ {
		client, err := ethclient.Dial(rpcURL)
		if err == nil {
			log.Printf("✅ Connected to Ganache at %s (attempt %d/%d)\n", rpcURL, i, attempts)
			return client
		}
		log.Printf("⏳ Waiting for Ganache (%d/%d): %v — retry in %s\n", i, attempts, err, delay)
		time.Sleep(delay)
	}
	log.Fatalf("❌ Could not connect to Ganache at %s after %d attempts", rpcURL, attempts)
	return nil
}

func main() {
	// -----------------------------------------------------
	// 1. Connect to Ganache
	// -----------------------------------------------------

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		rpcURL = "http://ganache:8545"
	}
	client := connectWithRetry(rpcURL, 30, 1*time.Second)
	// (если нужен объект err, можно извлечь client == nil)

	//client, err := ethclient.Dial("http://ganache:8545")
	/*if err != nil {
		log.Fatal("Failed to connect to Ganache:", err)
	}*/
	fmt.Println("Connected to Ganache")

	// -----------------------------------------------------
	// 2. Use first Ganache account (from docker logs)
	// -----------------------------------------------------
	ganachePrivKey := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	privateKey, err := crypto.HexToECDSA(ganachePrivKey)
	if err != nil {
		log.Fatal("Invalid private key:", err)
	}

	publicKey := privateKey.Public()
	pubECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("public key conversion failed")
	}

	fromAddress := crypto.PubkeyToAddress(*pubECDSA)
	fmt.Println("Using address:", fromAddress.Hex())

	// -----------------------------------------------------
	// 3. Build transactor (NO MANUAL NONCE!)
	// -----------------------------------------------------
	chainID := big.NewInt(1337) // Ganache chain ID
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(6_000_000)

	// -----------------------------------------------------
	// 4. Deploy Smart Contract
	// -----------------------------------------------------
	contractAddr, tx, contract, err := medical.DeployMedicalAccess(auth, client)
	if err != nil {
		log.Fatal("Deployment failed:", err)
	}

	fmt.Println("----------------------------------------------------")
	fmt.Println(" Smart Contract Deployed")
	fmt.Println(" Address:", contractAddr.Hex())
	fmt.Println(" Tx Hash:", tx.Hash().Hex())
	fmt.Println("----------------------------------------------------")

	log.Printf("🚀 MEDICAL ACCESS CONTRACT DEPLOYED AT: %s\n", contractAddr.Hex())

	time.Sleep(time.Second * 2)

	// -----------------------------------------------------
	// 5. Request Access (patient → blockchain)
	// -----------------------------------------------------
	requestID := "r1"

	tx1, err := contract.RequestAccess(auth, requestID, "rec1", "patient1", "clinic1")
	if err != nil {
		log.Fatal("RequestAccess failed:", err)
	}
	fmt.Println("RequestAccess TX:", tx1.Hash().Hex())

	time.Sleep(time.Second * 2)

	// -----------------------------------------------------
	// 6. Clinic Approves (generate tokenHash)
	// -----------------------------------------------------
	token := "patient-secret-token-999"
	tokenHash := crypto.Keccak256Hash([]byte(token))

	expiresAt := big.NewInt(time.Now().Add(10 * time.Minute).Unix())

	// auth.Nonce = nil <-- Let Go-Ethereum select the correct nonce
	tx2, err := contract.ApproveAccess(auth, requestID, tokenHash, expiresAt)
	if err != nil {
		log.Fatal("ApproveAccess failed:", err)
	}
	fmt.Println("ApproveAccess TX:", tx2.Hash().Hex())
	fmt.Println("Token (give to patient):", token)
	fmt.Println("Token Hash Stored On-Chain:", tokenHash.Hex())

	time.Sleep(time.Second * 2)

	// -----------------------------------------------------
	// 7. Patient tries to access → verifyToken
	// -----------------------------------------------------
	callOpts := &bind.CallOpts{Context: context.Background()}

	okVerify, err := contract.VerifyToken(callOpts, requestID, tokenHash)
	if err != nil {
		log.Fatal("VerifyToken failed:", err)
	}
	fmt.Println("VerifyToken result:", okVerify)

	// -----------------------------------------------------
	// 8. Get On-Chain State
	// -----------------------------------------------------
	rec, err := contract.GetRequest(callOpts, requestID)
	if err != nil {
		log.Fatal("GetRequest failed:", err)
	}

	fmt.Println("------------------ On-Chain State ------------------")
	fmt.Println("Record ID:", rec.RecordId)
	fmt.Println("Patient:", rec.Patient)
	fmt.Println("Clinic:", rec.Clinic)

	tokenHashOnChain := common.BytesToHash(rec.TokenHash[:])
	fmt.Println("Token Hash:", tokenHashOnChain.Hex())
	fmt.Println("ExpiresAt:", rec.ExpiresAt)
	fmt.Println("Status:", rec.Status)
	fmt.Println("----------------------------------------------------")
}
