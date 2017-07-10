package domain

type Resident struct {
	BaseModel
	Name         string `json:"name"`
	PhoneNumber  string `json:"phone"`
	EmailAddress string `json:"email"`
}
