package core

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/osamikoyo/yoconf/casher"
	"github.com/osamikoyo/yoconf/logger"
	"github.com/osamikoyo/yoconf/models"
	"github.com/osamikoyo/yoconf/retrier"
	"github.com/osamikoyo/yoconf/storage"
	"go.uber.org/zap"
)

const RetrierCount = 5

var ErrNilInput = errors.New("nil input")

type Core struct {
	casher  *casher.Casher
	storage *storage.Storage
	logger  *logger.Logger

	timeout time.Duration
}

func NewCore(
	casher *casher.Casher,
	storage *storage.Storage,
	logger *logger.Logger,
	timeout time.Duration,
) *Core {
	return &Core{
		casher:  casher,
		storage: storage,
		logger:  logger,
		timeout: timeout,
	}
}

func (c *Core) Close() error {
	return c.casher.Close()
}

func (c *Core) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}

func (c *Core) NewConfig(chunk *models.Chunk) error {
	if chunk == nil {
		return ErrNilInput
	}

	ctx, cancel := c.context()
	defer cancel()

	err := retrier.Try(RetrierCount, func() error {
		return c.storage.CreateNewChunk(ctx, chunk)
	})
	if err != nil {
		c.logger.Error("failed create chunk", zap.Error(err))

		return err
	}

	err = retrier.Try(RetrierCount, func() error {
		return c.casher.CreateChunk(ctx, chunk)
	})
	if err != nil {
		c.logger.Error("failed create chunk in cash", zap.Error(err))

		return err
	}

	return nil
}

func (c *Core) RollOn(project string, version int) error {
	if project == "" || version < 1 {
		return ErrNilInput
	}

	ctx, cancel := c.context()
	defer cancel()

	err := retrier.Try(RetrierCount, func() error {
		return c.storage.RollChunkOn(ctx, project, version)
	})
	if err != nil {
		c.logger.Error("failed roll chunk on", zap.Error(err))

		return err
	}

	data, err := c.casher.GetData(ctx, project)
	if err != nil {
		c.logger.Error("failed roll chunk on", zap.Error(err))
	}

	err = retrier.Try(RetrierCount, func() error {
		return c.casher.DeleteChunk(ctx, project)
	})
	if err != nil {
		c.logger.Error("failed roll chunk on", zap.Error(err))

		return err
	}

	err = retrier.Try(RetrierCount, func() error {
		return c.storage.CreateNewChunk(ctx, &models.Chunk{
			Data:    data,
			Version: version,
			Project: project,
			InUse:   false,
		})
	})
	if err != nil {
		c.logger.Error("failed roll chunk on", zap.Error(err))

		return err
	}

	return nil
}

func (c *Core) GetConfig(project string) (*models.Chunk, error) {
	ctx, cancel := c.context()
	defer cancel()

	data, err := c.casher.GetData(ctx, project)
	if err == nil {
		c.logger.Info("successfully fetched config", zap.String("data", data))

		chunk := models.Chunk{}
		if err = json.Unmarshal([]byte(data), &chunk); err != nil {
			return nil, err
		}

		return &chunk, nil
	}

	chunk, err := c.storage.GetChunk(ctx, project)
	if err != nil {
		c.logger.Error("failed get config", zap.Error(err))
	}

	c.logger.Info("successfully fetched config", zap.Any("chunk", chunk))
	return chunk, nil
}

func (c *Core) DeleteChunk(project string, version int) error {
	ctx, cancel := c.context()
	defer cancel()

	if err := c.storage.DeleteConfig(ctx, project, version); err != nil {
		c.logger.Error("failed delete config", zap.Error(err))

		return err
	}

	return nil
}
