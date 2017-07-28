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
	ApartmentBasePath string
)

func CreateApartmentRoutes(router *gin.RouterGroup, dbFactory db.DBFactory) {
	endpoint := ConfigureBasicTenantEndpoint(router, utils.BasicEndpointConfig{
		Path:           "apartments",
		Domain:         tenant.Apartment{},
		DBFactory:      dbFactory,
		RouterHandlers: []gin.HandlerFunc{utils.MustBeTenant(), MustAuthenticate()},
	})
	ApartmentBasePath = endpoint.BasePath()
	utils.AddLinkFactory(&tenant.Apartment{}, apartmentLinkFactory)
}

func apartmentLinkFactory(c *gin.Context, halResource domain.HalResource, basePath string, detailed bool) {
	apartment := halResource.(*tenant.Apartment)
	halResource.AddLink(domain.REL_SELF, fmt.Sprintf("%s/%d", basePath, apartment.ID))

	if detailed {
		permission := ResolvePermission(c)
		if permission.GlobalAdmin || permission.UpdateApartments {
			halResource.AddLink(domain.REL_UPDATE, fmt.Sprintf("%s/%d", basePath, apartment.ID))
		}
		if permission.GlobalAdmin || permission.DeleteApartments {
			halResource.AddLink(domain.REL_DELETE, fmt.Sprintf("%s/%d", basePath, apartment.ID))
		}
		if permission.GlobalAdmin || permission.ReadResidents {
			halResource.AddLink(tenant.REL_RESIDENTS, fmt.Sprintf("%s?apartment_id=%d", ResidentBasePath, apartment.ID))
		}
	}
}
