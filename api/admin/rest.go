package admin

import (
	"github.com/cohousing/cohousing-api/api/utils"
	"github.com/gin-gonic/gin"
)

func ConfigureBasicAdminEndpoint(router *gin.RouterGroup, config utils.BasicEndpointConfig) *gin.RouterGroup {
	return utils.ConfigureBasicEndpoint(router, config)
}
