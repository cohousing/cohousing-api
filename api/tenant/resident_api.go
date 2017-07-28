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
		DBFactory:      dbFactory,
		RouterHandlers: []gin.HandlerFunc{utils.MustBeTenant(), MustAuthenticate()},
	})
	ResidentBasePath = endpoint.BasePath()
	utils.AddLinkFactory(&tenant.Resident{}, residentLinkFactory)
}

func residentLinkFactory(c *gin.Context, halResource domain.HalResource, basePath string, detailed bool) {
	resident := halResource.(*tenant.Resident)
	resident.AddLink(domain.REL_SELF, fmt.Sprintf("%s/%d", basePath, resident.ID))

	if detailed {
		permission := ResolvePermission(c)
		if permission.GlobalAdmin || permission.UpdateResidents {
			resident.AddLink(domain.REL_UPDATE, fmt.Sprintf("%s/%d", basePath, resident.ID))
		}
		if permission.GlobalAdmin || permission.DeleteResidents {
			resident.AddLink(domain.REL_DELETE, fmt.Sprintf("%s/%d", basePath, resident.ID))
		}

		if permission.GlobalAdmin || resident.ApartmentID != nil && permission.ReadApartments {
			resident.AddLink(tenant.REL_APARTMENT, fmt.Sprintf("%s/%d", ApartmentBasePath, *resident.ApartmentID))
		}
	}
}
