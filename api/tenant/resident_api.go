package tenant

import (
	"fmt"
	"github.com/cohousing/cohousing-api/api/utils"
	"github.com/cohousing/cohousing-api/db"
	"github.com/cohousing/cohousing-api/domain"
	"github.com/cohousing/cohousing-api/domain/tenant"
	"github.com/gin-gonic/gin"
)

var (
	ResidentBasePath string
)

func CreateResidentRoutes(router *gin.RouterGroup, dbFactory db.DBFactory) {
	endpoint := ConfigureBasicTenantEndpoint(router, utils.BasicEndpointConfig{
		Path:           "residents",
		Domain:         tenant.Resident{},
		LinkFactory:    residentLinkFactory,
		DBFactory:      dbFactory,
		RouterHandlers: []gin.HandlerFunc{utils.MustBeTenant(), MustAuthenticate()},
	})
	ResidentBasePath = endpoint.BasePath()
}

func residentLinkFactory(halResource domain.HalResource, basePath string, detailed bool) {
	resident := halResource.(*tenant.Resident)
	resident.AddLink(domain.REL_SELF, fmt.Sprintf("%s/%d", basePath, resident.ID))

	if detailed && resident.ApartmentID != nil {
		resident.AddLink(domain.REL_UPDATE, fmt.Sprintf("%s/%d", basePath, resident.ID))
		resident.AddLink(domain.REL_DELETE, fmt.Sprintf("%s/%d", basePath, resident.ID))
		resident.AddLink(tenant.REL_APARTMENT, fmt.Sprintf("%s/%d", ApartmentBasePath, *resident.ApartmentID))
	}
}
