package domain

const (
	REL_APARTMENT RelType = "apartment"
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
