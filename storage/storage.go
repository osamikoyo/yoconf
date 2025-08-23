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

func (s *Storage) GetChunk(project string) (*models.Chunk, error) {
	var chunk models.Chunk

	res := s.db.Where(&models.Chunk{
		InUse:   true,
		Project: project,
	}).First(&chunk)
	if err := res.Error; err != nil {
		s.logger.Error("failed fetch chunk",
			zap.String("project", project),
			zap.Error(err))

		return nil, fmt.Errorf("failed fetch chunk: %v", err)
	}

	return &chunk, nil
}

func (s *Storage) ListVersions(project string) ([]int, error) {
	var versions []models.Chunk

	res := s.db.Where(&models.Chunk{
		Project: project,
	}).Find(&versions)
	if err := res.Error; err != nil {
		s.logger.Error("failed find chunks",
			zap.String("project", project),
			zap.Error(err))

		return nil, fmt.Errorf("failed find chunks: %v", err)
	}

	resp := make([]int, len(versions))
	for i, v := range versions {
		resp[i] = v.Version
	}

	s.logger.Info("fetched versions",
		zap.String("project", project))

	return resp, nil
}

func (s *Storage) RollChunkOn(project string, version int) error {
	res := s.db.Where(&models.Chunk{
		Project: project,
		InUse:   true,
	}).Update("in_use", false)
	if err := res.Error; err != nil {
		s.logger.Error("failed to update old version",
			zap.String("project", project),
			zap.Int("version", version),
			zap.Error(err))

		return fmt.Errorf("failed to update old version: %v", err)
	}

	res = s.db.Where(&models.Chunk{
		Project: project,
		Version: version,
	}).Update("in_use", true)
	if err := res.Error; err != nil {
		s.logger.Error("failed roll chunk on",
			zap.String("project", project),
			zap.Int("version", version),
			zap.Error(err))

		return fmt.Errorf("failed to update new version: %v", err)
	}

	s.logger.Info("successfully roll chunk on",
		zap.String("project", project),
		zap.Int("version", version))

	return nil
}

func (s *Storage) DeleteConfig(project string, version int) error {
	res := s.db.Where(&models.Chunk{
		Project: project,
		Version: version,
	}).Delete(&models.Chunk{})
	if err := res.Error; err != nil {
		s.logger.Error("failed delete chunk",
			zap.String("project", project),
			zap.Int("version", version),
			zap.Error(err))

		return fmt.Errorf("failed delete config: %v", err)
	}

	s.logger.Info("chunk deleted successfully",
		zap.String("project", project),
		zap.Int("version", version))

	return nil
}
