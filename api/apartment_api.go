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
	ApartmentBasePath string
)

func CreateApartmentRoutes(router *gin.RouterGroup, dbFactory db.DBFactory) {
	endpoint := ConfigureBasicTenantEndpoint(router, utils.BasicEndpointConfig{
		Path:           "apartments",
		Domain:         domain.Apartment{},
		DBFactory:      dbFactory,
		RouterHandlers: []gin.HandlerFunc{MustBeTenant(), MustAuthenticate()},
	})
	ApartmentBasePath = endpoint.BasePath()
	utils.AddLinkFactory(&domain.Apartment{}, apartmentLinkFactory)
}

func apartmentLinkFactory(c *gin.Context, halResource domain2.HalResource, basePath string, detailed bool) {
	apartment := halResource.(*domain.Apartment)
	halResource.AddLink(domain2.REL_SELF, fmt.Sprintf("%s/%d", basePath, apartment.ID))

	if detailed {
		permission := ResolvePermission(c)
		if permission.GlobalAdmin || permission.UpdateApartments {
			halResource.AddLink(domain2.REL_UPDATE, fmt.Sprintf("%s/%d", basePath, apartment.ID))
		}
		if permission.GlobalAdmin || permission.DeleteApartments {
			halResource.AddLink(domain2.REL_DELETE, fmt.Sprintf("%s/%d", basePath, apartment.ID))
		}
		if permission.GlobalAdmin || permission.ReadResidents {
			halResource.AddLink(domain.REL_RESIDENTS, fmt.Sprintf("%s?apartment_id=%d", ResidentBasePath, apartment.ID))
		}
	}
}
