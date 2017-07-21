package api

import (
	"fmt"
	"github.com/cohousing/cohousing-api/domain"
	"github.com/gin-gonic/gin"
)

var (
	ResidentBasePath string
)

func CreateResidentRoutes(router *gin.RouterGroup) {
	endpoint := ConfigureBasicTenantEndpoint(router, "residents", domain.Resident{}, residentLinkFactory)
	ResidentBasePath = endpoint.BasePath()
}

func residentLinkFactory(halResource domain.HalResource, basePath string, detailed bool) {
	resident := halResource.(*domain.Resident)
	resident.AddLink(domain.REL_SELF, fmt.Sprintf("%s/%d", basePath, resident.ID))

	if detailed && resident.ApartmentID != nil {
		resident.AddLink(domain.REL_APARTMENT, fmt.Sprintf("%s/%d", ApartmentBasePath, *resident.ApartmentID))
	}
}
