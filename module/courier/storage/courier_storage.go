package storage

import (
	"context"
	"encoding/json"
	"log"
	"projects/LDmitryLD/geotask/module/courier/models"

	"github.com/redis/go-redis/v9"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=CourierStorager
type CourierStorager interface {
	Save(ctx context.Context, courier models.Courier) error
	GetOne(ctx context.Context) (*models.Courier, error)
}

type CourierStorage struct {
	storage *redis.Client
}

func NewCourierStorage(storage *redis.Client) CourierStorager {
	return &CourierStorage{storage: storage}
}

func (c *CourierStorage) Save(ctx context.Context, courier models.Courier) error {
	courieRaw, err := json.Marshal(courier)
	if err != nil {
		log.Println("ошибка при кодированни структуры curier:", err)
		return err
	}

	key := "courier"
	err = c.storage.Set(ctx, key, string(courieRaw), 0).Err()
	if err != nil {
		log.Println("ошибка при сохранении курьера в Redis:", err)
		return err
	}

	return nil
}

func (c *CourierStorage) GetOne(ctx context.Context) (*models.Courier, error) {
	var courier models.Courier

	key := "courier"
	courierRaw, err := c.storage.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(courierRaw), &courier)
	if err != nil {
		return nil, err
	}

	return &courier, nil
}
