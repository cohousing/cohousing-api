package api

import (
	"fmt"
	domain2 "github.com/cohousing/cohousing-api-utils/domain"
	"github.com/cohousing/cohousing-tenant-api/domain"
	"github.com/cohousing/location"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TenantHome struct {
	Context string `json:"context"`
	ApiUrl  string `json:"apiurl"`
	Name    string `json:"name"`
	domain2.DefaultHalResource
}

func CreateHomeRoutes(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		basePath := router.BasePath()
		t := GetTenantFromContext(c)
		url := location.Get(c)

		tenantHome := &TenantHome{
			Context: t.Context,
			ApiUrl:  fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, basePath),
			Name:    t.Name,
		}

		tenantHome.AddLink(domain2.REL_SELF, basePath)

		if IsAuthenticated(c) {
			permission := ResolvePermission(c)
			if permission.GlobalAdmin || permission.ReadApartments {
				tenantHome.AddLink(domain.REL_APARTMENTS, ApartmentBasePath)
			}
			if permission.GlobalAdmin || permission.ReadResidents {
				tenantHome.AddLink(domain.REL_RESIDENTS, ResidentBasePath)
			}
			if permission.GlobalAdmin || permission.ReadUsers {
				tenantHome.AddLink(domain.REL_USERS, UserBasePath)
			}
			if permission.GlobalAdmin || permission.ReadGroups {
				tenantHome.AddLink(domain.REL_GROUPS, GroupBasePath)
			}
		}

		c.JSON(http.StatusOK, tenantHome)
	})
}
