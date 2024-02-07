package service

import (
	"context"
	"fmt"
	"math"
	"projects/LDmitryLD/geotask/geo"
	"projects/LDmitryLD/geotask/module/courier/models"
	"projects/LDmitryLD/geotask/module/courier/storage/mocks"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var (
	testAllowedZone    = geo.NewAllowedZone()
	testDisalowedZone1 = geo.NewDisAllowedZone1()
	testDisalowedZone2 = geo.NewDisAllowedZone2()
	testDisalowedZones = []geo.PolygonChecker{testDisalowedZone1, testDisalowedZone2}
	testCtx            = context.Background()
	testScore          = 1
	testBadPoint       = models.Point{Lat: 11.1, Lng: 22.2}
	testGoodPoint      = models.Point{Lat: DefaultCourierLat, Lng: DefaultCourierLng}
	testBadCourier     = models.Courier{
		Score:    testScore,
		Location: testBadPoint,
	}
	testGoodCourier = models.Courier{
		Score:    testScore,
		Location: testGoodPoint,
	}
	testZoom      = 14
	testPrecision = 0.001 / math.Pow(2, float64(testZoom-14))
)

func TestCourierService_GetCourier(t *testing.T) {
	mockStorage := mocks.NewCourierStorager(t)
	mockStorage.On("GetOne", testCtx).Return(&testGoodCourier, nil)
	mockStorage.On("Save", testCtx, testGoodCourier).Return(nil)
	service := NewCourierService(mockStorage, testAllowedZone, testDisalowedZones)

	got, err := service.GetCourier(testCtx)

	assert.Nil(t, err)
	assert.Equal(t, got, &testGoodCourier)
}

func TestCourierService_GetCourier_RedisNil(t *testing.T) {
	mockStorage := mocks.NewCourierStorager(t)
	mockStorage.On("GetOne", testCtx).Return(nil, redis.Nil)
	testGoodCourier.Score = 0
	defer func() {
		testGoodCourier.Score = testScore
	}()
	mockStorage.On("Save", testCtx, testGoodCourier).Return(nil)
	service := NewCourierService(mockStorage, testAllowedZone, testDisalowedZones)

	got, err := service.GetCourier(testCtx)

	assert.Nil(t, err)
	assert.Equal(t, got, &testGoodCourier)
}

func TestCourierService_GetCourier_Error(t *testing.T) {
	mockStorage := mocks.NewCourierStorager(t)
	mockStorage.On("GetOne", testCtx).Return(nil, fmt.Errorf("test error"))

	service := NewCourierService(mockStorage, testAllowedZone, testDisalowedZones)

	_, err := service.GetCourier(testCtx)

	assert.NotNil(t, err)
}

func TestCourierService_GetCourier_Error2(t *testing.T) {
	mockStorage := mocks.NewCourierStorager(t)
	mockStorage.On("GetOne", testCtx).Return(&testGoodCourier, nil)
	mockStorage.On("Save", testCtx, testGoodCourier).Return(fmt.Errorf("test error"))
	service := NewCourierService(mockStorage, testAllowedZone, testDisalowedZones)

	_, err := service.GetCourier(testCtx)

	assert.NotNil(t, err)
}

func TestCourierService_MovCourier(t *testing.T) {
	mockStorage := mocks.NewCourierStorager(t)
	service := NewCourierService(mockStorage, testAllowedZone, testDisalowedZones)

	testCases := []struct {
		direction int
		expectLat float64
		expectLng float64
	}{
		{
			direction: DirectionDown,
			expectLat: DefaultCourierLat - testPrecision,
			expectLng: DefaultCourierLng,
		},
		{
			direction: DirectionUp,
			expectLat: DefaultCourierLat + testPrecision,
			expectLng: DefaultCourierLng,
		},
		{
			direction: DirectionLeft,
			expectLat: DefaultCourierLat,
			expectLng: DefaultCourierLng - testPrecision,
		},
		{
			direction: DirectionRight,
			expectLat: DefaultCourierLat,
			expectLng: DefaultCourierLng + testPrecision,
		},
	}

	for _, tt := range testCases {
		expect := models.Courier{
			Score: testScore,
			Location: models.Point{
				Lat: tt.expectLat,
				Lng: tt.expectLng,
			},
		}
		mockStorage.On("Save", testCtx, expect).Return(nil)
		got, err := service.MoveCourier(testGoodCourier, tt.direction, testZoom)

		assert.Nil(t, err)
		assert.Equal(t, expect, got)
	}
}

func TestCourierService_MovCourier_Save(t *testing.T) {
	mockStorage := mocks.NewCourierStorager(t)
	mockStorage.On("Save", testCtx, testGoodCourier).Return(nil)
	service := NewCourierService(mockStorage, testAllowedZone, testDisalowedZones)

	err := service.Save(testCtx, testGoodCourier)

	assert.Nil(t, err)
}
