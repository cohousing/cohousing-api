package domain

import (
	"fmt"
)

var (
	ApartmentBasePath string
)

type Apartment struct {
	BaseModel
	Address string `json:"address"`
	DefaultHalResource
}

func (r *Apartment) Populate(detailed bool) {
	r.AddLink(REL_SELF, fmt.Sprintf("%s/%d", ApartmentBasePath, r.ID))
}
