package domain

import "time"

type BaseModel struct {
	ID        uint64     `gorm:"primary_key" json:"id"`
	CreatedAt *time.Time `json:"created,omitempty"`
	UpdatedAt *time.Time `json:"updated,omitempty"`
	DeletedAt *time.Time `sql:"index" json:"deleted,omitempty"`
}
