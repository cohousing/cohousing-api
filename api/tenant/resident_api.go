package tenant

import (
	"fmt"
	"github.com/cohousing/cohousing-api/db"
	"github.com/cohousing/cohousing-api/domain"
	"github.com/cohousing/cohousing-api/domain/tenant"
	"github.com/gin-gonic/gin"
)

var (
	ResidentBasePath string
)

func CreateResidentRoutes(router *gin.RouterGroup, dbFactory db.DBFactory) {
	endpoint := ConfigureBasicTenantEndpoint(router, "residents", tenant.Resident{}, residentLinkFactory, dbFactory)
	ResidentBasePath = endpoint.BasePath()
}

func residentLinkFactory(halResource domain.HalResource, basePath string, detailed bool) {
	resident := halResource.(*tenant.Resident)
	resident.AddLink(domain.REL_SELF, fmt.Sprintf("%s/%d", basePath, resident.ID))

	if detailed && resident.ApartmentID != nil {
		resident.AddLink(tenant.REL_APARTMENT, fmt.Sprintf("%s/%d", ApartmentBasePath, *resident.ApartmentID))
	}
}
