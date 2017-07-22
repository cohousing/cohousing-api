package tenant

import (
	"github.com/cohousing/cohousing-api/api/utils"
	"github.com/cohousing/cohousing-api/db"
	"github.com/gin-gonic/gin"
)

func ConfigureBasicTenantEndpoint(router *gin.RouterGroup, path string, domain interface{}, linkFactory utils.LinkFactory, dbFactory db.DBFactory) *gin.RouterGroup {
	return utils.ConfigureBasicEndpoint(router, path, domain, linkFactory, dbFactory, utils.MustBeTenant())
}
