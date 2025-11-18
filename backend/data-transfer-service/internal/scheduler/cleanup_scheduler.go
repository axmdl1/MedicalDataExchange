package scheduler

import (
	"context"
	"time"

	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/service"
	"github.com/rs/zerolog/log"
)

type CleanupScheduler struct {
	patientAccessService service.PatientAccessService
	interval             time.Duration
	stopChan             chan struct{}
}

func NewCleanupScheduler(patientAccessService service.PatientAccessService) *CleanupScheduler {
	return &CleanupScheduler{
		patientAccessService: patientAccessService,
		interval:             5 * time.Minute, // Запускаем каждые 5 минут
		stopChan:             make(chan struct{}),
	}
}

func (s *CleanupScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Info().
		Dur("interval", s.interval).
		Msg("Cleanup scheduler started")

	// Запускаем сразу при старте
	s.runCleanup(ctx)

	for {
		select {
		case <-ticker.C:
			s.runCleanup(ctx)
		case <-s.stopChan:
			log.Info().Msg("Cleanup scheduler stopped")
			return
		case <-ctx.Done():
			log.Info().Msg("Cleanup scheduler stopped by context")
			return
		}
	}
}

func (s *CleanupScheduler) Stop() {
	close(s.stopChan)
}

func (s *CleanupScheduler) runCleanup(ctx context.Context) {
	log.Debug().Msg("Running cleanup of expired temporary data")

	if err := s.patientAccessService.CleanupExpiredData(ctx); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to cleanup expired data")
	} else {
		log.Debug().Msg("Cleanup completed successfully")
	}
}
