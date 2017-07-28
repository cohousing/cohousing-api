package tenant

import (
	"github.com/cohousing/cohousing-tenant-api/api/utils"
	"github.com/gin-gonic/gin"
)

func ConfigureBasicTenantEndpoint(router *gin.RouterGroup, config utils.BasicEndpointConfig) *gin.RouterGroup {
	config.GetListHandlers = append([]gin.HandlerFunc{AuthorizeDomainObject(config.Domain, PERM_READ)}, config.GetListHandlers...)
	config.GetHandlers = append([]gin.HandlerFunc{AuthorizeDomainObject(config.Domain, PERM_READ)}, config.GetHandlers...)
	config.CreateHandlers = append([]gin.HandlerFunc{AuthorizeDomainObject(config.Domain, PERM_CREATE)}, config.CreateHandlers...)
	config.UpdateHandlers = append([]gin.HandlerFunc{AuthorizeDomainObject(config.Domain, PERM_UPDATE)}, config.UpdateHandlers...)
	config.DeleteHandlers = append([]gin.HandlerFunc{AuthorizeDomainObject(config.Domain, PERM_DELETE)}, config.DeleteHandlers...)

	return utils.ConfigureBasicEndpoint(router, config)
}
