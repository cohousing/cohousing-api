package tenant

import "github.com/cohousing/cohousing-tenant-api/domain"

const (
	REL_RESIDENTS domain.RelType = "residents"
)

type Apartment struct {
	domain.BaseModel
	Address string `json:"address"`
	domain.DefaultHalResource
}
