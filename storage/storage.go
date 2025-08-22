package storage

import (
	"fmt"

	"github.com/osamikoyo/yoconf/config"
	"github.com/osamikoyo/yoconf/logger"
	"github.com/osamikoyo/yoconf/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Storage struct {
	logger *logger.Logger
	cfg    *config.Config
	db     *gorm.DB
}

func NewStorage(db *gorm.DB, logger *logger.Logger, cfg *config.Config) *Storage {
	return &Storage{
		logger: logger,
		cfg:    cfg,
		db:     db,
	}
}

func (s *Storage) CreateNewChunk(chunk *models.Chunk) error {
	res := s.db.Where(&models.Chunk{
		InUse:   true,
		Project: chunk.Project,
	}).Updates(&models.Chunk{InUse: false})
	if err := res.Error; err != nil {
		s.logger.Error("failed update old chunk",
			zap.Any("chunk", chunk),
			zap.Error(err))

		return fmt.Errorf("failed update old chunk: %v", err)
	}

	res = s.db.Create(chunk)
	if err := res.Error; err != nil {
		s.logger.Error("failed create new chunk",
			zap.Any("chunk", chunk),
			zap.Error(err))

		return fmt.Errorf("failed create new chunk: %v", err)
	}

	s.logger.Info("successfully create new chunk", zap.Any("chunk", chunk))

	return nil
}
