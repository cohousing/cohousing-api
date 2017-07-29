package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/cohousing/cohousing-api-utils/api"
	"github.com/cohousing/cohousing-api-utils/db"
	"github.com/cohousing/cohousing-tenant-api/domain"
	domain2 "github.com/cohousing/cohousing-api-utils/domain"
)

var (
	ResidentBasePath string
)

func CreateResidentRoutes(router *gin.RouterGroup, dbFactory db.DBFactory) {
	endpoint := ConfigureBasicTenantEndpoint(router, utils.BasicEndpointConfig{
		Path:           "residents",
		Domain:         domain.Resident{},
		DBFactory:      dbFactory,
		RouterHandlers: []gin.HandlerFunc{MustBeTenant(), MustAuthenticate()},
	})
	ResidentBasePath = endpoint.BasePath()
	utils.AddLinkFactory(&domain.Resident{}, residentLinkFactory)
}

func residentLinkFactory(c *gin.Context, halResource domain2.HalResource, basePath string, detailed bool) {
	resident := halResource.(*domain.Resident)
	resident.AddLink(domain2.REL_SELF, fmt.Sprintf("%s/%d", basePath, resident.ID))

	if detailed {
		permission := ResolvePermission(c)
		if permission.GlobalAdmin || permission.UpdateResidents {
			resident.AddLink(domain2.REL_UPDATE, fmt.Sprintf("%s/%d", basePath, resident.ID))
		}
		if permission.GlobalAdmin || permission.DeleteResidents {
			resident.AddLink(domain2.REL_DELETE, fmt.Sprintf("%s/%d", basePath, resident.ID))
		}

		if permission.GlobalAdmin || resident.ApartmentID != nil && permission.ReadApartments {
			resident.AddLink(domain.REL_APARTMENT, fmt.Sprintf("%s/%d", ApartmentBasePath, *resident.ApartmentID))
		}
	}
}
