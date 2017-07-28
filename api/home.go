package api

import (
	"fmt"
	"github.com/cohousing/cohousing-api/api/admin"
	"github.com/cohousing/cohousing-api/api/tenant"
	"github.com/cohousing/cohousing-api/api/utils"
	"github.com/cohousing/cohousing-api/domain"
	adminDomain "github.com/cohousing/cohousing-api/domain/admin"
	tenantDomain "github.com/cohousing/cohousing-api/domain/tenant"
	"github.com/cohousing/location"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateHomeRoutes(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		if utils.IsTenantRequest(c) {
			renderTenantHome(c, router.BasePath())
		} else {
			renderAdminHome(c, router.BasePath())
		}
	})
}

type TenantHome struct {
	Context string `json:"context"`
	ApiUrl  string `json:"apiurl"`
	Name    string `json:"name"`
	domain.DefaultHalResource
}

func renderTenantHome(c *gin.Context, basePath string) {
	t := utils.GetTenantFromContext(c)
	url := location.Get(c)

	tenantHome := &TenantHome{
		Context: t.Context,
		ApiUrl:  fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, basePath),
		Name:    t.Name,
	}

	tenantHome.AddLink(domain.REL_SELF, basePath)
	tenantHome.AddLink(tenantDomain.REL_APARTMENTS, tenant.ApartmentBasePath)
	tenantHome.AddLink(tenantDomain.REL_RESIDENTS, tenant.ResidentBasePath)
	tenantHome.AddLink(tenantDomain.REL_USERS, tenant.UserBasePath)
	tenantHome.AddLink(tenantDomain.REL_GROUPS, tenant.GroupBasePath)

	c.JSON(http.StatusOK, tenantHome)
}

type AdminHome struct {
	ApiUrl string `json:"apiurl"`
	domain.DefaultHalResource
}

func renderAdminHome(c *gin.Context, basePath string) {
	url := location.Get(c)

	adminHome := &AdminHome{
		ApiUrl: fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, basePath),
	}

	adminHome.AddLink(domain.REL_SELF, basePath)
	adminHome.AddLink(adminDomain.REL_TENANTS, admin.TenantAdminBasePath)

	c.JSON(http.StatusOK, adminHome)
}
