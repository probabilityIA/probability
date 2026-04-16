package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
)

type DriverOptionResponse struct {
	ID             uint   `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Phone          string `json:"phone"`
	Identification string `json:"identification"`
	Status         string `json:"status"`
	LicenseType    string `json:"license_type"`
}

type VehicleOptionResponse struct {
	ID           uint   `json:"id"`
	Type         string `json:"type"`
	LicensePlate string `json:"license_plate"`
	Brand        string `json:"brand"`
	VehicleModel string `json:"vehicle_model"`
	Status       string `json:"status"`
}

type AssignableOrderResponse struct {
	ID            string    `json:"id"`
	OrderNumber   string    `json:"order_number"`
	CustomerName  string    `json:"customer_name"`
	CustomerPhone string    `json:"customer_phone"`
	Address       string    `json:"address"`
	City          string    `json:"city"`
	Lat           *float64  `json:"lat"`
	Lng           *float64  `json:"lng"`
	TotalAmount   float64   `json:"total_amount"`
	ItemCount     int       `json:"item_count"`
	CreatedAt     time.Time `json:"created_at"`
}

func DriversFromDTOs(drivers []dtos.DriverOption) []DriverOptionResponse {
	result := make([]DriverOptionResponse, len(drivers))
	for i, d := range drivers {
		result[i] = DriverOptionResponse{
			ID:             d.ID,
			FirstName:      d.FirstName,
			LastName:       d.LastName,
			Phone:          d.Phone,
			Identification: d.Identification,
			Status:         d.Status,
			LicenseType:    d.LicenseType,
		}
	}
	return result
}

func VehiclesFromDTOs(vehicles []dtos.VehicleOption) []VehicleOptionResponse {
	result := make([]VehicleOptionResponse, len(vehicles))
	for i, v := range vehicles {
		result[i] = VehicleOptionResponse{
			ID:           v.ID,
			Type:         v.Type,
			LicensePlate: v.LicensePlate,
			Brand:        v.Brand,
			VehicleModel: v.VehicleModel,
			Status:       v.Status,
		}
	}
	return result
}

func AssignableOrdersFromDTOs(orders []dtos.AssignableOrder) []AssignableOrderResponse {
	result := make([]AssignableOrderResponse, len(orders))
	for i, o := range orders {
		result[i] = AssignableOrderResponse{
			ID:            o.ID,
			OrderNumber:   o.OrderNumber,
			CustomerName:  o.CustomerName,
			CustomerPhone: o.CustomerPhone,
			Address:       o.Address,
			City:          o.City,
			Lat:           o.Lat,
			Lng:           o.Lng,
			TotalAmount:   o.TotalAmount,
			ItemCount:     o.ItemCount,
			CreatedAt:     o.CreatedAt,
		}
	}
	return result
}
