package api

import (
	"github.com/cohousing/cohousing-api/domain"
	"github.com/gin-gonic/gin"
)

func CreateApartmentRoutes(router *gin.RouterGroup) {
	endpoint := ConfigureBasicTenantEndpoint(router, "apartments", domain.Apartment{})
	domain.ApartmentBasePath = endpoint.BasePath()
}
