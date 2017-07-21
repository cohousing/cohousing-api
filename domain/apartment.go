package domain

const (
	REL_RESIDENTS RelType = "residents"
)

type Apartment struct {
	BaseModel
	Address string `json:"address"`
	DefaultHalResource
}
