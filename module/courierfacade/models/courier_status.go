package models

import (
	cm "projects/LDmitryLD/geotask/module/courier/models"
	om "projects/LDmitryLD/geotask/module/order/models"
)

type CourierStatus struct {
	Courier cm.Courier `json:"courier"`
	Orders  []om.Order `json:"orders"`
}
