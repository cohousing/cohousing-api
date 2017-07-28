package admin

import (
	"github.com/cohousing/cohousing-tenant-api/domain"
)

const (
	REL_TENANTS domain.RelType = "tenants"
)

type Tenant struct {
	domain.BaseModel
	Context   string `gorm:"size:100"`
	Name      string
	TenantUrl string `gorm:"-"`
	CustomUrl string
	domain.DefaultHalResource
}

func (tenant *Tenant) GetUrl() string {
	if tenant.CustomUrl != "" {
		return tenant.CustomUrl
	} else {
		return tenant.TenantUrl
	}
}
