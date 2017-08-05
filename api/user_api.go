package api

import (
	"fmt"
	"github.com/cohousing/cohousing-api-utils/api"
	"github.com/cohousing/cohousing-api-utils/db"
	domain2 "github.com/cohousing/cohousing-api-utils/domain"
	"github.com/cohousing/cohousing-tenant-api/domain"
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
		Domain:         domain.User{},
		DBFactory:      dbFactory,
		RouterHandlers: []gin.HandlerFunc{MustBeTenant(), MustAuthenticate()},
	})
	UserBasePath = endpoint.BasePath()

	utils.AddLinkFactory(&domain.User{}, userLinkFactory)

	endpoint.GET("/:id/groups", AuthorizeDomainObject(domain.User{}, PERM_READ), getGroupsForUser(dbFactory))
}

func userLinkFactory(c *gin.Context, halResource domain2.HalResource, basePath string, detailed bool) {
	u := halResource.(*domain.User)
	u.AddLink(domain2.REL_SELF, fmt.Sprintf("%s/%d", basePath, u.ID))

	if detailed {
		permission := GetPermissionsFromContext(c)
		if permission.GlobalAdmin || permission.UpdateUsers {
			u.AddLink(domain2.REL_UPDATE, fmt.Sprintf("%s/%d", basePath, u.ID))
		}
		if permission.GlobalAdmin || permission.DeleteUsers {
			u.AddLink(domain2.REL_DELETE, fmt.Sprintf("%s/%d", basePath, u.ID))
		}
		if permission.GlobalAdmin || permission.ReadGroups {
			u.AddLink(domain.REL_GROUPS, fmt.Sprintf("%s/%d/groups", basePath, u.ID))
		}
	}
}

func getGroupsForUser(dbFactory db.DBFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
			u := domain.User{}
			u.ID = id
			var groups []domain.Group
			dbFactory(c).Model(&u).Related(&groups, "Groups")
			utils.AddLinks(c, groups, GroupBasePath, false)
			c.JSON(http.StatusOK, groups)
		}
	}
}
