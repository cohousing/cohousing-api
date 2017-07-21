package domain

type Apartment struct {
	BaseModel
	Address string `json:"address"`
	DefaultHalResource
}
