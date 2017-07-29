package api

import (
	"github.com/cohousing/location"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/cohousing/cohousing-tenant-api/db"
)

func CreateRouter() {
	router := gin.Default()
	router.Use(ContextResolver())
	router.Use(location.Default())

	apiV1 := router.Group("api/v1")

	CreateHomeRoutes(apiV1)
	CreateApartmentRoutes(apiV1, dbFactory)
	CreateResidentRoutes(apiV1, dbFactory)
	CreateUserRoutes(apiV1, dbFactory)
	CreateGroupRoutes(apiV1, dbFactory)

	CreateFixtureRoutes(apiV1)

	router.Run(":8080")
}

func dbFactory(c *gin.Context) *gorm.DB {
	t := GetTenantFromContext(c)
	return db.GetTenantDB(t)
}
