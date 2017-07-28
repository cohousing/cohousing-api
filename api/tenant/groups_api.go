package tenant

import (
	"fmt"
	"github.com/cohousing/cohousing-tenant-api/api/utils"
	"github.com/cohousing/cohousing-tenant-api/db"
	"github.com/cohousing/cohousing-tenant-api/domain"
	"github.com/cohousing/cohousing-tenant-api/domain/tenant"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	GroupBasePath string
)

func CreateGroupRoutes(router *gin.RouterGroup, dbFactory db.DBFactory) {
	endpoint := ConfigureBasicTenantEndpoint(router, utils.BasicEndpointConfig{
		Path:           "groups",
		Domain:         tenant.Group{},
		DBFactory:      dbFactory,
		RouterHandlers: []gin.HandlerFunc{utils.MustBeTenant(), MustAuthenticate()},
	})
	GroupBasePath = endpoint.BasePath()

	utils.AddLinkFactory(&tenant.Group{}, groupLinkFactory)

	endpoint.GET("/:id/users", AuthorizeDomainObject(tenant.Group{}, PERM_READ), getUsersForGroup(dbFactory))
}

func groupLinkFactory(c *gin.Context, halResource domain.HalResource, basePath string, detailed bool) {
	u := halResource.(*tenant.Group)
	u.AddLink(domain.REL_SELF, fmt.Sprintf("%s/%d", basePath, u.ID))

	if detailed {
		permission := ResolvePermission(c)
		if permission.GlobalAdmin || permission.UpdateGroups {
			u.AddLink(domain.REL_UPDATE, fmt.Sprintf("%s/%d", basePath, u.ID))
		}
		if permission.GlobalAdmin || permission.DeleteGroups {
			u.AddLink(domain.REL_DELETE, fmt.Sprintf("%s/%d", basePath, u.ID))
		}
		if permission.GlobalAdmin || permission.ReadUsers {
			u.AddLink(tenant.REL_USERS, fmt.Sprintf("%s/%d/users", basePath, u.ID))
		}
	}
}

func getUsersForGroup(dbFactory db.DBFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
			u := tenant.User{}
			u.ID = id
			var users []tenant.User
			dbFactory(c).Model(&u).Related(&users, "Users")
			utils.AddLinks(c, users, UserBasePath, false)
			c.JSON(http.StatusOK, users)
		}
	}
}
