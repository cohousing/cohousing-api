package domain

import "github.com/cohousing/cohousing-api-utils/domain"

const (
	REL_RESIDENTS domain.RelType = "residents"
)

type Apartment struct {
	domain.BaseModel
	Address string `json:"address"`
	domain.DefaultHalResource
}
