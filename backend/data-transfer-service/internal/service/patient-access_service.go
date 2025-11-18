package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/model"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/repository"
	"github.com/rs/zerolog"
)

type PatientAccessService interface {
	// Пациент создает запрос на просмотр своих данных
	CreateAccessRequest(ctx context.Context, patientID, clinicID, medicalDataID int64) (*model.PatientAccessRequest, error)

	// Клиника одобряет запрос и создает временный доступ
	ApproveAccessRequest(ctx context.Context, requestID, approverID int64) (*model.TemporaryPatientData, error)

	// Клиника отклоняет запрос
	RejectAccessRequest(ctx context.Context, requestID int64) error

	// Получить запрос
	GetAccessRequest(ctx context.Context, requestID int64) (*model.PatientAccessRequest, error)

	// Список запросов пациента
	ListPatientAccessRequests(ctx context.Context, patientID *int64, clinicID *int64, status *string) ([]*model.PatientAccessRequest, error)

	// Получить временные данные по токену
	GetTemporaryData(ctx context.Context, accessToken string) (*model.MedicalData, error)

	// Отозвать доступ
	RevokeAccess(ctx context.Context, accessToken string) error

	// Очистить истекшие данные (вызывается периодически)
	CleanupExpiredData(ctx context.Context) error
}

type patientAccessService struct {
	repo           repository.DataTransferRepository
	encryptionKey  []byte
	accessDuration time.Duration // 15 минут
}

func NewPatientAccessService(repo repository.DataTransferRepository, encryptionSecret string) PatientAccessService {
	// Генерируем ключ из секрета
	hash := sha256.Sum256([]byte(encryptionSecret))

	return &patientAccessService{
		repo:           repo,
		encryptionKey:  hash[:],
		accessDuration: 15 * time.Minute,
	}
}

func (s *patientAccessService) CreateAccessRequest(ctx context.Context, patientID, clinicID, medicalDataID int64) (*model.PatientAccessRequest, error) {
	log := zerolog.Ctx(ctx).With().Str("component", "patient_access_service").Str("method", "CreateAccessRequest").Logger()

	// Проверяем, что медицинские данные принадлежат этой клинике
	medData, err := s.repo.GetMedicalData(ctx, medicalDataID)
	if err != nil {
		return nil, fmt.Errorf("medical data not found: %w", err)
	}

	if medData.ClinicID != clinicID {
		return nil, fmt.Errorf("medical data does not belong to this clinic")
	}

	if medData.UserID != patientID {
		return nil, fmt.Errorf("medical data does not belong to this patient")
	}

	request := &model.PatientAccessRequest{
		PatientID:     patientID,
		ClinicID:      clinicID,
		MedicalDataID: medicalDataID,
		Status:        "pending",
		RequestedAt:   time.Now(),
	}

	// Сохраняем в БД
	if err := s.repo.CreatePatientAccessRequest(ctx, request); err != nil {
		log.Error().Err(err).Msg("failed to create access request")
		return nil, err
	}

	log.Info().
		Int64("patient_id", patientID).
		Int64("clinic_id", clinicID).
		Int64("medical_data_id", medicalDataID).
		Msg("access request created")

	return request, nil
}

func (s *patientAccessService) ApproveAccessRequest(ctx context.Context, requestID, approverID int64) (*model.TemporaryPatientData, error) {
	log := zerolog.Ctx(ctx).With().Str("component", "patient_access_service").Str("method", "ApproveAccessRequest").Logger()

	// Получаем запрос
	request, err := s.repo.GetPatientAccessRequest(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("access request not found: %w", err)
	}

	if request.Status != "pending" {
		return nil, fmt.Errorf("request is not in pending status")
	}

	// Получаем медицинские данные
	medData, err := s.repo.GetMedicalData(ctx, request.MedicalDataID)
	if err != nil {
		return nil, fmt.Errorf("medical data not found: %w", err)
	}

	// Шифруем данные
	encryptedData, err := s.encryptMedicalData(medData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Генерируем токен доступа
	accessToken := s.generateAccessToken(request.PatientID, request.MedicalDataID)

	now := time.Now()
	expiresAt := now.Add(s.accessDuration)

	// Создаем временные данные
	tempData := &model.TemporaryPatientData{
		AccessToken:   accessToken,
		PatientID:     request.PatientID,
		ClinicID:      request.ClinicID,
		MedicalDataID: request.MedicalDataID,
		EncryptedData: encryptedData,
		GrantedAt:     now,
		ExpiresAt:     expiresAt,
		IsRevoked:     false,
	}

	if err := s.repo.CreateTemporaryPatientData(ctx, tempData); err != nil {
		return nil, fmt.Errorf("failed to create temporary data: %w", err)
	}

	// Обновляем статус запроса
	request.Status = "approved"
	request.ApprovedAt = &now
	request.ExpiresAt = &expiresAt

	if err := s.repo.UpdatePatientAccessRequest(ctx, request); err != nil {
		log.Error().Err(err).Msg("failed to update request status")
		// Не возвращаем ошибку, т.к. временные данные уже созданы
	}

	// TODO: Интеграция с блокчейном (опционально)
	// Если blockchain-service доступен, записываем транзакцию
	// Это можно сделать через gRPC вызов к blockchain-service
	// Пока просто логируем для аудита
	log.Info().
		Int64("request_id", requestID).
		Int64("approver_id", approverID).
		Str("access_token", accessToken).
		Time("expires_at", expiresAt).
		Msg("access request approved - ready for blockchain logging")

	// В будущем здесь будет вызов к blockchain-service:
	// blockchainReq := &blockchainpb.CreateAccessRequestRequest{
	//     Id:            fmt.Sprintf("%d", request.ID),
	//     PatientId:     fmt.Sprintf("%d", request.PatientID),
	//     ClinicId:      fmt.Sprintf("%d", request.ClinicID),
	//     MedicalDataId: fmt.Sprintf("%d", request.MedicalDataID),
	//     RequestType:   "patient_access",
	// }
	// blockchainResp, err := s.blockchainClient.CreateAccessRequest(ctx, blockchainReq)
	// if err == nil {
	//     request.BlockchainTxID = blockchainResp.Id
	//     s.repo.UpdatePatientAccessRequest(ctx, request)
	// }

	return tempData, nil
}

func (s *patientAccessService) RejectAccessRequest(ctx context.Context, requestID int64) error {
	request, err := s.repo.GetPatientAccessRequest(ctx, requestID)
	if err != nil {
		return fmt.Errorf("access request not found: %w", err)
	}

	if request.Status != "pending" {
		return fmt.Errorf("request is not in pending status")
	}

	request.Status = "rejected"
	return s.repo.UpdatePatientAccessRequest(ctx, request)
}

func (s *patientAccessService) GetAccessRequest(ctx context.Context, requestID int64) (*model.PatientAccessRequest, error) {
	return s.repo.GetPatientAccessRequest(ctx, requestID)
}

func (s *patientAccessService) ListPatientAccessRequests(ctx context.Context, patientID *int64, clinicID *int64, status *string) ([]*model.PatientAccessRequest, error) {
	return s.repo.ListPatientAccessRequests(ctx, patientID, clinicID, status)
}

func (s *patientAccessService) GetTemporaryData(ctx context.Context, accessToken string) (*model.MedicalData, error) {
	log := zerolog.Ctx(ctx).With().Str("component", "patient_access_service").Str("method", "GetTemporaryData").Logger()

	// Получаем временные данные
	tempData, err := s.repo.GetTemporaryDataByToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("temporary data not found: %w", err)
	}

	// Проверяем срок действия
	if time.Now().After(tempData.ExpiresAt) {
		// Данные истекли, удаляем их
		s.repo.DeleteTemporaryPatientData(ctx, tempData.ID)
		return nil, fmt.Errorf("access token has expired")
	}

	if tempData.IsRevoked {
		return nil, fmt.Errorf("access token has been revoked")
	}

	// Расшифровываем данные
	medData, err := s.decryptMedicalData(tempData.EncryptedData)
	if err != nil {
		log.Error().Err(err).Msg("failed to decrypt data")
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return medData, nil
}

func (s *patientAccessService) RevokeAccess(ctx context.Context, accessToken string) error {
	tempData, err := s.repo.GetTemporaryDataByToken(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("temporary data not found: %w", err)
	}

	tempData.IsRevoked = true
	return s.repo.UpdateTemporaryPatientData(ctx, tempData)
}

func (s *patientAccessService) CleanupExpiredData(ctx context.Context) error {
	log := zerolog.Ctx(ctx).With().Str("component", "patient_access_service").Str("method", "CleanupExpiredData").Logger()

	// Удаляем все истекшие временные данные
	deleted, err := s.repo.DeleteExpiredTemporaryData(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete expired data")
		return err
	}

	// Обновляем статус истекших запросов
	updated, err := s.repo.UpdateExpiredAccessRequests(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to update expired requests")
		return err
	}

	log.Info().
		Int("deleted_temp_data", deleted).
		Int("updated_requests", updated).
		Msg("cleanup completed")

	return nil
}

// Вспомогательные функции для шифрования

func (s *patientAccessService) encryptMedicalData(data *model.MedicalData) (string, error) {
	// Сериализуем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Создаем AES cipher
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	// Создаем GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Генерируем nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Шифруем
	ciphertext := gcm.Seal(nonce, nonce, jsonData, nil)

	// Кодируем в base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *patientAccessService) decryptMedicalData(encryptedData string) (*model.MedicalData, error) {
	// Декодируем из base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	// Создаем AES cipher
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return nil, err
	}

	// Создаем GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Расшифровываем
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	// Десериализуем JSON
	var data model.MedicalData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *patientAccessService) generateAccessToken(patientID, medicalDataID int64) string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("%d:%d:%d", patientID, medicalDataID, timestamp)
	hash := sha256.Sum256([]byte(data + string(s.encryptionKey)))
	return base64.URLEncoding.EncodeToString(hash[:])
}
