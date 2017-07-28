package tenant

import (
	"fmt"
	"github.com/cohousing/cohousing-api/api/utils"
	"github.com/cohousing/cohousing-api/db"
	"github.com/cohousing/cohousing-api/domain"
	"github.com/cohousing/cohousing-api/domain/tenant"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	UserBasePath string
)

func CreateUserRoutes(router *gin.RouterGroup, dbFactory db.DBFactory) {
	endpoint := ConfigureBasicTenantEndpoint(router, utils.BasicEndpointConfig{
		Path:           "users",
		Domain:         tenant.User{},
		DBFactory:      dbFactory,
		RouterHandlers: []gin.HandlerFunc{utils.MustBeTenant(), MustAuthenticate()},
	})
	UserBasePath = endpoint.BasePath()

	utils.AddLinkFactory(&tenant.User{}, userLinkFactory)

	endpoint.GET("/:id/groups", AuthorizeDomainObject(tenant.User{}, PERM_READ), getGroupsForUser(dbFactory))
}

func userLinkFactory(c *gin.Context, halResource domain.HalResource, basePath string, detailed bool) {
	u := halResource.(*tenant.User)
	u.AddLink(domain.REL_SELF, fmt.Sprintf("%s/%d", basePath, u.ID))

	if detailed {
		permission := ResolvePermission(c)
		if permission.GlobalAdmin || permission.UpdateUsers {
			u.AddLink(domain.REL_UPDATE, fmt.Sprintf("%s/%d", basePath, u.ID))
		}
		if permission.GlobalAdmin || permission.DeleteUsers {
			u.AddLink(domain.REL_DELETE, fmt.Sprintf("%s/%d", basePath, u.ID))
		}
		if permission.GlobalAdmin || permission.ReadGroups {
			u.AddLink(tenant.REL_GROUPS, fmt.Sprintf("%s/%d/groups", basePath, u.ID))
		}
	}
}

func getGroupsForUser(dbFactory db.DBFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
			u := tenant.User{}
			u.ID = id
			var groups []tenant.Group
			dbFactory(c).Model(&u).Related(&groups, "Groups")
			utils.AddLinks(c, groups, GroupBasePath, false)
			c.JSON(http.StatusOK, groups)
		}
	}
}
