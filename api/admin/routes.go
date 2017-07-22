package admin

import (
	"github.com/cohousing/cohousing-api/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func CreateAdminRoutes(router *gin.RouterGroup) {
	CreateTenantAdminRoutes(router, dbFactory)
}

func dbFactory(c *gin.Context) *gorm.DB {
	return db.GetConfDB()
}
