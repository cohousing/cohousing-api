package domain

import "github.com/cohousing/cohousing-api-utils/domain"

const (
	REL_APARTMENT  domain.RelType = "apartment"
	REL_APARTMENTS domain.RelType = "apartments"
)

type Resident struct {
	domain.BaseModel
	Name         string    `json:"name"`
	PhoneNumber  string    `json:"phone"`
	EmailAddress string    `json:"email"`
	Apartment    Apartment `json:"-"`
	ApartmentID  *uint64   `json:"-"`
	domain.DefaultHalResource
}
