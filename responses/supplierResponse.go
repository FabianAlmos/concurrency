package response

import model "concuLec/models"

type SupplierResponse struct {
	Id           int                `json:"id"`
	Image        string             `json:"image"`
	Name         string             `json:"name"`
	Type         string             `json:"type"`
	WorkingHours model.WorkingHours `json:"workingHours"`
}
