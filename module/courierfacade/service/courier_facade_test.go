package service

import (
	"context"
	"fmt"
	"testing"

	cmodels "projects/LDmitryLD/geotask/module/courier/models"
	"projects/LDmitryLD/geotask/module/courier/service"
	smocks "projects/LDmitryLD/geotask/module/courier/service/mocks"
	cfm "projects/LDmitryLD/geotask/module/courierfacade/models"
	omodels "projects/LDmitryLD/geotask/module/order/models"
	omocks "projects/LDmitryLD/geotask/module/order/service/mocks"

	"github.com/stretchr/testify/assert"
)

var (
	testCtx   = context.Background()
	testScore = 1
	testPoint = cmodels.Point{
		Lat: service.DefaultCourierLat,
		Lng: service.DefaultCourierLng,
	}
	testCourier = cmodels.Courier{
		Score:    testScore,
		Location: testPoint,
	}
	testID    = int64(1)
	testOrder = omodels.Order{
		ID: testID,
	}
)

func TestCourierFacer_GetStatus(t *testing.T) {
	courierService := smocks.NewCourierer(t)
	courierService.On("GetCourier", testCtx).Return(&testCourier, nil)

	orderService := omocks.NewOrderer(t)
	orderService.On("GetByRadius", testCtx, testCourier.Location.Lng, testCourier.Location.Lat, float64(CourierVisibilityRadius), unitMeters).Return([]omodels.Order{testOrder}, nil)

	expect := cfm.CourierStatus{
		Courier: testCourier,
		Orders:  []omodels.Order{testOrder},
	}

	facade := NewCourierFacade(courierService, orderService)

	got := facade.GetStatus(testCtx)

	assert.Equal(t, expect, got)
}

func TestCourierFacer_GetStatus_Error(t *testing.T) {
	courierService := smocks.NewCourierer(t)
	courierService.On("GetCourier", testCtx).Return(nil, fmt.Errorf("test error"))

	orderService := omocks.NewOrderer(t)

	facade := NewCourierFacade(courierService, orderService)

	expect := cfm.CourierStatus{}

	got := facade.GetStatus(testCtx)

	assert.Equal(t, expect, got)
}

func TestCourierFacer_GetStatus_Error2(t *testing.T) {
	courierService := smocks.NewCourierer(t)
	courierService.On("GetCourier", testCtx).Return(&testCourier, nil)

	orderService := omocks.NewOrderer(t)
	orderService.On("GetByRadius", testCtx, testCourier.Location.Lng, testCourier.Location.Lat, float64(CourierVisibilityRadius), unitMeters).Return(nil, fmt.Errorf("test error"))

	facade := NewCourierFacade(courierService, orderService)

	expect := cfm.CourierStatus{
		Courier: testCourier,
	}

	got := facade.GetStatus(testCtx)

	assert.Equal(t, expect, got)
}
