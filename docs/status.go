package docs

import "projects/LDmitryLD/geotask/module/courierfacade/models"

// swagger:route GET /api/status StatusRequest
// Получение статуса сервиса.
// responses:
//  - 200 : StatusResponse

// swagger:response StatusResponse
type StatusResponse struct {
	// in:body
	Body models.CourierStatus
}
