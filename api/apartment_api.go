package api

import (
	"fmt"
	"github.com/cohousing/cohousing-api/domain"
	"github.com/gin-gonic/gin"
)

var (
	ApartmentBasePath string
)

func CreateApartmentRoutes(router *gin.RouterGroup) {
	endpoint := ConfigureBasicTenantEndpoint(router, "apartments", domain.Apartment{}, apartmentLinkFactory)
	ApartmentBasePath = endpoint.BasePath()
}

func apartmentLinkFactory(halResource domain.HalResource, basePath string, detailed bool) {
	apartment := halResource.(*domain.Apartment)
	halResource.AddLink(domain.REL_SELF, fmt.Sprintf("%s/%d", basePath, apartment.ID))

	if detailed {
		halResource.AddLink(domain.REL_RESIDENTS, fmt.Sprintf("%s?apartment_id=%d", ResidentBasePath, apartment.ID))
	}
}
