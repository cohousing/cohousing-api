package domain

import (
	"fmt"
)

const (
	REL_APARTMENT RelType = "apartment"
)

var (
	ResidentBasePath string
)

type Resident struct {
	BaseModel
	Name         string    `json:"name"`
	PhoneNumber  string    `json:"phone"`
	EmailAddress string    `json:"email"`
	Apartment    Apartment `json:"-"`
	ApartmentID  *uint64   `json:"apartment_id"`
	DefaultHalResource
}

func (r *Resident) Populate(detailed bool) {
	r.AddLink(REL_SELF, fmt.Sprintf("%s/%d", ResidentBasePath, r.ID))

	if detailed && r.ApartmentID != nil {
		r.AddLink(REL_APARTMENT, fmt.Sprintf("%s/%d", ApartmentBasePath, *r.ApartmentID))
	}
}
