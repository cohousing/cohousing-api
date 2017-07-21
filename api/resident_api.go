package api

import (
	"github.com/cohousing/cohousing-api/domain"
	"github.com/gin-gonic/gin"
)

func CreateResidentRoutes(router *gin.RouterGroup) {
	endpoint := ConfigureBasicTenantEndpoint(router, "residents", domain.Resident{})
	domain.ResidentBasePath = endpoint.BasePath()
}
