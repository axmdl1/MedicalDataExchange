package handler

import (
	"encoding/json"
	"net/http"

	"github.com/axmdl1/MedicalDataExchange/core-service/internal/blockchain"
	"github.com/axmdl1/MedicalDataExchange/core-service/internal/repository"
	"github.com/ethereum/go-ethereum/crypto"
)

type DataHandler struct {
	Repo repository.MedicalRepository
	BC   *blockchain.Blockchain
}

func NewDataHandler(repo repository.MedicalRepository, bc *blockchain.Blockchain) *DataHandler {
	return &DataHandler{Repo: repo, BC: bc}
}

func (h *DataHandler) GetMedicalData(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("request_id")
	token := r.URL.Query().Get("token")
	recordID := r.URL.Query().Get("record_id")

	if requestID == "" || token == "" || recordID == "" {
		http.Error(w, "missing params", 400)
		return
	}

	// Hash token like clinic did
	tokenHash := crypto.Keccak256Hash([]byte(token))

	// 1) CHECK BLOCKCHAIN PERMISSION
	ok, err := h.BC.VerifyToken(requestID, tokenHash)
	if err != nil || !ok {
		http.Error(w, "access denied: expired or invalid token", 403)
		return
	}

	// 2) GET DATA FROM CLINIC DB
	data, err := h.Repo.GetMedicalData(recordID)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	json.NewEncoder(w).Encode(data)
}
