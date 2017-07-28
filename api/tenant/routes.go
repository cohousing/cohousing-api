package tenant

import (
	"github.com/cohousing/cohousing-tenant-api/api/utils"
	"github.com/cohousing/cohousing-tenant-api/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func CreateTenantRoutes(router *gin.RouterGroup) {
	CreateApartmentRoutes(router, dbFactory)
	CreateResidentRoutes(router, dbFactory)
	CreateUserRoutes(router, dbFactory)
	CreateGroupRoutes(router, dbFactory)

	CreateFixtureRoutes(router)
}

func dbFactory(c *gin.Context) *gorm.DB {
	t := utils.GetTenantFromContext(c)
	return db.GetTenantDB(t)
}
