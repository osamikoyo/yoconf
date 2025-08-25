package casher

import (
	"context"
	"fmt"
	"time"

	"github.com/osamikoyo/yoconf/logger"
	"github.com/osamikoyo/yoconf/models"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Casher struct {
	client *redis.Client
	logger *logger.Logger
}

var ExpTime = 2 * time.Hour

func NewCasher(client *redis.Client, logger *logger.Logger) *Casher {
	return &Casher{
		client: client,
		logger: logger,
	}
}

func (c *Casher) Close() error {
	return c.client.Close()
}

func (c *Casher) CreateChunk(ctx context.Context, chunk *models.Chunk) error {
	_, err := c.client.Set(ctx, chunk.Project, chunk.Data, ExpTime).Result()
	if err != nil {
		c.logger.Error("failed set",
			zap.String("key", chunk.Project),
			zap.Error(err))

		return fmt.Errorf("failed set: %v", err)
	}

	c.logger.Info("successfully create chunk",
		zap.Any("chunk", chunk))

	return nil
}

func (c *Casher) GetData(ctx context.Context, project string) (string, error) {
	data, err := c.client.Get(ctx, project).Result()
	if err != nil {
		c.logger.Error("failed get data",
			zap.String("key", project),
			zap.Error(err))

		return "", err
	}

	c.logger.Info("successfully fetch data",
		zap.String("key", project))

	return data, nil
}

func (c *Casher) DeleteChunk(ctx context.Context, project string) error {
	_, err := c.client.Del(ctx, project).Result()
	if err != nil {
		c.logger.Error("failed delete",
			zap.String("key", project),
			zap.Error(err))

		return fmt.Errorf("failed delete: %v", err)
	}

	c.logger.Info("successfully delete chunk",
		zap.String("key", project))

	return nil
}
