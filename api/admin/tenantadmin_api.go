package admin

import (
	"fmt"
	"github.com/cohousing/cohousing-api/api/utils"
	"github.com/cohousing/cohousing-api/db"
	"github.com/cohousing/cohousing-api/domain"
	"github.com/cohousing/cohousing-api/domain/admin"
	"github.com/gin-gonic/gin"
)

var (
	TenantAdminBasePath string
)

func CreateTenantAdminRoutes(router *gin.RouterGroup, dbFactory db.DBFactory) {
	endpoint := ConfigureBasicAdminEndpoint(router, utils.BasicEndpointConfig{
		Path:           "tenants",
		Domain:         admin.Tenant{},
		DBFactory:      dbFactory,
		RouterHandlers: []gin.HandlerFunc{utils.MustBeAdminDomain()},
	})
	TenantAdminBasePath = endpoint.BasePath()
	utils.AddLinkFactory(&admin.Tenant{}, tenantAdminLinkFactory)
}

func tenantAdminLinkFactory(c *gin.Context, halResource domain.HalResource, basePath string, detailed bool) {
	tenant := halResource.(*admin.Tenant)
	halResource.AddLink(domain.REL_SELF, fmt.Sprintf("%s/%d", basePath, tenant.ID))
}
