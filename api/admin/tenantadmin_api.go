package admin

import (
	"fmt"
	"github.com/cohousing/cohousing-api/db"
	"github.com/cohousing/cohousing-api/domain"
	"github.com/cohousing/cohousing-api/domain/admin"
	"github.com/gin-gonic/gin"
)

var (
	TenantAdminBasePath string
)

func CreateTenantAdminRoutes(router *gin.RouterGroup, dbFactory db.DBFactory) {
	endpoint := ConfigureBasicAdminEndpoint(router, "tenants", admin.Tenant{}, tenantAdminLinkFactory, dbFactory)
	TenantAdminBasePath = endpoint.BasePath()
}

func tenantAdminLinkFactory(halResource domain.HalResource, basePath string, detailed bool) {
	tenant := halResource.(*admin.Tenant)
	halResource.AddLink(domain.REL_SELF, fmt.Sprintf("%s/%d", basePath, tenant.ID))
}
