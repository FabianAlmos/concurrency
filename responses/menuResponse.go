package response

import model "concuLec/models"

type MenuResponse struct {
	Menu []model.Menu `json:"menu"`
}
