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

func getKey(project string, version int, use bool) string {
	return fmt.Sprintf("%s:%d:%t", project, version, use)
}

func (c *Casher) CreateChunk(ctx context.Context, chunk *models.Chunk) error {
	key := getKey(chunk.Project, chunk.Version, chunk.InUse)

	_, err := c.client.Set(ctx, key, chunk.Data, ExpTime).Result()
	if err != nil {
		c.logger.Error("failed set",
			zap.String("key", key),
			zap.Error(err))

		return fmt.Errorf("failed set: %v", err)
	}

	c.logger.Info("successfully create chunk",
		zap.Any("chunk", chunk))

	return nil
}

func (c *Casher) GetData(ctx context.Context, project string, version int) (string, error) {
	key := getKey(project, version, true)

	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		c.logger.Error("failed get data",
			zap.String("key", key),
			zap.Error(err))

		return "", err
	}

	c.logger.Info("successfully fetch data",
		zap.String("key", "key"))

	return data, nil
}

func (c *Casher) DeleteChunk(ctx context.Context, project string, version int, use bool) error {
	key := getKey(project, version, use)

	_, err := c.client.Del(ctx, key).Result()
	if err != nil {
		c.logger.Error("failed delete",
			zap.String("key", key),
			zap.Error(err))

		return fmt.Errorf("failed delete: %v", err)
	}

	c.logger.Info("successfully delete chunk",
		zap.String("key", key))

	return nil
}
