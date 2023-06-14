package model

type Menu struct {
	Id          int      `json:"id"`
	Image       string   `json:"image"`
	Ingredients []string `json:"ingredients"`
	Name        string   `json:"name"`
	Price       float32  `json:"price"`
	Type        string   `json:"type"`
}
