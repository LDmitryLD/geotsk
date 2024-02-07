package service

import (
	"context"
	"fmt"
	"projects/LDmitryLD/geotask/geo"
	"projects/LDmitryLD/geotask/module/courier/service"
	"projects/LDmitryLD/geotask/module/order/models"
	"projects/LDmitryLD/geotask/module/order/storage/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testAllowedZone     = geo.NewAllowedZone()
	testDisAllowedZone1 = geo.NewDisAllowedZone1()
	testDisAllowedZone2 = geo.NewDisAllowedZone2()
	testDisAllowedZones = []geo.PolygonChecker{testDisAllowedZone1, testDisAllowedZone2}
	testCtx             = context.Background()
	testLng             = service.DefaultCourierLng
	testLat             = service.DefaultCourierLat
	testRadius          = 1000.0
	testUnit            = "km"
	testID              = int64(1)
	testOrder           = models.Order{
		ID:  testID,
		Lng: testLng,
		Lat: testLat,
	}
	testCount = 10
)

func TestOrderService_GetByRadius(t *testing.T) {
	expect := []models.Order{testOrder}

	storageMock := mocks.NewOrderStorager(t)
	storageMock.On("GetByRadius", testCtx, testLng, testLat, testRadius, testUnit).Return(expect, nil)
	service := NewOrderService(storageMock, testAllowedZone, testDisAllowedZones)

	got, err := service.GetByRadius(testCtx, testLng, testLat, testRadius, testUnit)

	assert.Nil(t, err)
	assert.Equal(t, expect, got)
}

func TestOrderService_Save(t *testing.T) {
	storageMock := mocks.NewOrderStorager(t)
	storageMock.On("Save", testCtx, testOrder, orderMaxAge).Return(nil)

	service := NewOrderService(storageMock, testAllowedZone, testDisAllowedZones)

	err := service.Save(testCtx, testOrder)

	assert.Nil(t, err)
}

func TestOrderService_GetCount(t *testing.T) {
	storageMock := mocks.NewOrderStorager(t)
	storageMock.On("GetCount", testCtx).Return(testCount, nil)

	service := NewOrderService(storageMock, testAllowedZone, testDisAllowedZones)

	got, err := service.GetCount(testCtx)

	assert.Nil(t, err)
	assert.Equal(t, testCount, got)
}

func TestOrderService_RemoveOldOrders(t *testing.T) {
	storageMock := mocks.NewOrderStorager(t)
	storageMock.On("RemoveOldOrders", testCtx, orderMaxAge).Return(nil)

	service := NewOrderService(storageMock, testAllowedZone, testDisAllowedZones)

	err := service.RemoveOldOrders(testCtx)

	assert.Nil(t, err)
}

func TestOrderService_GenerateOrder(t *testing.T) {
	storageMock := mocks.NewOrderStorager(t)
	storageMock.On("GenerateUniqueID", testCtx).Return(testID, nil)
	storageMock.On("Save", testCtx, mock.Anything, orderMaxAge).Return(nil)

	service := NewOrderService(storageMock, testAllowedZone, testDisAllowedZones)

	err := service.GenerateOrder(testCtx)

	assert.Nil(t, err)
}

func TestOrderService_GenerateOrder_Error(t *testing.T) {
	storageMock := mocks.NewOrderStorager(t)
	storageMock.On("GenerateUniqueID", testCtx).Return(int64(0), fmt.Errorf("test error"))

	service := NewOrderService(storageMock, testAllowedZone, testDisAllowedZones)

	err := service.GenerateOrder(testCtx)

	assert.NotNil(t, err)
}

func TestOrderService_RemoveOrder(t *testing.T) {
	storageMock := mocks.NewOrderStorager(t)
	storageMock.On("RemoveOrder", testCtx, testOrder).Return(nil)

	service := NewOrderService(storageMock, testAllowedZone, testDisAllowedZones)

	err := service.RemoveOrder(testCtx, testOrder)

	assert.Nil(t, err)
}
